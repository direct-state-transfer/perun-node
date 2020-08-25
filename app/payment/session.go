// Copyright (c) 2020 - for information on the respective copyright owner
// see the NOTICE file and/or the repository at
// https://github.com/hyperledger-labs/perun-node
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package payment

import (
	"context"
	"fmt"
	"math/big"

	ppayment "perun.network/go-perun/apps/payment"
	pchannel "perun.network/go-perun/channel"

	"github.com/hyperledger-labs/perun-node"
	"github.com/hyperledger-labs/perun-node/blockchain/ethereum"
	"github.com/hyperledger-labs/perun-node/currency"
)

// init() initializes the payment app in go-perun.
func init() {
	wb := ethereum.NewWalletBackend()
	emptyAddr, err := wb.ParseAddr("0x0")
	if err != nil {
		panic("Error parsing zero address for app payment def: " + err.Error())
	}
	ppayment.SetAppDef(emptyAddr) // dummy app def.
}

// OpenPayCh opens a payment channel using the given sessionAPI instance with the specified parameters.
func OpenPayCh(pctx context.Context,
	s perun.SessionAPI,
	peerAlias string,
	openingBals perun.BalInfo,
	challengeDurSecs uint64) (PayChInfo, error) {
	paymentApp := perun.App{
		Def:  ppayment.AppDef(),
		Data: &ppayment.NoData{},
	}

	chInfo, err := s.OpenCh(pctx, peerAlias, openingBals, paymentApp, challengeDurSecs)
	if err != nil {
		return PayChInfo{}, err
	}
	return PayChInfo{
		ChannelID: chInfo.ChannelID,
		BalInfo:   balsFromState(chInfo.Currency, chInfo.State, chInfo.Parts),
		Version:   fmt.Sprintf("%d", chInfo.State.Version),
	}, nil
}

// GetPayChs returns a list of payment channel info for all the channels in this session.
func GetPayChs(s perun.SessionAPI) []PayChInfo {
	chInfos := s.GetChInfos()

	payChInfos := make([]PayChInfo, len(chInfos))
	for i := range chInfos {
		payChInfos[i] = PayChInfo{
			ChannelID: chInfos[i].ChannelID,
			BalInfo:   balsFromState(chInfos[i].Currency, chInfos[i].State, chInfos[i].Parts),
			Version:   fmt.Sprintf("%d", chInfos[i].State.Version),
		}
	}
	return payChInfos
}

func balsFromState(currency string, state *pchannel.State, parts []string) perun.BalInfo {
	return balsFromBigInt(currency, state.Balances[0], parts)
}

func balsFromBigInt(chCurrency string, bigInt []*big.Int, parts []string) perun.BalInfo {
	balInfo := perun.BalInfo{
		Currency: chCurrency,
		Bals:     make(map[string]string, len(parts)),
	}

	parser := currency.NewParser(chCurrency)
	for i := range parts {
		balInfo.Bals[parts[i]] = parser.Print(bigInt[i])
		balInfo.Bals[parts[i]] = parser.Print(bigInt[i])
	}
	return balInfo
}
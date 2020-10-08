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

package internal

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/pkg/errors"
	pethchannel "perun.network/go-perun/backend/ethereum/channel"
	pethwallet "perun.network/go-perun/backend/ethereum/wallet"
	pchannel "perun.network/go-perun/channel"
	pwallet "perun.network/go-perun/wallet"
)

// ChainBackend provides ethereum specific contract backend functionality.
type ChainBackend struct {
	// Cb is the instance of contract backend that will be used for all on-chain communications.
	Cb *pethchannel.ContractBackend
	// TxTimeout is the max time to wait for confirmation of transactions on blockchain.
	// If this expires, a transactions is considered failed.
	// Use sufficiently large values when connecting to mainnet.
	TxTimeout time.Duration
}

// NewFunder initializes and returns an instance of ethereum funder.
func (cb *ChainBackend) NewFunder(assetAddr pwallet.Address, cred pwallet.Address) pchannel.Funder {
	asset := pethwallet.AsWalletAddr(pethwallet.AsEthAddr(assetAddr))
	acc := accounts.Account{Address: pethwallet.AsEthAddr(cred)}
	accounts := map[pethchannel.Asset]accounts.Account{*asset: acc}
	depositors := map[pethchannel.Asset]pethchannel.Depositor{*asset: new(pethchannel.ETHDepositor)}
	return pethchannel.NewFunder(*cb.Cb, accounts, depositors)
}

// NewAdjudicator initializes and returns an instance of ethereum adjudicator.
func (cb *ChainBackend) NewAdjudicator(adjAddr, acct pwallet.Address) pchannel.Adjudicator {
	acc := accounts.Account{Address: pethwallet.AsEthAddr(acct)}
	return pethchannel.NewAdjudicator(*cb.Cb, pethwallet.AsEthAddr(adjAddr),
		pethwallet.AsEthAddr(acct), acc)
}

// ValidateContracts validates the integrity of given adjudicator and asset holder contracts.
func (cb *ChainBackend) ValidateContracts(adjAddr, assetAddr pwallet.Address) error {
	ctx, cancel := context.WithTimeout(context.Background(), cb.TxTimeout)
	defer cancel()

	// Integrity of Adjudicator is implicitly done during validation of asset holder contract.
	err := pethchannel.ValidateAssetHolderETH(ctx, *cb.Cb, pethwallet.AsEthAddr(assetAddr), pethwallet.AsEthAddr(adjAddr))
	if pethchannel.IsContractBytecodeError(err) {
		return errors.Wrap(err, "invalid contracts at given addresses")
	}
	return errors.Wrap(err, "validating contracts")
}

// DeployAdjudicator deploys the adjudicator contract.
func (cb *ChainBackend) DeployAdjudicator(onChainAddr pwallet.Address) (pwallet.Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cb.TxTimeout)
	defer cancel()

	acc := accounts.Account{Address: pethwallet.AsEthAddr(onChainAddr)}
	addr, err := pethchannel.DeployAdjudicator(ctx, *cb.Cb, acc)
	return pethwallet.AsWalletAddr(addr), errors.Wrap(err, "deploying adjudicator contract")
}

// DeployAsset deploys the asset holder contract, setting the adjudicator address to given value.
func (cb *ChainBackend) DeployAsset(adjAddr, cred pwallet.Address) (pwallet.Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cb.TxTimeout)
	defer cancel()

	acc := accounts.Account{Address: pethwallet.AsEthAddr(cred)}
	addr, err := pethchannel.DeployETHAssetholder(ctx, *cb.Cb, pethwallet.AsEthAddr(adjAddr), acc)
	return pethwallet.AsWalletAddr(addr), errors.Wrap(err, "deploying asset contract")
} // nolint:gofumpt // unknown error, maybe a false positive

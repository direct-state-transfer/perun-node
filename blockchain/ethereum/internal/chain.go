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
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	pperuntoken "perun.network/go-perun/backend/ethereum/bindings/peruntoken"
	pethchannel "perun.network/go-perun/backend/ethereum/channel"
	pethwallet "perun.network/go-perun/backend/ethereum/wallet"
	pchannel "perun.network/go-perun/channel"
	pwallet "perun.network/go-perun/wallet"

	"github.com/hyperledger-labs/perun-node/blockchain"
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
func (cb *ChainBackend) NewFunder(assetETHAddr pwallet.Address, txSender pwallet.Address) pchannel.Funder {
	assetETH := pethwallet.AsWalletAddr(pethwallet.AsEthAddr(assetETHAddr))
	txSenderAcc := accounts.Account{Address: pethwallet.AsEthAddr(txSender)}
	txSenderAccs := map[pethchannel.Asset]accounts.Account{*assetETH: txSenderAcc}
	depositors := map[pethchannel.Asset]pethchannel.Depositor{*assetETH: new(pethchannel.ETHDepositor)}
	return pethchannel.NewFunder(*cb.Cb, txSenderAccs, depositors)
}

// NewAdjudicator initializes and returns an instance of ethereum adjudicator.
func (cb *ChainBackend) NewAdjudicator(adjAddr, txSender pwallet.Address) pchannel.Adjudicator {
	txSenderAcc := accounts.Account{Address: pethwallet.AsEthAddr(txSender)}
	return pethchannel.NewAdjudicator(*cb.Cb, pethwallet.AsEthAddr(adjAddr),
		pethwallet.AsEthAddr(txSender), txSenderAcc)
}

// ERC20Info reads the symbol and number of decimal values decimals for the
// ERC20 token from the blockchain.
func (cb *ChainBackend) ERC20Info(addr pwallet.Address) (symbol string, decimal uint8, _ error) {
	tokenERC20, err := pperuntoken.NewPerunToken(pethwallet.AsEthAddr(addr), *cb.Cb)
	if err != nil {
		// This errors only when an invalid ABI JSON is provided when binding.
		// As we include the ABI JSON in the compiled binary, it never errors.
		return "", 0, errors.Wrap(err, "binding to erc20 token contract")
	}

	opts := &bind.CallOpts{
		Pending: false,
		Context: context.Background(),
	}
	symbol, err = tokenERC20.Symbol(opts)
	if err != nil {
		return "", 0, errors.Wrap(err, "reading symbol from the contract")
	}

	decimal, err = tokenERC20.Decimals(opts)
	return symbol, decimal, errors.Wrap(err, "reading decimals from the contract")
}

// ValidateAdjudicator validates the integrity of adjudicator contract at the
// given address.
func (cb *ChainBackend) ValidateAdjudicator(adjAddr pwallet.Address) error {
	ctx, cancel := context.WithTimeout(context.Background(), cb.TxTimeout)
	defer cancel()
	err := pethchannel.ValidateAdjudicator(ctx, *cb.Cb, pethwallet.AsEthAddr(adjAddr))
	if pethchannel.IsErrInvalidContractCode(err) {
		return blockchain.NewInvalidContractError(blockchain.Adjudicator, adjAddr.String(), err)
	}
	return errors.Wrap(err, "validating adjudicator contract")
}

// ValidateAssetETH validates the integrity of adjudicator and asset ETH
// contracts at the given addresses.
//
// TODO: Submit a suggestion to go-perun to not validate the adjudicator contract in ValidateAssetHolder.
// If accepted, then update this function to
// validate only the asset ETH contract.
func (cb *ChainBackend) ValidateAssetETH(adjAddr, assetETHAddr pwallet.Address) error {
	ctx, cancel := context.WithTimeout(context.Background(), cb.TxTimeout)
	defer cancel()
	// Though integrity of adjudicator is implicitly checked by ValidateAssetHolderETH,
	// we do it before that call to identify this type of error.
	err := pethchannel.ValidateAdjudicator(ctx, *cb.Cb, pethwallet.AsEthAddr(adjAddr))
	if pethchannel.IsErrInvalidContractCode(err) {
		return blockchain.NewInvalidContractError(blockchain.Adjudicator, adjAddr.String(), err)
	}

	err = pethchannel.ValidateAssetHolderETH(ctx, *cb.Cb,
		pethwallet.AsEthAddr(assetETHAddr), pethwallet.AsEthAddr(adjAddr))
	if pethchannel.IsErrInvalidContractCode(err) {
		return blockchain.NewInvalidContractError(blockchain.AssetETH, assetETHAddr.String(), err)
	}
	return errors.Wrap(err, "validating asset ETH contract")
}

// ValidateAssetERC20 validates the integrity of adjudicator and asset ERC20
// contracts at the given addresses. TokenERC20 is the address of ERC20 token
// contract.
//
// TODO: Submit a suggestion to go-perun to not validate the adjudicator contract in ValidateAssetHolder.
// If accepted, then update this function to validate only the asset ERC20 contract.
func (cb *ChainBackend) ValidateAssetERC20(adj, tokenERC20, assetERC20 pwallet.Address) (
	symbol string, decimals uint8, _ error) {
	ctx, cancel := context.WithTimeout(context.Background(), cb.TxTimeout)
	defer cancel()
	var err error
	symbol, decimals, err = cb.ERC20Info(tokenERC20)
	if err != nil {
		return "", 0, errors.WithMessage(err, "reading symbol and decimal values from token contract")
	}
	// Though integrity of adjudicator is implicitly checked by ValidateAssetHolderERC20,
	// we do it before that call to identify this type of error.
	err = pethchannel.ValidateAdjudicator(ctx, *cb.Cb, pethwallet.AsEthAddr(adj))
	if pethchannel.IsErrInvalidContractCode(err) {
		return "", 0, blockchain.NewInvalidContractError(blockchain.Adjudicator, adj.String(), err)
	}
	err = pethchannel.ValidateAssetHolderERC20(ctx, *cb.Cb,
		pethwallet.AsEthAddr(assetERC20), pethwallet.AsEthAddr(adj), pethwallet.AsEthAddr(tokenERC20))
	if err != nil && pethchannel.IsErrInvalidContractCode(err) {
		return "", 0, blockchain.NewInvalidContractError(symbol, assetERC20.String(), err)
	}

	return symbol, decimals, errors.Wrap(err, "validating asset ERC20 contract")
}

// DeployAdjudicator deploys the adjudicator contract.
func (cb *ChainBackend) DeployAdjudicator(txSender pwallet.Address) (pwallet.Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cb.TxTimeout)
	defer cancel()

	txSenderAcc := accounts.Account{Address: pethwallet.AsEthAddr(txSender)}
	addr, err := pethchannel.DeployAdjudicator(ctx, *cb.Cb, txSenderAcc)
	return pethwallet.AsWalletAddr(addr), errors.Wrap(err, "deploying adjudicator contract")
}

// DeployAssetETH deploys the asset ETH contract, setting the adjudicator address to given value.
func (cb *ChainBackend) DeployAssetETH(adjAddr, txSender pwallet.Address) (pwallet.Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cb.TxTimeout)
	defer cancel()

	txSenderAcc := accounts.Account{Address: pethwallet.AsEthAddr(txSender)}
	addr, err := pethchannel.DeployETHAssetholder(ctx, *cb.Cb, pethwallet.AsEthAddr(adjAddr), txSenderAcc)
	return pethwallet.AsWalletAddr(addr), errors.Wrap(err, "deploying asset ETH contract")
}

// DeployPerunToken deploys the perun ERC20 token contract.
func (cb *ChainBackend) DeployPerunToken(initAccs []pwallet.Address, initBal *big.Int, txSender pwallet.Address) (
	pwallet.Address, error) {
	initAccsETH := make([]common.Address, len(initAccs))
	for i := range initAccs {
		initAccsETH[i] = pethwallet.AsEthAddr(initAccs[i])
	}

	ctx, cancel := context.WithTimeout(context.Background(), cb.TxTimeout)
	defer cancel()
	txSenderAcc := accounts.Account{Address: pethwallet.AsEthAddr(txSender)}
	addr, err := pethchannel.DeployPerunToken(ctx, *cb.Cb, txSenderAcc, initAccsETH, initBal)
	return pethwallet.AsWalletAddr(addr), errors.Wrap(err, "deploying adjudicator contract")
}

// DeployAssetERC20 deploys the asset ERC20 contract, setting the adjudicator
// and erc20 token addresses to given values.
func (cb *ChainBackend) DeployAssetERC20(adj, tokenERC20, txSender pwallet.Address) (pwallet.Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cb.TxTimeout)
	defer cancel()

	txSenderAcc := accounts.Account{Address: pethwallet.AsEthAddr(txSender)}
	addr, err := pethchannel.DeployERC20Assetholder(ctx, *cb.Cb,
		pethwallet.AsEthAddr(adj), pethwallet.AsEthAddr(tokenERC20), txSenderAcc)
	return pethwallet.AsWalletAddr(addr), errors.Wrap(err, "deploying asset ERC20 contract")
}

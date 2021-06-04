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

package session

import (
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type (
	// Config defines the parameters required to configure a session.
	Config struct {
		User UserConfig

		IDProviderType     string        // Type of ID provider.
		IDProviderURL      string        // URL for accessing the ID provider.
		ChainURL           string        // URL of the blockchain node.
		ChainID            int           // See client/config.go.
		Asset, Adjudicator string        // Address of the Asset and Adjudicator contracts.
		ChainConnTimeout   time.Duration // Timeout for connecting to blockchain node.
		OnChainTxTimeout   time.Duration // Timeout to wait for confirmation of on-chain tx.
		ResponseTimeout    time.Duration // Timeout to wait for a response from the peer / user.

		DatabaseDir string // Path to directory containing persistence database.
		// Timeout for re-establishing all open channels (if any) that was persisted during the
		// previous running instance of the node.
		PeerReconnTimeout time.Duration
	}

	// UserConfig defines the parameters required to configure a user.
	// Address strings should be parsed using the wallet backend.
	UserConfig struct {
		Alias string

		OnChainAddr   string
		OnChainWallet WalletConfig

		PartAddrs      []string
		OffChainAddr   string
		OffChainWallet WalletConfig

		CommAddr string
		CommType string
	}

	// WalletConfig defines the parameters required to configure a wallet.
	WalletConfig struct {
		KeystorePath string
		Password     string
	}

	// clientConfig represents the configuration parameters for state channel client.
	clientConfig struct {
		Chain chainConfig

		// Path to directory containing persistence database.
		DatabaseDir string
		// Timeout for re-establishing all open channels (if any) that was persisted during the
		// previous running instance of the node.
		PeerReconnTimeout time.Duration
	}

	// chainConfig represents the configuration parameters for connecting to blockchain.
	chainConfig struct {
		// Addresses of on-chain contracts used for establishing state channel network.
		Adjudicator string
		Asset       string

		// URL for connecting to the blockchain node.
		URL string
		// ChainID is the unique identifier for different chains in the ethereum ecosystem.
		ChainID int
		// ConnTimeout is the timeout used when dialing for new connections to the on-chain node.
		ConnTimeout time.Duration
		// OnChainTxTimeout is the timeout to wait for a blockchain transaction to be finalized.
		OnChainTxTimeout time.Duration
	}
)

// ParseConfig parses the session configuration from a file.
func ParseConfig(configFile string) (Config, error) {
	v := viper.New()
	v.SetConfigFile(filepath.Clean(configFile))

	var cfg Config
	err := v.ReadInConfig()
	if err != nil {
		return Config{}, errors.Wrap(err, "reading from source")
	}
	return cfg, errors.Wrap(v.Unmarshal(&cfg), "unmarshalling")
}

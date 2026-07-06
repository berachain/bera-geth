// Copyright 2026 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package utils

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/node"
	"github.com/urfave/cli/v2"
)

func testCLIContext(t *testing.T, flags []cli.Flag, args ...string) *cli.Context {
	t.Helper()

	app := cli.NewApp()
	app.Flags = flags
	set := flag.NewFlagSet("test", flag.ContinueOnError)
	for _, f := range flags {
		if err := f.Apply(set); err != nil {
			t.Fatalf("failed to apply flag: %v", err)
		}
	}
	if err := set.Parse(args); err != nil {
		t.Fatalf("failed to parse args: %v", err)
	}
	return cli.NewContext(app, set, nil)
}

var setEthConfigTestFlags = []cli.Flag{
	BerachainFlag,
	BepoliaFlag,
	MainnetFlag,
	SepoliaFlag,
	HoleskyFlag,
	HoodiFlag,
	DeveloperFlag,
	OverrideGenesisFlag,
	NetworkIdFlag,
	GCModeFlag,
	CacheFlag,
	CacheDatabaseFlag,
	CacheTrieFlag,
	CacheGCFlag,
	CacheSnapshotFlag,
	CachePreimagesFlag,
	SnapshotFlag,
	CryptoKZGFlag,
	FDLimitFlag,
	StateSizeTrackingFlag,
	VMWitnessStatsFlag,
}

func TestMakeGenesisNetworkPresets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		args        []string
		wantChainID uint64
	}{
		{name: "default", args: nil, wantChainID: 80094},
		{name: "berachain", args: []string{"--berachain"}, wantChainID: 80094},
		{name: "bepolia", args: []string{"--bepolia"}, wantChainID: 80069},
		{name: "mainnet", args: []string{"--mainnet"}, wantChainID: 1},
		{name: "sepolia", args: []string{"--sepolia"}, wantChainID: 11155111},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := testCLIContext(t, NetworkFlags, tt.args...)
			genesis := MakeGenesis(ctx)
			if genesis == nil || genesis.Config == nil {
				t.Fatal("expected non-nil genesis config")
			}
			if got := genesis.Config.ChainID.Uint64(); got != tt.wantChainID {
				t.Fatalf("chain id mismatch: got %d, want %d", got, tt.wantChainID)
			}
		})
	}
}

func TestSetEthConfigNetworkPresets(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantNet     uint64
		wantChainID uint64
	}{
		{name: "default", args: nil, wantNet: 0, wantChainID: 80094},
		{name: "berachain", args: []string{"--berachain"}, wantNet: 80094, wantChainID: 80094},
		{name: "bepolia", args: []string{"--bepolia"}, wantNet: 80069, wantChainID: 80069},
		{name: "mainnet", args: []string{"--mainnet"}, wantNet: 1, wantChainID: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ncfg := node.DefaultConfig
			ncfg.DataDir = "" // keep tests hermetic; avoid touching the host filesystem
			stack, err := node.New(&ncfg)
			if err != nil {
				t.Fatalf("failed to create node: %v", err)
			}
			defer stack.Close()

			ctx := testCLIContext(t, setEthConfigTestFlags, tt.args...)
			cfg := ethconfig.Defaults
			SetEthConfig(ctx, stack, &cfg)
			if cfg.NetworkId != tt.wantNet {
				t.Fatalf("network id mismatch: got %d, want %d", cfg.NetworkId, tt.wantNet)
			}
			if cfg.Genesis == nil || cfg.Genesis.Config == nil {
				t.Fatal("expected non-nil genesis config")
			}
			if got := cfg.Genesis.Config.ChainID.Uint64(); got != tt.wantChainID {
				t.Fatalf("chain id mismatch: got %d, want %d", got, tt.wantChainID)
			}
		})
	}
}

func TestSetDataDirBepolia(t *testing.T) {
	t.Parallel()

	ctx := testCLIContext(t, NetworkFlags, "--bepolia")
	cfg := node.DefaultConfig
	SetDataDir(ctx, &cfg)
	want := filepath.Join(node.DefaultDataDir(), "bepolia")
	if cfg.DataDir != want {
		t.Fatalf("datadir mismatch: got %s, want %s", cfg.DataDir, want)
	}
}

func TestIsNetworkPresetBerachainBepolia(t *testing.T) {
	t.Parallel()

	if IsNetworkPreset(testCLIContext(t, NetworkFlags)) {
		t.Fatal("expected no preset without flags")
	}
	if !IsNetworkPreset(testCLIContext(t, NetworkFlags, "--berachain")) {
		t.Fatal("expected berachain to be a network preset")
	}
	if !IsNetworkPreset(testCLIContext(t, NetworkFlags, "--bepolia")) {
		t.Fatal("expected bepolia to be a network preset")
	}
}

func TestDefaultsBerachainGenesis(t *testing.T) {
	t.Parallel()

	// Genesis is resolved lazily in SetEthConfig to avoid decoding the large
	// Berachain prealloc at ethconfig package init time.
	if ethconfig.Defaults.Genesis != nil {
		t.Fatal("expected nil genesis in Defaults; use SetEthConfig to resolve")
	}
}

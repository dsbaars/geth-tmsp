package main

// NOTE: This list of flags was co-opted from geth at github.com/ethereum/go-ethereum,
// but we only run the daemon, and don't offer any of the other commands.
// The resulting application should still work
// (ie. `geth attach` on a `geth-tmsp` daemon should work.
// but anything related to the chain itself won't (ie. now need to talk to the tendermint daemon)

import (
	"os"

	"github.com/eris-ltd/geth-tmsp/app"
	"github.com/ethereum/go-ethereum/cmd/utils"
	. "github.com/tendermint/go-common"
	"github.com/tendermint/tmsp/server"

	"github.com/codegangsta/cli"
)

const Version = "1.3.3-tmsp"

func main() {
	app := utils.NewApp(Version, "the geth-tmsp command line interface")
	app.Action = run
	app.Flags = []cli.Flag{
		utils.IdentityFlag,
		utils.UnlockedAccountFlag,
		utils.PasswordFileFlag,
		utils.GenesisFileFlag,
		utils.BootnodesFlag,
		utils.DataDirFlag,
		utils.BlockchainVersionFlag,
		utils.OlympicFlag,
		utils.FastSyncFlag,
		utils.CacheFlag,
		utils.LightKDFFlag,
		utils.JSpathFlag,
		utils.ListenPortFlag,
		utils.MaxPeersFlag,
		utils.MaxPendingPeersFlag,
		utils.EtherbaseFlag,
		utils.GasPriceFlag,
		utils.MinerThreadsFlag,
		utils.MiningEnabledFlag,
		utils.MiningGPUFlag,
		utils.AutoDAGFlag,
		utils.NATFlag,
		utils.NatspecEnabledFlag,
		utils.NoDiscoverFlag,
		utils.NodeKeyFileFlag,
		utils.NodeKeyHexFlag,
		utils.RPCEnabledFlag,
		utils.RPCListenAddrFlag,
		utils.RPCPortFlag,
		utils.RpcApiFlag,
		utils.IPCDisabledFlag,
		utils.IPCApiFlag,
		utils.IPCPathFlag,
		utils.ExecFlag,
		utils.WhisperEnabledFlag,
		utils.DevModeFlag,
		utils.TestNetFlag,
		utils.VMDebugFlag,
		utils.VMForceJitFlag,
		utils.VMJitCacheFlag,
		utils.VMEnableJitFlag,
		utils.NetworkIdFlag,
		utils.RPCCORSDomainFlag,
		utils.VerbosityFlag,
		utils.BacktraceAtFlag,
		utils.LogVModuleFlag,
		utils.LogFileFlag,
		utils.PProfEanbledFlag,
		utils.PProfPortFlag,
		utils.MetricsEnabledFlag,
		utils.SolcPathFlag,
		utils.GpoMinGasPriceFlag,
		utils.GpoMaxGasPriceFlag,
		utils.GpoFullBlockRatioFlag,
		utils.GpobaseStepDownFlag,
		utils.GpobaseStepUpFlag,
		utils.GpobaseCorrectionFactorFlag,
		utils.ExtraDataFlag,
	}
	app.Before = func(ctx *cli.Context) error {
		utils.SetupLogger(ctx)
		utils.SetupNetwork(ctx)
		utils.SetupVM(ctx)
		if ctx.GlobalBool(utils.PProfEanbledFlag.Name) {
			utils.StartPProf(ctx)
		}
		return nil
	}
	app.Run(os.Args)
}

func run(ctx *cli.Context) {
	// Start the tmsp listener
	// all flags are passed through to the app
	_, err := server.StartListener("tcp://0.0.0.0:46658", app.NewEthereumApplication(ctx))
	if err != nil {
		Exit(err.Error())
	}

	// Wait forever
	TrapSignal(func() {
		// TODO: Cleanup
	})
}

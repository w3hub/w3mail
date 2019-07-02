package main

import (
	"encoding/json"
	"io"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/auth/genaccounts"
	genaccscli "github.com/cosmos/cosmos-sdk/x/auth/genaccounts/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/w3hub/w3mail"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	tmtypes "github.com/tendermint/tendermint/types"
)

func main() {
	cobra.EnableCommandSorting = false

	cdc := w3mail.MakeCodec()

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "nsd",
		Short:             "nameservice App Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}
	// CLI commands to initialize the chain
	rootCmd.AddCommand(
		genutilcli.InitCmd(ctx, cdc, w3mail.ModuleBasics, w3mail.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(ctx, cdc, genaccounts.AppModuleBasic{}, w3mail.DefaultNodeHome),
		genutilcli.GenTxCmd(ctx, cdc, w3mail.ModuleBasics, staking.AppModuleBasic{}, genaccounts.AppModuleBasic{}, w3mail.DefaultNodeHome, w3mail.DefaultCLIHome),
		genutilcli.ValidateGenesisCmd(ctx, cdc, w3mail.ModuleBasics),
		// AddGenesisAccountCmd allows users to add accounts to the genesis file
		genaccscli.AddGenesisAccountCmd(ctx, cdc, w3mail.DefaultNodeHome, w3mail.DefaultCLIHome),
	)

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "NS", w3mail.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return w3mail.NewW3MailServiceApp(logger, db)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

	if height != -1 {
		w3mailApp := w3mail.NewW3MailServiceApp(logger, db)
		err := w3mailApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return w3mailApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	w3mailApp := w3mail.NewW3MailServiceApp(logger, db)

	return w3mailApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}

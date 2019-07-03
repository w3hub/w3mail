package w3mail

import (
	"encoding/json"

	"github.com/tendermint/tendermint/libs/log"

	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	dbm "github.com/tendermint/tendermint/libs/db"
	tmtypes "github.com/tendermint/tendermint/types"
)

const appName = "w3mail"

var (
	// DefaultCLIHome default home directories for the application CLI
	DefaultCLIHome = os.ExpandEnv("$HOME/.w3mailcli")

	// DefaultNodeHome sets the folder where the applcation data and configuration will be stored
	DefaultNodeHome = os.ExpandEnv("$HOME/.w3mail")

	// ModuleBasics The ModuleBasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics module.BasicManager
)

func init() {
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		slashing.AppModuleBasic{},
		distribution.AppModuleBasic{},
	)
}

// ServiceApp .
type ServiceApp struct {
	*baseapp.BaseApp
	cdc *codec.Codec

	keyMain    *sdk.KVStoreKey
	keyParams  *sdk.KVStoreKey
	tkeyParams *sdk.TransientStoreKey

	accountKeeper       auth.AccountKeeper
	bankKeeper          bank.Keeper
	stakingKeeper       staking.Keeper
	slashingKeeper      slashing.Keeper
	distrKeeper         distribution.Keeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	paramsKeeper        params.Keeper
}

// MakeCodec custom tx codec
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

// NewW3MailServiceApp .
func NewW3MailServiceApp(logger log.Logger, db dbm.DB) *ServiceApp {

	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := baseapp.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc))

	app := &ServiceApp{
		BaseApp:    bApp,
		cdc:        cdc,
		keyMain:    sdk.NewKVStoreKey(baseapp.MainStoreKey),
		keyParams:  sdk.NewKVStoreKey(params.StoreKey),
		tkeyParams: sdk.NewTransientStoreKey(params.TStoreKey),
	}

	app.paramsKeeper = params.NewKeeper(app.cdc, app.keyParams, app.tkeyParams, params.DefaultCodespace)

	return app
}

// LoadHeight .
func (app *ServiceApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keyMain)
}

// ExportAppStateAndValidators .
func (app *ServiceApp) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {

	// as if they could withdraw from the start of the next block
	// ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	// if forZeroHeight {
	// 	app.prepForZeroHeightGenesis(ctx, jailWhiteList)
	// }

	// genState := app.mm.ExportGenesis(ctx)
	// appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	// if err != nil {
	// 	return nil, nil, err
	// }
	// validators = staking.WriteValidators(ctx, app.stakingKeeper)
	return appState, validators, nil
}

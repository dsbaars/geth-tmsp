package app

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"sync"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc/api"
	"github.com/ethereum/go-ethereum/rpc/codec"
	"github.com/ethereum/go-ethereum/rpc/comms"
	"github.com/ethereum/go-ethereum/rpc/shared"
	"github.com/ethereum/go-ethereum/rpc/useragent"
	"github.com/ethereum/go-ethereum/xeth"

	"github.com/codegangsta/cli"
	client "github.com/tendermint/go-rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tmsp/types"
)

const (
	ClientIdentifier = "Geth-TMSP"
	Version          = "1.3.3"
	nodeNameVersion  = Version // TODO: pass git commit here ...
)

// eth.Ethereum is a monolith handling everything from peers to compiling solidity.
// We need it for sanity, but we need to manage state independently.
// The Xeth is modified to use the latest committed state only, and all block related functionality is removed (at least for now).
type EthereumApplication struct {
	// all hail the one holy object
	ethereum *eth.Ethereum

	// committed and working states
	mtx       sync.Mutex
	stateRoot []byte         // last committed state root
	stateDB   *state.StateDB // for applying txs

	// client to the tendermint core rpc
	client *client.ClientURI

	// NOTE: we need to keep another state for check tx (txpool).
	// (for now we KISS so don't expect the rpc to be smart or robust)
}

func NewEthereumApplication(ctx *cli.Context) *EthereumApplication {
	cfg := utils.MakeEthConfig(ClientIdentifier, nodeNameVersion, ctx)
	ethereum, err := eth.New(cfg)
	if err != nil {
		utils.Fatalf("%v", err)
	}

	stateDB, err := ethereum.BlockChain().State()
	if err != nil {
		utils.Fatalf("%v", err)
	}

	ethApp := &EthereumApplication{
		ethereum: ethereum,
		stateDB:  stateDB,
		client:   client.NewClientURI(fmt.Sprintf("http://%s", ctx.String(TendermintCoreHostFlag.Name))),
	}

	// NOTE: RPC/IPC should only be enabled on local nodes
	if !ctx.GlobalBool(utils.IPCDisabledFlag.Name) {
		if err := ethApp.StartIPC(ctx); err != nil {
			utils.Fatalf("%v", err)
		}
	}

	if ctx.GlobalBool(utils.RPCEnabledFlag.Name) {
		if err := ethApp.StartRPC(ctx); err != nil {
			utils.Fatalf("%v", err)
		}
	}
	return ethApp
}

// Start the local RPC server
func (app *EthereumApplication) StartRPC(ctx *cli.Context) error {
	config := comms.HttpConfig{
		ListenAddress: ctx.GlobalString(utils.RPCListenAddrFlag.Name),
		ListenPort:    uint(ctx.GlobalInt(utils.RPCPortFlag.Name)),
		CorsDomain:    ctx.GlobalString(utils.RPCCORSDomainFlag.Name),
	}

	eth := app.ethereum
	xeth := xeth.NewFromApp(app, nil)
	codec := codec.JSON

	// We only run a few of the APIs.
	// For now, all block related functionality will fail (ie. only state!)
	// Eventually, replace with calls to tendermint core
	ethApi := api.NewEthApi(xeth, eth, codec)
	personalApi := api.NewPersonalApi(xeth, eth, codec)
	web3Api := api.NewWeb3Api(xeth, codec)
	return comms.StartHttp(config, codec, api.Merge(ethApi, personalApi, web3Api))
}

func (app *EthereumApplication) StartIPC(ctx *cli.Context) error {
	eth := app.ethereum
	config := comms.IpcConfig{
		Endpoint: utils.IpcSocketPath(ctx),
	}

	initializer := func(conn net.Conn) (comms.Stopper, shared.EthereumApi, error) {
		fe := useragent.NewRemoteFrontend(conn, eth.AccountManager())
		xeth := xeth.NewFromApp(app, fe)
		ethApi := api.NewEthApi(xeth, eth, codec.JSON)
		personalApi := api.NewPersonalApi(xeth, eth, codec.JSON)
		web3Api := api.NewWeb3Api(xeth, codec.JSON)
		return xeth, api.Merge(ethApi, personalApi, web3Api), nil
	}

	return comms.StartIpc(config, codec.JSON, initializer)
}

//--------------------------------------------------------
// Implements xeth.ethApp interface
// Allows state to be managed by the EthereumApplication, and txs to be broadcast to the tendermint core
// eventually we need to redirect all chain/tx queries there too

func (app *EthereumApplication) Ethereum() *eth.Ethereum {
	return app.ethereum
}

func (app *EthereumApplication) StateRoot() common.Hash {
	app.mtx.Lock()
	defer app.mtx.Unlock()
	return common.BytesToHash(app.stateRoot)
}

func (app *EthereumApplication) BroadcastTx(tx *ethtypes.Transaction) error {
	var result ctypes.TMResult
	buf := new(bytes.Buffer)
	if err := tx.EncodeRLP(buf); err != nil {
		return err
	}
	params := map[string]interface{}{
		"tx": hex.EncodeToString(buf.Bytes()),
	}
	_, err := app.client.Call("broadcast_tx", params, &result)
	return err
}

// Equivalent of the tx-pool state; updated by CheckTx and replaced on commit
// For now, we just return the last block state
func (app *EthereumApplication) ManagedState() *state.ManagedState {
	st, err := state.New(app.StateRoot(), app.ethereum.ChainDb())
	if err != nil {
		return nil // (?!)
	}
	return state.ManageState(st)
}

//--------------------------------------------------------
// Implements TMSP App

func (app *EthereumApplication) Info() string {
	// TODO: cache some info about the state
	return ""
}

func (app *EthereumApplication) Query(query []byte) (result []byte, log string) {
	return nil, ""
}

func (app *EthereumApplication) SetOption(key string, value string) string {
	// TODO: gas limits, other params
	return ""
}

func (app *EthereumApplication) AppendTx(txBytes []byte) (retCode types.RetCode, result []byte, log string) {
	// decode and run tx
	tx := new(ethtypes.Transaction)
	rlpStream := rlp.NewStream(bytes.NewBuffer(txBytes), 0)
	if err := tx.DecodeRLP(rlpStream); err != nil {
		return types.RetCodeEncodingError, result, log
	}

	gpi := big.NewInt(1000000000) // a billion ... TODO: configurable
	gp := core.GasPool(*gpi)      // XXX: this feels so wrong!?
	ret, gas, err := core.ApplyMessage(NewEnv(app.stateDB, tx), tx, &gp)
	if err != nil {
		if err == ethtypes.ErrInvalidSig || err == ethtypes.ErrInvalidPubKey {
			return types.RetCodeUnauthorized, result, err.Error()
		} else if core.IsNonceErr(err) {
			return types.RetCodeBadNonce, result, err.Error()
		} else if core.IsInvalidTxErr(err) {
			return types.RetCodeInsufficientFees, result, err.Error() // bad gas or value transfer
		} else {
			return types.RetCodeUnauthorized, result, err.Error() // bad pubkey recovery
		}
	}
	_, _ = ret, gas
	return types.RetCodeOK, result, log
}

func (app *EthereumApplication) CheckTx(tx []byte) (retCode types.RetCode, result []byte, log string) {
	// TODO: check nonce, balance, etc
	return retCode, result, log
}

// Commit
func (app *EthereumApplication) GetHash() (hashBytes []byte, log string) {
	hash, err := app.stateDB.Commit()
	if err != nil {
		return nil, err.Error()
	}

	app.mtx.Lock()
	app.stateRoot = hash.Bytes()
	app.mtx.Unlock()

	return hash.Bytes(), log
}

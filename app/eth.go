package app

import (
	"bytes"
	"math/big"
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
	"github.com/ethereum/go-ethereum/xeth"

	"github.com/codegangsta/cli"
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
	}

	// NOTE: RPC should only be enabled on local nodes
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

	// XXX: we need to give it the eth app so it can get latest state hash,
	// instead of from the ethereum obj
	eth := app.ethereum
	xeth := xeth.New(app, eth, nil)
	codec := codec.JSON

	// We only run a few of the APIs.
	// For now, all block related functionality removed (ie. only state!)
	// Eventually, replace with calls to tendermint core
	ethApi := api.NewEthApi(xeth, eth, codec)
	personalApi := api.NewPersonalApi(xeth, eth, codec)
	web3Api := api.NewWeb3Api(xeth, codec)
	return comms.StartHttp(config, codec, api.Merge(ethApi, personalApi, web3Api))
}

// So XEth can load state with the latest state root
func (app *EthereumApplication) StateRoot() common.Hash {
	app.mtx.Lock()
	defer app.mtx.Unlock()
	return common.BytesToHash(app.stateRoot)
}

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

	gpi := big.NewInt(0)
	gp := core.GasPool(*gpi) // XXX: this feels so wrong!?
	ret, gas, err := core.ApplyMessage(NewEnv(app.stateDB, tx), tx, &gp)
	if err != nil {
		if err == ethtypes.ErrInvalidSig || err == ethtypes.ErrInvalidPubKey {
			return types.RetCodeUnauthorized, result, log
		} else if core.IsNonceErr(err) {
			return types.RetCodeBadNonce, result, log
		} else if core.IsInvalidTxErr(err) {
			return types.RetCodeInsufficientFees, result, log // bad gas or value transfer
		} else {
			return types.RetCodeUnauthorized, result, log // bad pubkey recovery
		}
	}
	_, _ = ret, gas
	return types.RetCodeOK, result, log
}

func (app *EthereumApplication) CheckTx(tx []byte) (retCode types.RetCode, result []byte, log string) {
	// TODO: check nonce, balance, etc
	return retCode, result, log
}

// NOTE: the only way to get the hash is to commit (?!)
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

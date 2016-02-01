package xeth

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/miner"
)

// Same as New but using the tmsp ethApp
func NewFromApp(ethApp ethStateApp, frontend Frontend) *XEth {
	ethereum := ethApp.Ethereum()
	xeth := &XEth{
		ethApp:           ethApp, // so we can get state roots
		backend:          ethereum,
		frontend:         frontend,
		quit:             make(chan struct{}),
		filterManager:    filters.NewFilterSystem(ethereum.EventMux()),
		logQueue:         make(map[int]*logQueue),
		blockQueue:       make(map[int]*hashQueue),
		transactionQueue: make(map[int]*hashQueue),
		messages:         make(map[int]*whisperFilter),
		agent:            miner.NewRemoteAgent(),
		gpo:              eth.NewGasPriceOracle(ethereum),
	}
	if ethereum.Whisper() != nil {
		xeth.whisper = NewWhisper(ethereum.Whisper())
	}
	ethereum.Miner().Register(xeth.agent)
	if frontend == nil {
		xeth.frontend = dummyFrontend{}
	}
	state, _ := xeth.backend.BlockChain().State()
	xeth.state = NewState(xeth, state)
	go xeth.start()
	return xeth
}

//----------
// functionality needed so xeth can work with the tmsp ethereum app and eg talk to tendermint core
type ethStateApp interface {
	Ethereum() *eth.Ethereum
	StateRoot() common.Hash
	BroadcastTx(tx *types.Transaction) error
	ManagedState() *state.ManagedState
}

func (self *XEth) StateRoot() common.Hash {
	return self.ethApp.StateRoot()
}

func (self *XEth) LatestState() *XEth {
	st, err := state.New(self.StateRoot(), self.backend.ChainDb())
	if err != nil {
		return nil
	}
	return self.WithState(st)
}

func (self *XEth) BroadcastTx(tx *types.Transaction) error {
	return self.ethApp.BroadcastTx(tx)
}

func (self *XEth) ManagedState() *state.ManagedState {
	return self.ethApp.ManagedState()
}

//----------

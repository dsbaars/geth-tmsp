# geth-tmsp
Ethereum's geth over the TenderMint Socket Protocol

Geth rpc/ipc interfaces are provided as normal, and whatever flags can be supported are.
Some backend functions are torn out to talk to tendermint core.

Currently most of the chain related txs will just work with the genesis block, until we route those functions properly to tendermint.

Working with the current state is supported.

# Install

We need a slightly modified version of go-ethereum, and something is up with godeps (as usual):

```
go get github.com/ethereum/go-ethereum
cd $GOPATH/src/github.com/ethereum/go-ethereum
git remote add eris https://github.com/eris-ltd/go-ethereum
git fetch -a eris
git checkout eris/tmsp
make geth
```

Now we can install the `geth-tmsp` app:

```
go get github.com/eris-ltd/geth-tmsp
```

Run the test (`test/test.sh`) for an example. It requires tendermint be installed as well

TODO: dockerize the test, test contracts
TODO: mempool state, wiring up blocks, etc.


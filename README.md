# geth-tmsp
Ethereum's geth over the [Tendermint Socket Protocol](http://tendermint.com/tutorials/run-your-first-tmsp-application/)

Geth rpc/ipc interfaces are provided as normal, and whatever flags can be supported are.
Some backend functions are torn out to talk to tendermint core.

Currently most of the chain related txs will just work with the genesis block, until we route those functions properly to tendermint.

Working with the current state is supported.

A [fork of go-ethereum](https://github.com/eris-ltd/go-ethereum/tree/tmsp) is used with changes only to the `xeth` package to support calling the tendermint daemon for broadcasting transactions (and eventually getting block info). This code is vendored under ethereum/go-ethereum because golang has tyranical import rules.

# Install

```
go get github.com/eris-ltd/geth-tmsp
```

Normal `geth` can be used to `attach` to a running `geth-tmsp` to interact with the chain.

See the test in test/test.sh for more details. 

Run the test in a docker container with:

```
docker build -t eris/geth-tmsp-test -f test/Dockerfile .
docker run --rm -t eris/geth-tmsp-test
```

# TODO

- test contracts, mempool state, wiring up blocks, etc.


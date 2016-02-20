#! /bin/bash

export DATADIR=~/geth-tmsp-test
export TMROOT=~/geth-tmsp-tendermint

# geth prefix
geth() {
	$GOPATH/src/github.com/ethereum/go-ethereum/build/bin/geth --datadir $DATADIR "$@"
}
export -f geth

removeQuotes() {
	a=$1
	a="${a%\"}" 
	a="${a#\"}"
	echo $a
}

export -f removeQuotes

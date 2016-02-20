#! /bin/bash

export DATADIR=~/geth-tmsp-test
export TMROOT=~/geth-tmsp-tendermint

mkdir -p $DATADIR
mkdir -p $TMROOT

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

# create a new key 
KEYSTRING=`geth --password ./test/passwd account new` 

# get our address and make a genesis file with it
ADDR=`echo $KEYSTRING | awk '{print $2}' | sed 's/[{}]//g'` 
#ADDR=`geth --exec 'eth.accounts[0]' console 2>/dev/null`
#ADDR=`removeQuotes $ADDR`
ethgen $ADDR > genesis.json

cat genesis.json

# start the geth-tmsp app
geth-tmsp --datadir $DATADIR --genesis genesis.json --verbosity 6 > geth-tmsp.log 2>&1 &
sleep 4

# get the state hash
STATE=`geth --exec "web3.eth.getBlock(0).stateRoot" attach`
STATE=`removeQuotes $STATE`
STATE=${STATE:2} # drop the leading 0x
echo $STATE

# initalize tendermint files
# and set the state root in the genesis
tendermint init
cat $TMROOT/genesis.json | jq .app_hash=\"$STATE\" > mintgenesis.json 
mv mintgenesis.json $TMROOT/genesis.json
cat $TMROOT/genesis.json

# start tendermint node 
tendermint node --log_level=debug --fast_sync=false --skip_upnp > tendermint.log 2>&1 &
sleep 3

# run the eth test script
geth --exec 'loadScript("test/script1.js")' attach


sleep 2 # commit blocks TODO: sleep in script

RESULT=`geth --exec 'loadScript("test/script2.js")' attach`
echo $RESULT

if [[ "$RESULT" != "1" ]]; then
	exit 1
fi

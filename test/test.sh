#! /bin/bash

source test/util.sh

mkdir -p $DATADIR
mkdir -p $TMROOT

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
RESULT=`geth --exec 'loadScript("test/script.js")' attach`
RESULT=`echo $RESULT | awk '{print $1}'` # for some reason geth --exec always outputs true when it finishes

if [[ "$RESULT" != "1" ]]; then
	echo "Failed"
	exit 1
fi

echo "Passed"

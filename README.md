# geth-tmsp
Ethereum's geth over the TenderMint Socket Protocol

Geth rpc/ipc interfaces are provided as normal, and whatever flags can be supported are.
Some backend functions are torn out to talk to tendermint core.



Details
-------

We use an Ethereum (backend) object because its a monolith and holds everything, but we manage the state in a seperate stateDB.
Any reference to the state must be to this stateDB.
Any reference to the chain should be routed to tendermint core.





What if 

use ErisDB as local proxy
makes calls to eth rpc and to tendermint core

eth rpc funcs stay as is .. 



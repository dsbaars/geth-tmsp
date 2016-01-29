package app

import (
	"github.com/codegangsta/cli"
)

var (
	// TODO: need to get this another way for docker since core dials/links to the app.
	TendermintCoreHostFlag = cli.StringFlag{
		Name:  "mintcore",
		Usage: "RPC address of tendermint core for this instance of the application state",
		Value: "localhost:46657",
	}
)

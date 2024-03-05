//go:build rev2

package main

import (
	"context"
	"log"
	"os"

	"github.com/koron/funddb/internal/subcmd"
	"github.com/koron/funddb/subcmds/fund"
)

var commandSet = subcmd.DefineRootSet(fund.Set)

func main() {
	err := subcmd.Run(context.Background(), commandSet, os.Args[1:]...)
	if err != nil {
		log.Fatal(err)
	}
}

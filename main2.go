//go:build rev2

package main

import (
	"context"
	"log"
	"os"

	"github.com/koron/funddb/internal/subcmd"
	"github.com/koron/funddb/subcmds/fund"
	"github.com/koron/funddb/subcmds/price"
)

var commandSet = subcmd.DefineRootSet(fund.Set, price.Set)

func main() {
	err := subcmd.Run(context.Background(), commandSet, os.Args[1:]...)
	if err != nil {
		log.Fatal(err)
	}
}

//go:build rev2

package main

import (
	"context"
	"log"
	"os"

	"github.com/koron/funddb/internal/subcmd"
	"github.com/koron/funddb/subcmds/database"
	"github.com/koron/funddb/subcmds/fund"
	"github.com/koron/funddb/subcmds/price"
)

var commandSet = subcmd.DefineRootSet(price.Set, fund.Set, database.Set)

func main() {
	err := subcmd.Run(context.Background(), commandSet, os.Args[1:]...)
	if err != nil {
		log.Fatal(err)
	}
}

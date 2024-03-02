package main

import (
	"context"
	"log"
	"os"

	"github.com/koron/funddb/internal/subcmd"
	"github.com/koron/funddb/subcmds/fund"
)

var commandSet = subcmd.DefineSet(subcmd.RootName(), "", fund.Set)

func main2() {
	err := subcmd.Run(context.Background(), commandSet, os.Args[1:]...)
	//dbfile := flag.String("d", "fund2.db", `database file`)
	//flag.Parse()
	//_, err := orm.NewEngine(*dbfile, false)
	if err != nil {
		log.Fatal(err)
	}
}

package appcore

import (
	"context"
	"flag"
	"strings"

	"github.com/koron/funddb/internal/dataobj"
	"github.com/koron/funddb/internal/subcmd"
	"xorm.io/xorm"
)

type Core struct {
	ORM *xorm.Engine
}

type FlagHook func(fs *flag.FlagSet)

func New(ctx context.Context, args []string, flagHooks ...FlagHook) (ac *Core, flagArgs []string, err error) {
	name := strings.Join(subcmd.Names(ctx), " ")
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	dbfile := fs.String("dbfile", "fund2.db", "database file")
	for _, hook := range flagHooks {
		hook(fs)
	}
	fs.Parse(args)
	orm, err := dataobj.NewEngine(*dbfile, false)
	if err != nil {
		return nil, nil, err
	}
	return &Core{ORM: orm}, fs.Args(), nil
}

func (ac *Core) Close() error {
	return ac.ORM.Close()
}

package appcore

import (
	"context"
	"flag"
	"strings"

	"github.com/koron-go/subcmd"
	"github.com/koron/funddb/internal/dataobj"
	"xorm.io/xorm"
)

type Core struct {
	ORM     *xorm.Engine
	ShowSQL bool
}

type FlagHook func(fs *flag.FlagSet)

func New(ctx context.Context, args []string, flagHooks ...FlagHook) (ac *Core, flagArgs []string, err error) {
	name := strings.Join(subcmd.Names(ctx), " ")
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	dbfile := fs.String("dbfile", "fund.db", "database file")
	showsql := fs.Bool("showsql", false, "show SQL for debug")
	for _, hook := range flagHooks {
		hook(fs)
	}
	fs.Parse(args)
	orm, err := dataobj.NewEngine(*dbfile)
	if err != nil {
		return nil, nil, err
	}
	if *showsql {
		orm.ShowSQL(true)
	}
	return &Core{
		ORM:     orm,
		ShowSQL: *showsql,
	}, fs.Args(), nil
}

func (ac *Core) Close() error {
	return ac.ORM.Close()
}

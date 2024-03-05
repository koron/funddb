package fund

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/koron/funddb/internal/orm"
	"github.com/koron/funddb/internal/subcmd"
	"xorm.io/xorm"
)

func bootstrapORM(ctx context.Context, args []string, hooks ...func(fs *flag.FlagSet)) (*xorm.Engine, []string, error) {
	name := strings.Join(subcmd.Names(ctx), " ")
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	dbfile := fs.String("dbfile", "fund2.db", "database file")
	for _, hook := range hooks {
		hook(fs)
	}
	fs.Parse(args)
	engine, err := orm.NewEngine(*dbfile, false)
	if err != nil {
		return nil, nil, err
	}
	return engine, fs.Args(), nil
}

var Import = subcmd.DefineCommand("import", "import funds from file", func(ctx context.Context, args []string) error {
	engine, _, err := bootstrapORM(ctx, args)
	if err != nil {
		return err
	}
	// TODO:
	_ = engine
	return nil
})

var List = subcmd.DefineCommand("list", "list funds", func(ctx context.Context, args []string) error {
	engine, _, err := bootstrapORM(ctx, args)
	if err != nil {
		return err
	}
	defer engine.Close()

	return engine.Iterate(&orm.Fund{}, func(idx int, bean interface{}) error {
		f := bean.(*orm.Fund)
		fmt.Printf("%+v\n", *f)
		return nil
	})
})

var Add = subcmd.DefineCommand("add", "add a fund", func(ctx context.Context, args []string) error {
	// TODO:
	return nil
})

var Delete = subcmd.DefineCommand("delete", "delete a fund", func(ctx context.Context, args []string) error {
	engine, params, err := bootstrapORM(ctx, args)
	if err != nil {
		return err
	}
	defer engine.Close()

	if len(params) == 0 {
		return errors.New("required one or more ID of fund to be deleted")
	}

	session := engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	for _, id := range params {
		n, err := session.Delete(&orm.Fund{ID: id})
		if err != nil {
			session.Rollback()
			return err
		}
		if n == 0 {
			session.Rollback()
			return fmt.Errorf("no funds for id:%s", id)
		}
	}
	return session.Commit()
})

var Modify = subcmd.DefineCommand("modify", "modify a fund", func(ctx context.Context, args []string) error {
	// TODO:
	return nil
})

var Set = subcmd.DefineSet("fund", "operate funds",
	Import,
	List,
	//Add,
	//Delete,
	//Modify,
)

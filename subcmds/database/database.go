package database

import (
	"context"

	"github.com/koron/funddb/internal/appcore"
	"github.com/koron/funddb/internal/dataobj"
	"github.com/koron/funddb/internal/subcmd"
)

var InitSchema = subcmd.DefineCommand("initschema", "Initialize schema by statements", func(ctx context.Context, args []string) error {
	ac, _, err := appcore.New(ctx, args)
	if err != nil {
		return err
	}
	defer ac.Close()
	return dataobj.InitSchema(ac.ORM, ac.ShowSQL)
})

var SyncORM = subcmd.DefineCommand("syncorm", "Sync schema with ORM", func(ctx context.Context, args []string) error {
	ac, _, err := appcore.New(ctx, args)
	if err != nil {
		return err
	}
	defer ac.Close()
	return ac.ORM.Sync(dataobj.Beans...)
})

var Set = subcmd.DefineSet("database", "operate database",
	InitSchema,
	SyncORM,
)

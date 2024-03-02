package fund

import (
	"context"

	"github.com/koron/funddb/internal/subcmd"
)

var List = subcmd.DefineCommand("list", "list funds", func(ctx context.Context, args []string) error {
	// TODO:
	return nil
})

var Add = subcmd.DefineCommand("add", "add a fund", func(ctx context.Context, args []string) error {
	// TODO:
	return nil
})

var Delete = subcmd.DefineCommand("delete", "delete a fund", func(ctx context.Context, args []string) error {
	// TODO:
	return nil
})

var Modify = subcmd.DefineCommand("modify", "modify a fund", func(ctx context.Context, args []string) error {
	// TODO:
	return nil
})

var Set = subcmd.DefineSet("fund", "operate funds",
	List,
	Add,
	Delete,
	Modify,
)

package fund

import (
	"context"
	"errors"
	"fmt"

	"github.com/koron/funddb/internal/appcore"
	"github.com/koron/funddb/internal/dataobj"
	"github.com/koron/funddb/internal/subcmd"
)

var Import = subcmd.DefineCommand("import", "import funds from file", func(ctx context.Context, args []string) error {
	ac, _, err := appcore.New(ctx, args)
	if err != nil {
		return err
	}
	// TODO:
	_ = ac
	return nil
})

var List = subcmd.DefineCommand("list", "list funds", func(ctx context.Context, args []string) error {
	ac, _, err := appcore.New(ctx, args)
	if err != nil {
		return err
	}
	defer ac.Close()

	return ac.ORM.Iterate(&dataobj.Fund{}, func(idx int, bean interface{}) error {
		f := bean.(*dataobj.Fund)
		fmt.Printf("%+v\n", *f)
		return nil
	})
})

var Add = subcmd.DefineCommand("add", "add a fund", func(ctx context.Context, args []string) error {
	// TODO:
	return nil
})

var Delete = subcmd.DefineCommand("delete", "delete a fund", func(ctx context.Context, args []string) error {
	ac, params, err := appcore.New(ctx, args)
	if err != nil {
		return err
	}
	defer ac.Close()

	if len(params) == 0 {
		return errors.New("required one or more ID of fund to be deleted")
	}

	session := ac.ORM.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	for _, id := range params {
		n, err := session.Delete(&dataobj.Fund{ID: id})
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

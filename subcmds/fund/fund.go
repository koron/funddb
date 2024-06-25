package fund

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/koron-go/subcmd"
	"github.com/koron/funddb/internal/appcore"
	"github.com/koron/funddb/internal/dataobj"
	"github.com/koron/funddb/internal/xormhelper"
	"xorm.io/xorm"
)

func importFile(ctx context.Context, session *xorm.Session, fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.Comma = '\t'
	r.Comment = '#'
	r.ReuseRecord = true
	for {
		r.FieldsPerRecord = 0
		records, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		if len(records) < 3 {
			return errors.New("few records, require 3 at least")
		}
		fund := dataobj.Fund{
			ID:   strings.TrimSpace(records[0]),
			Name: strings.TrimSpace(records[1]),
			URL:  strings.TrimSpace(records[2]),
		}
		if len(records) >= 4 {
			fund.FetchID = strings.TrimSpace(records[3])
		}
		err = xormhelper.UpsertOne(session, fund.ID, &fund)
		if err != nil {
			return err
		}
	}
	return nil
}

var Import = subcmd.DefineCommand("import", "import funds from TSV file (id, name, url, fetch_id)", func(ctx context.Context, args []string) error {
	ac, files, err := appcore.New(ctx, args)
	if err != nil {
		return err
	}
	defer ac.Close()
	if len(files) == 0 {
		return errors.New("no files to import as fund")
	}
	return xormhelper.Tx(ac.ORM, func(session *xorm.Session) error {
		for _, f := range files {
			err := importFile(ctx, session, f)
			if err != nil {
				return err
			}
		}
		return nil
	})
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

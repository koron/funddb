package price

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/k0kubun/pp/v3"
	"github.com/koron-go/subcmd"
	"github.com/koron/funddb/internal/appcore"
	"github.com/koron/funddb/internal/dataobj"
	"github.com/koron/funddb/internal/xormhelper"
	"xorm.io/xorm"
)

var FetchTest = subcmd.DefineCommand("fetchtest", "test: fetch price data and print", func(ctx context.Context, args []string) error {
	ac, ids, err := appcore.New(ctx, args, func(fs *flag.FlagSet) {})
	if err != nil {
		return err
	}
	defer ac.Close()
	if len(ids) == 0 {
		return errors.New("require one or more fund IDs")
	}
	return xormhelper.Tx(ac.ORM, func(session *xorm.Session) error {
		ctx := context.Background()
		for _, id := range ids {
			var fund dataobj.Fund
			has, err := session.ID(id).Get(&fund)
			if err != nil {
				return err
			}
			if !has {
				return fmt.Errorf("no funds found for ID=%s", id)
			}
			p, err := fetchPrice(ctx, fund.FetchID)
			if err != nil {
				return err
			}
			pp.Print(p)
		}
		return nil
	})
})

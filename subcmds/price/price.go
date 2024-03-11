package price

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/koron-go/subcmd"
	"github.com/koron/funddb/internal/ammufg"
	"github.com/koron/funddb/internal/appcore"
	"github.com/koron/funddb/internal/dataobj"
	"github.com/koron/funddb/internal/fidelity"
	"github.com/koron/funddb/internal/xormhelper"
	"xorm.io/xorm"
	"xorm.io/xorm/schemas"
)

func fetchPrice(ctx context.Context, fetchID string) (*dataobj.Price, error) {
	parts := strings.SplitN(fetchID, ":", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid fetch ID, required format \"{scheme}:{id}\": %s", fetchID)
	}
	scheme, id := parts[0], parts[1]
	switch scheme {
	case "fidelity":
		d, err := fidelity.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		return &dataobj.Price{
			Date:  dataobj.DateFromTime(d.Date()),
			Value: int(d.Price()),
		}, nil

	case "ammufg":
		d, err := ammufg.Get(ctx, ammufg.CodeTypeFund, id)
		if err != nil {
			return nil, err
		}
		return &dataobj.Price{
			Date:  dataobj.DateFromTime(d.Date()),
			Value: int(d.Price()),
		}, nil

	default:
		return nil, fmt.Errorf("unknown scheme: %s", scheme)
	}
}

func upsertPrice(session *xorm.Session, p *dataobj.Price) error {
	var curr dataobj.Price
	ok, err := session.Where("id = ? AND date = ?", p.ID, p.Date).Get(&curr)
	if err != nil {
		return err
	}
	log.Printf("upsertPrice: p=%+v ok=%t curr=%+v", p, ok, curr)
	if ok {
		if curr.Value == p.Value {
			log.Printf("skip %+v, not updated", p)
			return nil
		}
		// update only value
		updated, err := session.Where("id = ? AND date = ?", p.ID, p.Date).Update(&dataobj.Price{Value: p.Value})
		if err != nil {
			return err
		}
		if updated != 1 {
			return fmt.Errorf("not 1 row updated on prices: %d", updated)
		}
		return nil
	}
	// inset new value
	inserted, err := session.Insert(p)
	if err != nil {
		return err
	}
	if inserted != 1 {
		return fmt.Errorf("expected 1 row inserted on prices, but %d rows inserted", inserted)
	}
	return nil
}

var FetchLatest = subcmd.DefineCommand("fetchlatest", "fetch latest price data", func(ctx context.Context, args []string) error {
	ac, filter, err := appcore.New(ctx, args)
	if err != nil {
		return err
	}
	defer ac.Close()
	const batchSize = 100
	var ids []any
	if len(filter) > 0 {
		ids = make([]any, len(filter))
		for i, f := range filter {
			ids[i] = f
		}
	}
	return xormhelper.Tx(ac.ORM, func(session *xorm.Session) error {
		// count target funds.
		if len(ids) > 0 {
			session.In("id", ids...)
		}
		fundCnt, err := session.Count(&dataobj.Fund{})
		if err != nil {
			return err
		}
		if fundCnt == 0 {
			return nil
		}
		// fetch latest price with batches.
		for i := 0; i < int(fundCnt); i += batchSize {
			session.OrderBy("id").Limit(batchSize, i)
			if len(ids) > 0 {
				session.In("id", ids...)
			}
			var fundList []dataobj.Fund
			if err := session.Find(&fundList); err != nil {
				return err
			}
			for _, fund := range fundList {
				if fund.FetchID == "" {
					continue
				}
				log.Printf("fetch latest price for %s", fund.FetchID)
				p, err := fetchPrice(ctx, fund.FetchID)
				if err != nil {
					return err
				}
				p.ID = fund.ID
				pk := schemas.PK{p.ID, p.Date}
				if err := xormhelper.UpsertOne(session, pk, p); err != nil {
					return err
				}
			}
		}
		return nil
	})
})

var Set = subcmd.DefineSet("price", "operate prices",
	FetchLatest,
)

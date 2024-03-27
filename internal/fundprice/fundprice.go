package fundprice

import "time"

type Price interface {
	Date() time.Time
	Price() int64
}

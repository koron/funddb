package orm

import "time"

type Fund struct {
	ID      string `xorm:"pk"`             // Association ID
	Name    string `xorm:"notnull unique"` // Display name
	URL     string `xorm:"notnull unique"` // URL for the fund
	FetchID string `xorm:"null unique"`    // Fetch ID
}

func (Fund) TableName() string {
	return "funds"
}

type Price struct {
	ID    string    `xorm:"notnull index unique(id_date)"` // FK:Fund.ID
	Date  time.Time `xorm:"notnull index unique(id_date)"`
	Value int       `xorm:"bigint not null"`
}

func (Price) TableName() string {
	return "prices"
}

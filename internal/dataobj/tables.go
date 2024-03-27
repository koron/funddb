package dataobj

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
	ID    string `xorm:"notnull index unique(id_date) pk"` // FK:Fund.ID
	Date  Date   `xorm:"notnull index unique(id_date) pk"`
	Value int64  `xorm:"bigint not null"`

	NetAssets int64 `xorm:"bigint null"`
}

func (Price) TableName() string {
	return "prices"
}

var Beans = []any{&Fund{}, &Price{}}

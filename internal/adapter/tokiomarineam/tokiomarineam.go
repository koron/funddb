package tokiomarineam

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Should implement fundprice.Price
type FundInfo struct {
	fundID string

	Dt    string      `json:"Dt"`
	Nav   int64       `json:"Nav"`
	Diff  int64       `json:"Diff"`
	Rtn   json.Number `json:"Rtn"`
	Asset int64       `json:"Asset"`

	Div1   int64  `json:"Div1"`
	DivDt1 string `json:"DivDt1"`
	Div2   int64  `json:"Div2"`
	DivDt2 string `json:"DivDt2"`
	Div3   int64  `json:"Div3"`
	DivDt3 string `json:"DivDt3"`

	ReportDt  string `json:"ReportDt"`
	ReportDt1 string `json:"ReportDt1"`
	ReportDt2 string `json:"ReportDt2"`

	Hansya []Hansya `json:"Hansya"`
}

func (f FundInfo) ID() string {
	return f.fundID
}

func (f FundInfo) Date() time.Time {
	loc, err := time.LoadLocation("Japan")
	if err != nil {
		log.Print(err)
		return time.Time{}
	}
	ti, err := time.ParseInLocation("2006/01/02", f.Dt, loc)
	if err != nil {
		log.Printf("invalid date %s: %s", f.Dt, err)
		return time.Time{}
	}
	return ti.Add(time.Hour * 18)
}

func (f FundInfo) Price() int64 {
	return f.Nav
}

func (f FundInfo) NetAssets() int64 {
	return f.Asset
}

type Hansya struct {
	ToriKbn string `json:"ToriKbn"`
	Type    string `json:"Type"`
	Url1    string `json:"Url1"`
	Url2    string `json:"Url2"`
	NmKana  string `json:"NmKana"`
	Nm      string `json:"Nm"`
}

func Get(ctx context.Context, fundId string, dummy *string) (*FundInfo, error) {
	// Compose and make a GET request.
	u := fmt.Sprintf("https://api.tokiomarineam.co.jp/hp/funds?FundId=%s", url.QueryEscape(fundId))
	if dummy != nil {
		u += "&_=" + url.QueryEscape(*dummy)
	}
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		all, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("[WARN] failed to parse error:%d response: %v", res.StatusCode, err)
		} else {
			log.Printf("[INFO] error:%d reponse:\n%s", res.StatusCode, string(all))
		}
		return nil, fmt.Errorf("failed HTTP with %d for: %q", res.StatusCode, u)
	}

	var data FundInfo
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	data.fundID = fundId

	return &data, nil
}

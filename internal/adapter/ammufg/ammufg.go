package ammufg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Number string

type Dataset struct {
	FundCD            string `json:"fund_cd"`
	AssociationFundCD string `json:"association_fund_cd"`
	ISINCd            string `json:"isin_cd"`

	FundName string `json:"fund_name"`

	BaseDate          string `json:"base_date"`
	CancellationPrice int64  `json:"cancellation_price"`

	NetAssets_                int64       `json:"netassets"`
	NetAssetsChangeCmpPrevDay json.Number `json:"netassets_change_cmp_prev_day"`

	Nav int64 `json:"nav"`

	NavMax1m     Number `json:"nav_max_1m"`
	NavMax1mDt   string `json:"nav_max_1m_dt"`
	NavMax3m     Number `json:"nav_max_3m"`
	NavMax3mDt   string `json:"nav_max_3m_dt"`
	NavMax6m     Number `json:"nav_max_6m"`
	NavMax6mDt   string `json:"nav_max_6m_dt"`
	NavMax1y     Number `json:"nav_max_1y"`
	NavMax1yDt   string `json:"nav_max_1y_dt"`
	NavMaxFull   Number `json:"nav_max_full"`
	NavMaxFullDt string `json:"nav_max_full_dt"`

	NavMin1m     Number `json:"nav_min_1m"`
	NavMin1mDt   string `json:"nav_min_1m_dt"`
	NavMin3m     Number `json:"nav_min_3m"`
	NavMin3mDt   string `json:"nav_min_3m_dt"`
	NavMin6m     Number `json:"nav_min_6m"`
	NavMin6mDt   string `json:"nav_min_6m_dt"`
	NavMin1y     Number `json:"nav_min_1y"`
	NavMin1yDt   string `json:"nav_min_1y_dt"`
	NavMinFull   Number `json:"nav_min_full"`
	NavMinFullDt string `json:"nav_min_full_dt"`

	PercentageChange Number `json:"percentage_change"`

	PercentageChange1m   Number `json:"percentage_change_1m"`
	PercentageChange3m   Number `json:"percentage_change_3m"`
	PercentageChange6m   Number `json:"percentage_change_6m"`
	PercentageChange1y   Number `json:"percentage_change_1y"`
	PercentageChangeFull Number `json:"percentage_change_full"`

	PercentageChangeMax1m   Number `json:"percentage_change_max_1m"`
	PercentageChangeMax3m   Number `json:"percentage_change_max_3m"`
	PercentageChangeMax6m   Number `json:"percentage_change_max_6m"`
	PercentageChangeMax1y   Number `json:"percentage_change_max_1y"`
	PercentageChangeMaxFull Number `json:"percentage_change_max_full"`

	PercentageChangeMin1m   Number `json:"percentage_change_min_1m"`
	PercentageChangeMin3m   Number `json:"percentage_change_min_3m"`
	PercentageChangeMin6m   Number `json:"percentage_change_min_6m"`
	PercentageChangeMin1y   Number `json:"percentage_change_min_1y"`
	PercentageChangeMinFull Number `json:"percentage_change_min_full"`

	Risk1y   Number `json:"risk_1y"`
	Risk3y   Number `json:"risk_3y"`
	RiskFull Number `json:"risk_full"`

	RiskReturn1y   Number `json:"risk_return_1y"`
	RiskReturn3y   Number `json:"risk_return_3y"`
	RiskReturnFull Number `json:"risk_return_full"`
}

func (ds Dataset) Scheme() string {
	return "ammufg"
}

func (ds Dataset) ID() string {
	return ds.FundCD
}

func (ds Dataset) Date() time.Time {
	loc, err := time.LoadLocation("Japan")
	if err != nil {
		log.Print(err)
		return time.Time{}
	}
	ti, err := time.ParseInLocation("20060102", ds.BaseDate, loc)
	if err != nil {
		log.Printf("invalid date %s: %s", ds.BaseDate, err)
		return time.Time{}
	}
	return ti.Add(time.Hour * 18)
}

func (ds Dataset) Price() int64 {
	return ds.CancellationPrice
}

func (ds Dataset) NetAssets() int64 {
	return ds.NetAssets_
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (err Error) Error() string {
	return fmt.Sprintf("%s (code: %s)", err.Message, err.Code)
}

type Errors struct {
	Count     int     `json:"count"`
	ErrorList []Error `json:"error_list"`
}

type Result struct {
	ErrCD    string `json:"errcd"`
	ErrMsg   string `json:"errmsg"`
	Function string `json:"function"`
	RetCount int    `json:"retcount"`
	Status   int    `json:"status"`
}

type FundInfo struct {
	Result   Result    `json:"result"`
	Errors   Errors    `json:"errors"`
	Datasets []Dataset `json:"datasets"`
}

type errorResult struct {
	Result Result
	Errors Errors
}

func (er errorResult) Error() string {
	if er.Result.ErrMsg != "" {
		return fmt.Sprintf("%s (code: %s)", er.Result.ErrMsg, er.Result.ErrCD)
	}
	if len(er.Errors.ErrorList) > 0 {
		return er.Error()
	}
	return fmt.Sprintf("failed something status:%d retcount:%d errors.count:%d", er.Result.Status, er.Result.RetCount, er.Errors.Count)
}

type CodeType string

const (
	CodeTypeAssociationFund CodeType = "association_fund_cd"
	CodeTypeISIN            CodeType = "isin_cd"
	CodeTypeFund            CodeType = "fund_cd"
)

// Get retrives latest fund information (Dataset) by code and its type.
func Get(ctx context.Context, ct CodeType, code string) (*Dataset, error) {
	u := fmt.Sprintf("https://developer.am.mufg.jp/fund_information_latest/%s/%s", ct, code)
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
	// more error handlings.
	if data.Result.Status != 200 || data.Errors.Count != 0 {
		return nil, &errorResult{Result: data.Result, Errors: data.Errors}
	}
	if len(data.Datasets) < 1 {
		return nil, errors.New("no datasets available in API response")
	}
	return &data.Datasets[0], nil
}

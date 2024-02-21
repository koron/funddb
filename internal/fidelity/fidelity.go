package fidelity

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type FundData struct {
	KeyID       string
	DisplayName string
	PriceData   struct {
		ChangeAbsolute json.Number
		ChangePercent  json.Number
		Nav            struct {
			Date  string
			Value json.Number
		}
		SellingPrice json.Number
	}
}

func (fd FundData) Scheme() string {
	return "fidelity"
}

func (fd FundData) ID() string {
	return fd.KeyID
}

func (fd FundData) Date() time.Time {
	loc, err := time.LoadLocation("Japan")
	if err != nil {
		log.Print(err)
		return time.Time{}
	}
	ti, err := time.ParseInLocation(time.DateOnly, fd.PriceData.Nav.Date, loc)
	if err != nil {
		log.Printf("invalid date %s: %s", fd.PriceData.Nav.Date, err)
		return time.Time{}
	}
	return ti.Add(time.Hour*15)
}

func (fd FundData) Price() int64 {
	v, err := fd.PriceData.SellingPrice.Int64()
	if err != nil {
		log.Printf("invalid price %s: %s", fd.PriceData.SellingPrice, err)
		return -1
	}
	return v
}

func Get(ctx context.Context, id string) (*FundData, error) {
	u := fmt.Sprintf("https://www.fidelity.co.jp/api/ce/fdh/FundData.json?id=%s&country=jp", id)
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
		return nil, fmt.Errorf("HTTP failed with status code %d", res.StatusCode)
	}
	var data map[string]FundData
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	target, ok := data[id]
	if !ok {
		return nil, fmt.Errorf("fund data %s not found", id)
	}
	target.KeyID = id
	return &target, nil
}

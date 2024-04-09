package pictet

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Data struct {
	id        string
	date      time.Time
	price     int64
	netAssets int64
}

func (d Data) Scheme() string {
	return "pictet"
}

func (d Data) ID() string {
	return d.id
}

func (d Data) Date() time.Time {
	return d.date
}

func (d Data) Price() int64 {
	return d.price
}

func (d Data) NetAssets() int64 {
	return d.netAssets
}

func parseDate(s string) (time.Time, error) {
	loc, err := time.LoadLocation("Japan")
	if err != nil {
		return time.Time{}, err
	}
	ti, err := time.ParseInLocation("基準日: 2006年01月02日", s, loc)
	if err != nil {
		return time.Time{}, err
	}
	return ti.Add(time.Hour * 18), nil
}

func parsePrice(s string) (int64, error) {
	var n int64
	_, err := fmt.Sscanf(strings.ReplaceAll(s, ",", ""), "%d円", &n)
	if err != nil {
		return 0, err
	}
	return n, err
}

func parseNetAssets(s string) (int64, error) {
	var n int64
	_, err := fmt.Sscanf(strings.ReplaceAll(s, ",", ""), "%d百万円", &n)
	if err != nil {
		return 0, err
	}
	return n * 1_000_000, err
}

func Get(ctx context.Context, name string) (*Data, error) {
	u := fmt.Sprintf("https://www.pictet.co.jp/fund/%s.html", name)
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
		return nil, fmt.Errorf("failed HTTP with %d for: %q", res.StatusCode, u)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	var (
		d             Data
		errs          []error
		flagDate      bool
		flagPrice     bool
		flagNetAssets bool
	)

	doc.Find(".cmp-funds__fund-summary .cmp-fund__fund-summary-value").Each(func(i int, s *goquery.Selection) {
		key := s.Prev().Text()
		value := s.Text()
		switch key {
		case "基本情報":
			date, err := parseDate(value)
			if err != nil {
				errs = append(errs, err)
				return
			}
			flagDate = true
			d.date = date
		case "基準価額":
			p, err := parsePrice(value)
			if err != nil {
				errs = append(errs, err)
				return
			}
			flagPrice = true
			d.price = p
		case "純資産総額":
			na, err := parseNetAssets(value)
			if err != nil {
				errs = append(errs, err)
				return
			}
			flagNetAssets = true
			d.netAssets = na
		}
	})

	if !flagDate {
		errs = append(errs, errors.New("not found date"))
	}
	if !flagPrice {
		errs = append(errs, errors.New("not found price"))
	}
	if !flagNetAssets {
		errs = append(errs, errors.New("not found net assets"))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	d.id = name
	return &d, nil
}

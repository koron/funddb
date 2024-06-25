package tokiomarineam_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/koron/funddb/internal/adapter/tokiomarineam"
)

func TestDecode(t *testing.T) {
	f, err := os.Open(filepath.Join("testdata", "tokiomarineam.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var data tokiomarineam.FundInfo
	err = json.NewDecoder(f).Decode(&data)
	if err != nil {
		t.Fatal(err)
	}
	want := tokiomarineam.FundInfo{
		Dt:        "2024/06/24",
		Nav:       17203,
		Diff:      88,
		Rtn:       json.Number("0.51"),
		Asset:     34666889549,
		Div1:      500,
		DivDt1:    "2024/04/22",
		Div2:      500,
		DivDt2:    "2023/10/20",
		Div3:      500,
		DivDt3:    "2023/04/20",
		ReportDt:  "2024/06/07",
		ReportDt1: "2024/05/09",
		ReportDt2: "2024/04/05",
		Hansya: []tokiomarineam.Hansya{
			{
				ToriKbn: "0",
				Type:    "bank",
				Url1:    "https://www.a-bank.jp/",
				NmKana:  "ｱｵﾓﾘ",
				Nm:      "青森銀行",
			},
			{
				ToriKbn: "0",
				Type:    "bank",
				Url1:    "https://www.aeonbank.co.jp/",
				NmKana:  "ｲｵﾝｷﾞﾝｺｳ（ｲﾀｸｷﾝﾕｳｼｮｳﾋﾝﾄﾘﾋｷｷﾞｮｳｼｬ ﾏﾈｯｸｽｼｮｳｹﾝ)",
				Nm:      "イオン銀行（委託金融商品取引業者 マネックス証券）",
			},
			{
				ToriKbn: "0",
				Type:    "sec",
				Url1:    "https://kabu.com/",
				NmKana:  "ｴｰﾕｰｶﾌﾞｺﾑｼｮｳｹﾝ",
				Nm:      "auカブコム証券",
			},
			{
				ToriKbn: "0",
				Type:    "sec",
				Url1:    "https://www.sbisec.co.jp/ETGate",
				NmKana:  "ｴｽﾋﾞｰｱｲｼｮｳｹﾝ",
				Nm:      "ＳＢＩ証券",
			},
			{
				ToriKbn: "0",
				Type:    "bank",
				Url1:    "https://www.sbishinseibank.co.jp",
				NmKana:  "ｴｽﾋﾞｰｱｲｼﾝｾｲｷﾞﾝｺｳ（ｲﾀｸｷﾝﾕｳｼｮｳﾋﾝﾄﾘﾋｷｷﾞｮｳｼｬ ｴｽﾋﾞｰｱｲｼｮｳｹﾝ)",
				Nm:      "ＳＢＩ新生銀行（委託金融商品取引業者  ＳＢＩ証券）",
			},
			{
				ToriKbn: "0",
				Type:    "bank",
				Url1:    "https://www.sbishinseibank.co.jp",
				NmKana:  "ｴｽﾋﾞｰｱｲｼﾝｾｲｷﾞﾝｺｳ（ｲﾀｸｷﾝﾕｳｼｮｳﾋﾝﾄﾘﾋｷｷﾞｮｳｼｬ ﾏﾈｯｸｽｼｮｳｹﾝ)",
				Nm:      "ＳＢＩ新生銀行（委託金融商品取引業者 マネックス証券）",
			},
			{
				ToriKbn: "0",
				Type:    "sec",
				Url1:    "https://www.16ttsec.co.jp/",
				NmKana:  "ｼﾞｭｳﾛｸﾃｨﾃｨｼｮｳｹﾝ",
				Nm:      "十六ＴＴ証券",
			},
			{
				ToriKbn: "0",
				Type:    "sec",
				Url1:    "http://www.tokaitokyo.co.jp/",
				NmKana:  "ﾄｳｶｲﾄｳｷｮｳｼｮｳｹﾝ",
				Nm:      "東海東京証券",
			},
			{
				ToriKbn: "0",
				Type:    "bank",
				Url1:    "https://www.nagasakibank.co.jp/",
				NmKana:  "ﾅｶﾞｻｷ",
				Nm:      "長崎銀行",
			},
			{
				ToriKbn: "0",
				Type:    "bank",
				Url1:    "https://www.ncbank.co.jp/",
				NmKana:  "ﾆｼﾆｯﾎﾟﾝｼﾃｨｷﾞﾝｺｳ",
				Nm:      "西日本シティ銀行",
			},
			{
				ToriKbn: "0",
				Type:    "bank",
				Url1:    "https://www.hyakugo.co.jp/",
				NmKana:  "ﾋｬｸｺﾞ",
				Nm:      "百五銀行",
			},
			{
				ToriKbn: "0",
				Type:    "sec",
				Url1:    "http://www.105sec.co.jp/",
				NmKana:  "ﾋｬｸｺﾞｼｮｳｹﾝ",
				Nm:      "百五証券",
			},
			{
				ToriKbn: "0",
				Type:    "sec",
				Url1:    "https://www.matsui.co.jp/",
				NmKana:  "ﾏﾂｲｼｮｳｹﾝ",
				Nm:      "松井証券",
			},
			{
				ToriKbn: "0",
				Type:    "sec",
				Url1:    "https://www.monex.co.jp/",
				NmKana:  "ﾏﾈｯｸｽｼｮｳｹﾝ",
				Nm:      "マネックス証券",
			},
			{
				ToriKbn: "0",
				Type:    "bank",
				Url1:    "https://www.michinokubank.co.jp/",
				NmKana:  "ﾐﾁﾉｸｷﾞﾝｺｳ",
				Nm:      "みちのく銀行",
			},
			{
				ToriKbn: "0",
				Type:    "bank",
				Url1:    "https://www.bk.mufg.jp/",
				NmKana:  "ﾐﾂﾋﾞｼUFJｷﾞﾝｺｳ",
				Nm:      "三菱ＵＦＪ銀行",
			},
			{
				ToriKbn: "0",
				Type:    "bank",
				Url1:    "https://www.bk.mufg.jp/",
				NmKana:  "ﾐﾂﾋﾞｼUFJｷﾞﾝｺｳ(ｲﾀｸｷﾝﾕｳｼｮｳﾋﾝﾄﾘﾋｷｷﾞｮｳｼｬ ﾐﾂﾋﾞｼUFJﾓﾙｶﾞﾝ･ｽﾀﾝﾚｰｼｮｳｹﾝ)",
				Nm:      "三菱ＵＦＪ銀行（委託金融商品取引業者 三菱ＵＦＪモルガン・スタンレー証券）",
			},
			{
				ToriKbn: "0",
				Type:    "bank",
				Url1:    "https://www.tr.mufg.jp/",
				NmKana:  "ﾐﾂﾋﾞｼUFJｼﾝﾀｸｷﾞﾝｺｳ(ｷｭｳUFJ)",
				Nm:      "三菱ＵＦＪ信託銀行",
			},
			{
				ToriKbn: "0",
				Type:    "sec",
				Url1:    "https://www.sc.mufg.jp/",
				NmKana:  "ﾐﾂﾋﾞｼUFJﾓﾙｶﾞﾝ･ｽﾀﾝﾚｰｼｮｳｹﾝ(ｷｭｳUFJ)",
				Nm:      "三菱ＵＦＪモルガン・スタンレー証券",
			},
			{
				ToriKbn: "0",
				Type:    "sec",
				Url1:    "https://www.rakuten-sec.co.jp/",
				NmKana:  "ﾗｸﾃﾝｼｮｳｹﾝ",
				Nm:      "楽天証券",
			},
		},
	}
	if d := cmp.Diff(want, data, cmpopts.IgnoreUnexported(tokiomarineam.FundInfo{})); d != "" {
		t.Errorf("unmatch: --want ++got\n%s", d)
	}

	// Test methods of fundprice.Price interface.
	t.Run("Price", func(t *testing.T) {
		if got, want := data.Price(), int64(17203); got != want {
			t.Errorf("unmatch Price: want=%d got=%d", want, got)
		}
	})
	t.Run("NetAssets", func(t *testing.T) {
		if got, want := data.NetAssets(), int64(34666889549); got != want {
			t.Errorf("unmatch NetAssets: want=%d got=%d", want, got)
		}
	})
	t.Run("Date", func(t *testing.T) {
		loc, err := time.LoadLocation("Japan")
		if err != nil {
			t.Error(err)
		}
		got, want := data.Date(), time.Date(2024, 6, 24, 18, 0, 0, 0, loc)
		if d := want.Compare(got); d != 0 {
			t.Errorf("unmatch Date: want=%s got=%s", want, got)
		}
	})
}

func TestGet(t *testing.T) {
	ctx := context.Background()

	for _, c := range []struct {
		fundId string
		dummy  *string
	}{
		{"635333", nil},
	} {
		d, err := tokiomarineam.Get(ctx, c.fundId, c.dummy)
		if err != nil {
			t.Errorf("failed to get %+v: %s", c, err)
			continue
		}
		if dt := d.Date(); dt.IsZero() {
			t.Errorf("invalid date %+v", c)
		}
		if id := d.ID(); id != c.fundId {
			t.Errorf("unmatched ID for %+v: got=%s", c, id)
		}
		if p := d.Price(); p <= 0 {
			t.Errorf("invalid price %+v: %d", c, p)
		}
		if na := d.NetAssets(); na <= 0 {
			t.Errorf("invalid net assets %+v: %d", c, na)
		}
	}
}

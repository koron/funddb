package ammufg_test

import (
	"context"
	"testing"

	"github.com/koron/funddb/internal/adapter/ammufg"
)

func TestGet(t *testing.T) {
	ctx := context.Background()

	for _, c := range []struct {
		ctype  ammufg.CodeType
		code   string
		wantID string
	}{
		{ammufg.CodeTypeFund, "253425", "253425"},
		{ammufg.CodeTypeAssociationFund, "0331418A", "253425"},
		{ammufg.CodeTypeISIN, "JP90C000H1T1", "253425"},
	} {
		d, err := ammufg.Get(ctx, c.ctype, c.code)
		loc := string(c.ctype) + ":" + c.code
		if err != nil {
			t.Errorf("failed to get %q: %v", loc, err)
		}
		if dt := d.Date(); dt.IsZero() {
			t.Errorf("invalid date %q", loc)
		}
		if id := d.ID(); id != c.wantID {
			t.Errorf("unmatched ID for %q: want=%s got=%s", loc, c.wantID, id)
		}
		if p := d.Price(); p <= 0 {
			t.Errorf("invalid price %q: %d", loc, p)
		}
		if na := d.NetAssets(); na <= 0 {
			t.Errorf("invalid net assets %q: %d", loc, na)
		}
	}
}

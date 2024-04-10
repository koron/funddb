package fidelity_test

import (
	"context"
	"testing"

	"github.com/koron/funddb/internal/adapter/fidelity"
)

func TestGet(t *testing.T) {
	ctx := context.Background()

	for _, name := range []string{
		"267002/F",
		"217001/F",
		"267007/F",
	} {
		d, err := fidelity.Get(ctx, name)
		if err != nil {
			t.Errorf("failed to get %q: %v", name, err)
			continue
		}
		if dt := d.Date(); dt.IsZero() {
			t.Errorf("invalid date %q", name)
		}
		if id := d.ID(); id != name {
			t.Errorf("unmatched ID: want=%s got=%s", name, id)
		}
		if p := d.Price(); p <= 0 {
			t.Errorf("invalid price %q: %d", name, p)
		}
		if na := d.NetAssets(); na <= 0 {
			t.Errorf("invalid net assets %q: %d", name, na)
		}
	}
}

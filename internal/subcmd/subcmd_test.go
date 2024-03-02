package subcmd_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/koron/funddb/internal/subcmd"
)

func TestCommand(t *testing.T) {
	var called bool
	cmd := subcmd.DefineCommand("foo", t.Name(), func(context.Context, []string) error {
		called = true
		return nil
	})
	err := subcmd.Run(context.Background(), cmd)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("command func is not called")
	}
}

func TestCommandNil(t *testing.T) {
	cmd := subcmd.DefineCommand("foo", t.Name(), nil)
	err := subcmd.Run(context.Background(), cmd)
	if err == nil {
		t.Fatal("unexpected succeed")
	}
	if d := cmp.Diff("no function declared for command: foo", err.Error()); d != "" {
		t.Errorf("error unmatch: -want +got\n%s", d)
	}
}

func TestSet(t *testing.T) {
	var (
		gotNames []string
		gotArgs  []string
	)
	record := func(ctx context.Context, args []string) error {
		gotNames = subcmd.Names(ctx)
		gotArgs = args
		return nil
	}

	set := subcmd.DefineSet("set", "",
		subcmd.DefineSet("user", "",
			subcmd.DefineCommand("list", "", record),
			subcmd.DefineCommand("add", "", record),
			subcmd.DefineCommand("delete", "", record),
		),
		subcmd.DefineSet("post", "",
			subcmd.DefineCommand("list", "", record),
			subcmd.DefineCommand("add", "", record),
			subcmd.DefineCommand("delete", "", record),
		),
	)

	for i, c := range []struct {
		args      []string
		wantNames []string
		wantArgs  []string
	}{
		{
			[]string{"user", "list"},
			[]string{"set", "user", "list"},
			[]string{},
		},
		{
			[]string{"user", "add", "-email", "foobar@example.com"},
			[]string{"set", "user", "add"},
			[]string{"-email", "foobar@example.com"},
		},
		{
			[]string{"user", "delete", "-id", "123"},
			[]string{"set", "user", "delete"},
			[]string{"-id", "123"},
		},
		{
			[]string{"post", "list"},
			[]string{"set", "post", "list"},
			[]string{},
		},
		{
			[]string{"post", "add", "-title", "Brown fox..."},
			[]string{"set", "post", "add"},
			[]string{"-title", "Brown fox..."},
		},
		{
			[]string{"post", "delete", "-id", "ABC"},
			[]string{"set", "post", "delete"},
			[]string{"-id", "ABC"},
		},
	} {
		err := subcmd.Run(context.Background(), set, c.args...)
		if err != nil {
			t.Fatalf("failed for case#%d (%+v): %s", i, c, err)
			continue
		}
		if d := cmp.Diff(c.wantNames, gotNames); d != "" {
			t.Errorf("unexpected names on #%d: -want +got\n%s", i, d)
		}
		if d := cmp.Diff(c.wantArgs, gotArgs); d != "" {
			t.Errorf("unexpected args on #%d: -want +got\n%s", i, d)
		}
	}
}

func TestSetFails(t *testing.T) {
	set := subcmd.DefineSet("fail", "",
		subcmd.DefineCommand("list", "", nil),
		subcmd.DefineCommand("add", "", nil),
		subcmd.DefineCommand("delete", "", nil),
	)
	for i, c := range []struct {
		args []string
		want string
	}{
		{[]string{}, "required one of name from: [list add delete]"},
		{[]string{"foo"}, "given \"foo\" is not one of name in: [list add delete]"},
	} {
		err := subcmd.Run(context.Background(), set, c.args...)
		if err == nil {
			t.Fatalf("unexpected succeed at #%d %+v", i, c)
		}
		got := err.Error()
		if d := cmp.Diff(c.want, got); d != "" {
			t.Errorf("unexpected error at #%d: -want +got\n%s", i, d)
		}
	}
}

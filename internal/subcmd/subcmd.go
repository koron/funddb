package subcmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Runner defines a base interface for Command and Set.
// Runner interface is defined for use only with DefineSet function.
type Runner interface {
	name() string
	desc() string
	run(ctx context.Context, args []string) error
}

type CommandFunc func(ctx context.Context, args []string) error

type Command struct {
	Name string
	Desc string
	Func CommandFunc
}

var _ Runner = Command{}

// DefineCommand defines a Command with name, desc, and function.
func DefineCommand(name, desc string, fn CommandFunc) Command {
	return Command{
		Name: name,
		Desc: desc,
		Func: fn,
	}
}

func (c Command) name() string {
	return c.Name
}

func (c Command) desc() string {
	return c.Desc
}

func (c Command) run(ctx context.Context, args []string) error {
	ctx = withName(ctx, c)
	if c.Func == nil {
		names := strings.Join(Names(ctx), " ")
		return fmt.Errorf("no function declared for command: %s", names)
	}
	return c.Func(ctx, args)
}

type Set struct {
	Name    string
	Desc    string
	Runners []Runner
}

var _ Runner = Set{}

// DefineSet defines a set of Runners with name, and desc.
func DefineSet(name, desc string, runners ...Runner) Set {
	return Set{
		Name:    name,
		Desc:    desc,
		Runners: runners,
	}
}

// DefineRootSet defines a set of Runners which used as root of Set (maybe
// passed to Run).
func DefineRootSet(runners ...Runner) Set {
	return Set{Name: rootName(), Runners: runners}
}

func (s Set) name() string {
	return s.Name
}

func (s Set) desc() string {
	return s.Desc
}

func (s Set) runnerNames() []string {
	a := make([]string, 0, len(s.Runners))
	for _, r := range s.Runners {
		if n := r.name(); n != "" {
			a = append(a, n)
		}
	}
	return a
}

// childRunner retrieves a child Runner with name
func (s Set) childRunner(name string) Runner {
	for _, r := range s.Runners {
		if r.name() == name {
			return r
		}
	}
	return nil
}

type errorSetRun struct {
	src Set
	msg string
}

func (err *errorSetRun) Error() string {
	// align width of name columns
	var w int = 12
	for _, r := range err.src.Runners {
		if n := len(r.name()) + 1; n > w {
			w = (n + 3) / 4 * 4
		}
	}
	// format error message
	bb := &bytes.Buffer{}
	fmt.Fprintf(bb, "%s.\n\nAvailable sub-commands are:\n", err.msg)
	for _, r := range err.src.Runners {
		fmt.Fprintf(bb, "\n\t%-*s%s", w, r.name(), r.desc())
	}
	return bb.String()
}

func (s Set) run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return &errorSetRun{src: s, msg: "no commands selected"}
	}
	name := args[0]
	child := s.childRunner(name)
	if child == nil {
		return &errorSetRun{src: s, msg: "command not found"}
	}
	return child.run(withName(ctx, s), args[1:])
}

// Run runs a Runner with ctx and args.
func Run(ctx context.Context, r Runner, args ...string) error {
	return r.run(ctx, args)
}

var keyNames = struct{}{}

// Names retrives names layer of current sub command.
func Names(ctx context.Context) []string {
	if names, ok := ctx.Value(keyNames).([]string); ok {
		return names
	}
	return nil
}

func withName(ctx context.Context, r Runner) context.Context {
	return context.WithValue(ctx, keyNames, append(Names(ctx), r.name()))
}

func rootName() string {
	exe, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("failed to obtain executable name: %s", err))
	}
	_, name := filepath.Split(exe)
	ext := filepath.Ext(name)
	if ext == ".exe" {
		return name[:len(name)-len(ext)]
	}
	return name
}

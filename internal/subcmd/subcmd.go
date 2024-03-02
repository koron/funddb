package subcmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var keyNames = struct{}{}

// Runner defines a base interface for Command and Set.
// Runner interface is defined for use only with DefineSet function.
type Runner interface {
	name() string
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

func DefineRootSet(runners ...Runner) Set {
	return Set{Name: rootName(), Runners: runners}
}

func (s Set) name() string {
	return s.Name
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

func (s Set) childRunner(name string) Runner {
	for _, r := range s.Runners {
		if r.name() == name {
			return r
		}
	}
	return nil
}

func (s Set) run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		names := s.runnerNames()
		return fmt.Errorf("required one of name from: %s", names)
	}
	name := args[0]
	child := s.childRunner(name)
	if child == nil {
		names := s.runnerNames()
		return fmt.Errorf("given %q is not one of name in: %s", name, names)
	}
	return child.run(withName(ctx, s), args[1:])
}

func withName(ctx context.Context, r Runner) context.Context {
	return context.WithValue(ctx, keyNames, append(Names(ctx), r.name()))
}

// Run runs a Runner with ctx and args.
func Run(ctx context.Context, r Runner, args ...string) error {
	return r.run(ctx, args)
}

func rootName() string {
	exe, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("failed to obtain executable name: %s", err))
	}
	_, name := filepath.Split(exe)
	ext := filepath.Ext(name)
	return name[:len(name)-len(ext)]
}

// Names retrives names layer of current sub command.
func Names(ctx context.Context) []string {
	if names, ok := ctx.Value(keyNames).([]string); ok {
		return names
	}
	return nil
}

func NewFlagSet(ctx context.Context) *flag.FlagSet {
	name := strings.Join(Names(ctx), " ")
	return flag.NewFlagSet(name, flag.ExitOnError)
}

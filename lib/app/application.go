package app

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sort"
	"time"

	"github.com/betam/glb/lib/di"
	"github.com/betam/glb/lib/try"
)

type commandTag struct{}

var (
	CommandTag = commandTag{}
)

type Configure interface{}

type App interface {
	Run()
	Stop()
}

func New() App {
	di.Wire[Configure](func() Configure { return struct{ Configure }{} }, di.Fallback())
	di.Wire[App](func(commands []Command, cfg Configure) App { return Init(commands) }, di.Defaults(map[int]any{
		0: di.Tags[Command](CommandTag),
	}))
	application, closer := di.NewWithCloser[App]()
	if appInstance, ok := application.(*app); ok {
		appInstance.closer = closer
	}
	return application
}

func Init(commands []Command) *app {
	var cmd CommandList
	for _, command := range commands {
		cmd.Add(command)
	}
	return &app{commands: cmd}
}

type app struct {
	commands  CommandList
	config    Configure
	closer    func()
	ctxCancel context.CancelFunc
}

func (a *app) Run() {
	opts := getopt.New()
	var help bool
	opts.FlagLong(&help, "help", 'h', "Display this help message")
	opts.SetParameters("command")

	opts.Parse(os.Args)
	command := opts.Arg(0)
	if command == "" {
		help = true
	}
	if help {
		a.usage(opts.PrintUsage)
	}

	if cmd := a.commands.Get(command); cmd != nil {
		ctx, cancel := NewContextWithCancelWithCommand(cmd.Name(), opts.Args())
		a.ctxCancel = cancel

		os.Args = opts.Args()
		try.Catch(
			func() {
				cmd.Run(ctx)
			},
			func(throwable error) {
				logrus.WithContext(ctx).Error(throwable)
				panic(throwable)
			},
		)
	} else {
		color.Red("Unknown command '%s'.\n", command)
		a.usage(opts.PrintUsage)
		os.Exit(2)
	}
}

func (a *app) Stop() {
	if a.ctxCancel != nil {
		a.ctxCancel()
		time.Sleep(1 * time.Second)
	}
	if a.closer != nil {
		a.closer()
	}
}

func (a *app) usage(printer func(w io.Writer)) {
	printer(os.Stderr)
	_, _ = fmt.Fprintf(os.Stderr, "Available commands:\n")
	var names []string
	for command := range a.commands {
		names = append(names, command)
	}
	sort.Strings(names)
	for cmd := range a.commands.List() {
		fmt.Printf("\t%s\t%s\n", color.New(color.FgGreen).Sprint(cmd.Name()), cmd.Description())
	}
	os.Exit(2)
}

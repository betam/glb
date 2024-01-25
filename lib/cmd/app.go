package cmd

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"

	"github.com/betam/glb/lib/try"
)

func Run(commandList CommandList) {
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
		usage(commandList, opts.PrintUsage)
	}

	if cmd, ok := commandList.Get(command); ok {
		ctx := NewContextWithCommand(cmd.Name, opts.Args())
		os.Args = opts.Args()
		try.Catch(
			func() {
				cmd.Handler(ctx)
			},
			func(throwable error) {
				logrus.WithContext(ctx).Error(throwable)
				panic(throwable)
			},
		)
	} else {
		color.Red("Unknown command '%s'.\n", command)
		usage(commandList, opts.PrintUsage)
		os.Exit(2)
	}
}

func usage(commandList CommandList, printer func(w io.Writer)) {
	printer(os.Stderr)
	_, _ = fmt.Fprintf(os.Stderr, "Available commands:\n")
	var names []string
	for command := range commandList {
		names = append(names, command)
	}
	sort.Strings(names)
	for cmd := range commandList.List() {
		fmt.Printf("\t%s\t%s\n", color.New(color.FgGreen).Sprint(cmd.Name), cmd.Description)
	}
	os.Exit(2)
}

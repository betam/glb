package cmd

import (
	"context"
	"fmt"
	"sort"

	"github.com/betam/glb/lib/list"
)

type Command struct {
	Name        string
	Description string
	Handler     func(ctx context.Context)
}

type CommandList map[string]*Command

func (l *CommandList) Add(name, description string, handler func(ctx context.Context)) {
	if name == "" {
		panic(fmt.Errorf("empty command name"))
	}
	if handler == nil {
		panic(fmt.Errorf("empty command handler"))
	}
	if *l == nil {
		*l = make(CommandList)
	}
	if _, ok := (*l)[name]; ok {
		panic(fmt.Errorf("command '%s' has been already registered", name))
	}
	(*l)[name] = &Command{Name: name, Handler: handler, Description: description}
}

func (l *CommandList) Command(command string) (func(ctx context.Context), bool) {
	if cmd, ok := l.Get(command); ok {
		return cmd.Handler, ok
	}
	return nil, false
}

func (l *CommandList) Get(command string) (*Command, bool) {
	details, ok := (*l)[command]
	return details, ok
}

func (l *CommandList) List() <-chan *Command {
	ch := make(chan *Command, len(*l))
	commands := list.Values(*l)
	sort.SliceStable(
		commands,
		func(i, j int) bool {
			return commands[i].Name < commands[j].Name
		},
	)
	for _, value := range commands {
		ch <- value
	}
	close(ch)
	return ch
}

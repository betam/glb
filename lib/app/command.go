package app

import (
	"context"
	"fmt"
	"github.com/betam/glb/lib/list"
	"sort"
)

type Command interface {
	Name() string
	Description() string
	Run(ctx context.Context)
}

type CommandList map[string]Command

func (l *CommandList) Add(command Command) {
	if command.Name() == "" {
		panic(fmt.Errorf("empty command name"))
	}
	if *l == nil {
		*l = make(CommandList)
	}
	if _, ok := (*l)[command.Name()]; ok {
		panic(fmt.Errorf("command '%s' has been already registered", command.Name()))
	}
	(*l)[command.Name()] = command
}

func (l *CommandList) Get(command string) Command {
	details, _ := (*l)[command]
	return details
}

func (l *CommandList) List() <-chan Command {
	ch := make(chan Command, len(*l))
	commands := list.Values(*l)
	sort.SliceStable(
		commands,
		func(i, j int) bool {
			return commands[i].Name() < commands[j].Name()
		},
	)
	for _, value := range commands {
		ch <- value
	}
	close(ch)
	return ch
}

func NewCommand(name, description string) *BaseCommand {
	return &BaseCommand{
		name:        name,
		description: description,
	}
}

type BaseCommand struct {
	name, description string
}

func (b *BaseCommand) Name() string {
	return b.name
}

func (b *BaseCommand) Description() string {
	return b.description
}

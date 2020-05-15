package command

import (
	"fmt"

	"github.com/urfave/cli"
)

const (
	defaultHelloString = "SomeDefaultHelloString"
)

var (
	helloString string
)

type Builder struct {
	commands   []cli.Command
}

// NewBuilder creates and returns a new command builder.
func NewBuilder() *Builder {
	return &Builder{}
}

// GetCommands returns the list of allowed commands.
func (b *Builder) GetCommands() []cli.Command {
	return b.commands
}

// Hello simply prints hello and a string.
func (b *Builder) Hello() *Builder {
	detailsCmd := cli.Command{
		Name:    "hello",
		Aliases: []string{"s"},
		Usage:   "Prints hello.",
		Flags:   []cli.Flag{
			cli.StringFlag{
				Name:        "hello, s",
				Usage:       "Simply print hello.",
				Value:       defaultHelloString,
				Destination: &helloString,
				Required:    false,
			},
		},
		Action: func(c *cli.Context) {
			fmt.Printf("Hello %v!!\n", helloString)
		},
	}

	b.commands = append(b.commands, detailsCmd)

	return b
}
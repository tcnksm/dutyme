package command

import (
	"strings"
)

type StartCommand struct {
	Meta
}

func (c *StartCommand) Run(args []string) int {
	// Write your code here

	return 0
}

func (c *StartCommand) Synopsis() string {
	return ""
}

func (c *StartCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}

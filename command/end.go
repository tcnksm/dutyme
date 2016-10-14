package command

import (
	"strings"
)

type EndCommand struct {
	Meta
}

func (c *EndCommand) Run(args []string) int {
	// Write your code here

	return 0
}

func (c *EndCommand) Synopsis() string {
	return ""
}

func (c *EndCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}

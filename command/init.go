package command

import (
	"strings"
)

type InitCommand struct {
	Meta
}

func (c *InitCommand) Run(args []string) int {
	// Write your code here

	return 0
}

func (c *InitCommand) Synopsis() string {
	return ""
}

func (c *InitCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}

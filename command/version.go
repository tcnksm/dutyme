package command

import (
	"bytes"
	"fmt"
)

type VersionCommand struct {
	Meta

	Name     string
	Version  string
	Revision string
}

func (c *VersionCommand) Run(args []string) int {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "%s version %s", c.Name, c.Version)
	if c.Revision != "" {
		fmt.Fprintf(&buf, " (%s)", c.Revision)
	}

	fmt.Fprintln(c.OutStream, buf.String())
	return 0
}

func (c *VersionCommand) Synopsis() string {
	return fmt.Sprintf("Print %s version and quit", c.Name)
}

func (c *VersionCommand) Help() string {
	return ""
}

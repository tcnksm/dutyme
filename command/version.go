package command

import (
	"bytes"
	"fmt"

	latest "github.com/tcnksm/go-latest"
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

	githubTag := &latest.GithubTag{
		Owner:             "tcnksm",
		Repository:        "dutyme",
		FixVersionStrFunc: latest.DeleteFrontV(),
	}

	res, err := latest.Check(githubTag, c.Version)
	if err != nil {
		Debugf("Failed to check latest version: %s", err)
		return 0
	}

	if res.Outdated {
		fmt.Fprintln(c.OutStream,
			"\n Your version of `dutyme` is out of date! The latest version is %s.\n"+
				"You can donwloand it from github.com/tcnksm/dutyme", res.Latest)
	}
	return 0
}

func (c *VersionCommand) Synopsis() string {
	return fmt.Sprintf("Print %s version and quit", c.Name)
}

func (c *VersionCommand) Help() string {
	return ""
}

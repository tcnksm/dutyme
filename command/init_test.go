package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestInitCommand_implement(t *testing.T) {
	var _ cli.Command = &InitCommand{}
}

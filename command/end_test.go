package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestEndCommand_implement(t *testing.T) {
	var _ cli.Command = &EndCommand{}
}

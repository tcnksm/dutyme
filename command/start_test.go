package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestStartCommand_implement(t *testing.T) {
	var _ cli.Command = &StartCommand{}
}

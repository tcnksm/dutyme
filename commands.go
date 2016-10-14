package main

import (
	"github.com/mitchellh/cli"
	"github.com/tcnksm/dutyme/command"
)

func Commands(meta *command.Meta) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"init": func() (cli.Command, error) {
			return &command.InitCommand{
				Meta: *meta,
			}, nil
		},
		"start": func() (cli.Command, error) {
			return &command.StartCommand{
				Meta: *meta,
			}, nil
		},
		"end": func() (cli.Command, error) {
			return &command.EndCommand{
				Meta: *meta,
			}, nil
		},

		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Meta:     *meta,
				Version:  Version,
				Revision: GitCommit,
				Name:     Name,
			}, nil
		},
	}
}

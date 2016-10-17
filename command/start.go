package command

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tcnksm/dutyme/config"
	"github.com/tcnksm/dutyme/dutyme"
	"github.com/tcnksm/go-input"
)

type StartCommand struct {
	Meta
}

const (
	EnvToken = "PD_SERVICE_KEY"
)

const (
	ExitCodeOK = iota
	ExitCodeError
)

func (c *StartCommand) Run(args []string) int {
	flags := c.Meta.NewFlagSet("start", c.Help())
	if err := flags.Parse(args); err != nil {
		return ExitCodeError
	}

	cfgPath, err := c.Meta.ConfigPath()
	if err != nil {
		log.Fatal(err)
		return ExitCodeError
	}

	cfg := &config.Config{}
	if _, err := os.Stat(cfgPath); err == nil {
		var err error
		cfg, err = config.ParseFile(cfgPath)
		if err != nil {
			log.Fatal(err)
			return ExitCodeError
		}
	}

	if v := os.Getenv(EnvToken); len(v) != 0 {
		cfg.Token = v
	}

	if len(cfg.Token) == 0 {
		token, err := c.Meta.AskToken("")
		if err != nil {
			log.Fatal(err)
			return ExitCodeError
		}
		cfg.Token = token
	}

	pd, err := dutyme.NewPDClient(cfg.Token)
	if err != nil {
		log.Fatal(err)
		return ExitCodeError
	}

	dutyme := &dutyme.Dutyme{
		UI: c.Meta.UI,
		PD: pd,
	}

	if cfg.User == nil {
		user, err := dutyme.GetUser("")
		if err != nil {
			log.Fatal(err)
			return ExitCodeError
		}

		cfg.User = user
	}

	if cfg.ScheduleID == "" {
		scheduleName, scheduleID, err := dutyme.GetSchedule("")
		if err != nil {
			log.Fatal(err)
			return ExitCodeError
		}

		cfg.ScheduleName = scheduleName
		cfg.ScheduleID = scheduleID
	}

	// TODO(tcnsm): Enable to set working time via flag
	DefaultWorkingTime := 1 * time.Hour
	start := time.Now()
	end := start.Add(DefaultWorkingTime)

	force := false
	log.Printf("Override user %q to schedule %q from %s to %s",
		cfg.User.Email, cfg.ScheduleName, start, end)
	override, err := dutyme.Override(cfg.ScheduleID, cfg.User, start, end, force)
	if err != nil {
		log.Fatal(err)
		return ExitCodeError
	}
	log.Println("Successfuly overrided:", override.ID)

	query := "Want to save login info? (you can skip input from next time) [Y/n]"
	ans, err := c.Meta.UI.Ask(query, &input.Options{
		Default:     "Y",
		Loop:        true,
		HideOrder:   true,
		HideDefault: true,
		ValidateFunc: func(s string) error {
			if s != "Y" && s != "y" && s != "N" && s != "n" {
				return fmt.Errorf("input must be Y or n")
			}
			return nil
		},
	})

	if ans == "N" || ans == "n" {
		return ExitCodeOK
	}

	indent := true
	if err := cfg.WriteFile(cfgPath, indent); err != nil {
		log.Fatal(err)
		return ExitCodeError
	}

	return ExitCodeOK
}

func (c *StartCommand) Synopsis() string {
	return ""
}

func (c *StartCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}

package command

import (
	"fmt"
	"os"
	"time"

	"github.com/tcnksm/dutyme/config"
	"github.com/tcnksm/dutyme/dutyme"
	"github.com/tcnksm/go-input"
)

const (
	// DefaultWorkingTime is default duration of override.
	DefaultWorkingTime = 1 * time.Hour

	// TimeFmt is used for displaying time
	TimeFmt = "2006-01-02 15:04:05-0700"
)

type StartCommand struct {
	Meta
}

func (c *StartCommand) Synopsis() string {
	return "Override schedule and assign on-call to you"
}

func (c *StartCommand) Help() string {
	helpText := fmt.Sprintf(`Usage: dutyme start [options...]

start overrides the schedule and assign on-call to you. By default,
it creates 1 hour override (you can change this via -working option).

To use dutyme command, you need a PagerDuty API v2 token.
The token must have full access to read, write, update, and delete.

Only account administrators have the ability to generate token.
(See more about token on official doc https://goo.gl/VPvlwB)

You can set the token via %q env var (then you can
skip the following input). Or after running first overriding, you can
save login info on config file.

Options:

  -working TIME  Working time (overriding time). By default, it's 1 hour.
                 TIME can be specified by decimal numbers with a unit suffix,
                 such "1.5h" or "2h45m". It must be positive value.

  -update        Update existing configuration file. It asks email and
                 schedule name again.

  -force         Force overriding without confirmation.

`, EnvToken)
	return helpText
}

func (c *StartCommand) Run(args []string) int {

	var (
		force       bool
		update      bool
		workingTime time.Duration

		useExisting bool
	)

	flags := c.Meta.NewFlagSet("start", c.Help())

	flags.BoolVar(&force, "force", false, "")
	flags.BoolVar(&force, "f", false, "")

	flags.BoolVar(&update, "update", false, "")
	flags.DurationVar(&workingTime, "working", DefaultWorkingTime, "")

	if err := flags.Parse(args); err != nil {
		return ExitCodeError
	}

	// Find configuration file for dutyme.
	cfgPath, err := c.Meta.ConfigPath()
	if err != nil {
		fmt.Fprintf(c.ErrStream, "Failed to read config path: %s\n", err)
		TracePrint(c.ErrStream, err)
		return ExitCodeError
	}

	cfg := &config.Config{}
	if _, err := os.Stat(cfgPath); err == nil {
		Debugf("Use existing configuration file: %s", cfgPath)

		var err error
		cfg, err = config.ParseFile(cfgPath)
		if err != nil {
			fmt.Fprintf(c.ErrStream, "Failed to parse configuration file: %s\n", err)
			TracePrint(c.ErrStream, err)
			return ExitCodeError
		}
		useExisting = true
	}

	// Override API token via env var if exsits.
	if v := os.Getenv(EnvToken); len(v) != 0 {
		Debugf("Read PD API token from env var: %s", v)
		cfg.Token = v
	}

	if len(cfg.Token) == 0 {
		token, err := c.Meta.AskToken()
		if err != nil {
			fmt.Fprintf(c.ErrStream, "Failed to ask API token: %s\n", err)
			TracePrint(c.ErrStream, err)
			return ExitCodeError
		}
		cfg.Token = token
	}

	// Construct Dutyme client with PD HTTP client
	pd, err := dutyme.NewPDClient(cfg.Token)
	if err != nil {
		fmt.Fprintf(c.ErrStream, "Failed to create PD HTTP client: %s\n", err)
		TracePrint(c.ErrStream, err)
		return ExitCodeError
	}

	dutyme := &dutyme.Dutyme{
		UI: c.Meta.UI,
		PD: pd,
	}

	// When configuration file is not exist (fisrt time to execute or not saved before).
	// or when -update flag is provided, ask/get user information.
	if cfg.IsEmpty() || update {
		user, err := dutyme.GetUser("")
		if err != nil {
			fmt.Fprintf(c.ErrStream, "Failed to get PagerDuty user: %s\n", err)
			TracePrint(c.ErrStream, err)
			return ExitCodeError
		}
		cfg.User = user

		scheduleName, scheduleID, err := dutyme.GetSchedule("")
		if err != nil {
			fmt.Fprintf(c.ErrStream, "Failed to get PagerDuty schedule: %s\n", err)
			TracePrint(c.ErrStream, err)
			return ExitCodeError
		}
		cfg.ScheduleName = scheduleName
		cfg.ScheduleID = scheduleID
	}
	Debugf("User: %s", cfg.User.Email)
	Debugf("Schedule: %s", cfg.ScheduleName)

	// Override time: from now to now + working time
	start := time.Now()
	end := start.Add(workingTime)

	fmt.Fprintf(c.OutStream, "Override schedule %q (%s) by user %q\n",
		cfg.ScheduleName, cfg.ScheduleID, cfg.User.Email)
	fmt.Fprintf(c.OutStream, "from %s to %s\n",
		start.Format(TimeFmt), end.Format(TimeFmt))

	override, err := dutyme.Override(cfg.ScheduleID, cfg.User, start, end, force)
	if err != nil {
		if IsCancel(err) {
			fmt.Fprintln(c.OutStream, "Override canceled")
			return ExitCodeError
		}

		fmt.Fprintf(c.ErrStream, "Failed to override: %s\n", err)
		TracePrint(c.ErrStream, err)
		return ExitCodeError
	}
	fmt.Fprintf(c.OutStream, "Successfuly overrided schedule (%s)\n", override.ID)

	// If it's used exsiting configuration file,
	// and -update flag is not provided, then skip the following section.
	if useExisting && !update {
		return ExitCodeOK
	}

	// Save override info on file
	query := "Want to save override info? (you can skip input from next time) [Y/n]"
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

	if err := cfg.WriteFile(cfgPath, true); err != nil {
		fmt.Fprintf(c.ErrStream, "Failed to save file: %s\n", err)
		TracePrint(c.ErrStream, err)
		return ExitCodeError
	}
	fmt.Fprintf(c.OutStream, "Successfuly saved file (%s)\n", cfgPath)

	return ExitCodeOK
}

type isCancel interface {
	IsCancel() bool
}

func IsCancel(err error) bool {
	e, ok := err.(isCancel)
	return ok && e.IsCancel()
}

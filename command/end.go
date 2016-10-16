package command

import (
	"log"
	"os"
	"strings"

	"github.com/tcnksm/dutyme/config"
	"github.com/tcnksm/dutyme/dutyme"
	"github.com/tcnksm/go-input"
)

type EndCommand struct {
	Meta
}

func (c *EndCommand) Run(args []string) int {
	cfgPath, err := c.Meta.ConfigPath()
	if err != nil {
		log.Fatal(err)
		return ExitCodeError
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatal(err)
		return ExitCodeError
	}

	cfg, err := config.ParseFile(cfgPath)
	if err != nil {
		log.Fatal(err)
		return ExitCodeError
	}

	if !cfg.IsOverrideExist() {
		log.Fatal("Currently no override exists")
		return ExitCodeError
	}

	if v := os.Getenv(EnvToken); len(v) != 0 {
		cfg.Token = v
	}

	if len(cfg.Token) == 0 {
		query := "Input PagerDuty API token"
		token, err := c.Meta.UI.Ask(query, &input.Options{
			Required:  true,
			Loop:      true,
			HideOrder: true,
			Mask:      true,
		})

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

	if err := pd.DeleteOverride(cfg.OverrideScheduleID, cfg.OverrideID); err != nil {
		log.Fatal(err)
		return ExitCodeError
	}
	log.Printf("Successfully delete override on %s", cfg.OverrideScheduleID)

	// Remove existing override from config file.
	cfg.OverrideID = ""
	cfg.OverrideScheduleID = ""

	indent := true
	if err := cfg.WriteFile(cfgPath, indent); err != nil {
		log.Fatal(err)
		return ExitCodeError
	}

	return ExitCodeError
}

func (c *EndCommand) Synopsis() string {
	return ""
}

func (c *EndCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}

package dutyme

import (
	"fmt"
	"strings"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/pkg/errors"
	"github.com/tcnksm/go-gitconfig"
	"github.com/tcnksm/go-input"
)

type Dutyme struct {
	PD PagerDuty
	UI *input.UI
}

func (d *Dutyme) GetUser(defaultEmail string) (*User, error) {
	if len(defaultEmail) == 0 {
		// PD email address may be same as git email address
		defaultEmail, _ = gitconfig.Email()
	}

	query := "Input PagerDuty account email address"
	email, err := d.UI.Ask(query, &input.Options{
		Default:   defaultEmail,
		Required:  true,
		Loop:      true,
		HideOrder: true,
	})

	if err != nil {
		return nil, errors.Wrap(err, "faield to ask PD email address")
	}

	return d.PD.GetUser(email)
}

func (d *Dutyme) GetSchedule(defaultQuery string) (string, string, error) {
	query := "Input PagerDuty schedule name which you want to override"
	scheduleQuery, err := d.UI.Ask(query, &input.Options{
		Default:   defaultQuery,
		Required:  true,
		Loop:      true,
		HideOrder: true,
	})

	schedules, err := d.PD.GetSchedules(scheduleQuery)
	if err != nil {
		return "", "", err
	}

	var (
		name string
		ID   string
	)

	// API may return multiple schedules. Then ask to user.
	if len(schedules) > 1 {
		targets := make([]string, 0, len(schedules))
		for _, schedule := range schedules {
			targets = append(targets, schedule.Name)
		}

		query := "Found multiple schedules. Select one."
		target, err := d.UI.Select(query, targets, &input.Options{
			Default: targets[0],
			Loop:    true,
		})

		if err != nil {
			return "", "", errors.Wrap(err, "failed to ask schedule from the given list")
		}

		for _, schedule := range schedules {
			if target == schedule.Name {
				name = schedule.Name
				ID = schedule.ID
			}
		}
	} else {
		name = schedules[0].Name
		ID = schedules[0].ID
	}

	return name, ID, nil
}

func (d *Dutyme) Override(scheduleID string, user *User, start, end time.Time, force bool) (*pagerduty.Override, error) {

	if !force {
		query := fmt.Sprintf("OK to override? [Y/n]")
		ans, err := d.UI.Ask(query, &input.Options{
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

		if err != nil {
			return nil, errors.Wrap(err, "failed to ask")
		}

		if ans == "N" || ans == "n" {
			return nil, &errCancel{}
		}
	}

	return d.PD.Override(scheduleID, user, start, end)
}

func (d *Dutyme) GetOverride(scheduleID string, user *User, since, until time.Time) (string, error) {
	overrides, err := d.PD.GetOverrides(scheduleID, since, until)
	if err != nil {
		return "", err
	}

	targets := make([]string, 0, len(overrides))
	for _, override := range overrides {
		// Filter override by user
		if override.User.ID == user.Obj.ID {
			s := fmt.Sprintf("%s: %s - %s", override.ID, override.Start, override.End)
			targets = append(targets, s)
		}
	}

	var target string
	if len(targets) > 1 {
		var err error
		query := "Found multiple overrides. Select one."
		target, err = d.UI.Select(query, targets, &input.Options{
			Default: targets[0],
			Loop:    true,
		})
		if err != nil {
			return "", errors.Wrap(err, "failed to select override from the given list")
		}

	} else {
		target = targets[0]
	}

	return target[:strings.Index(target, ":")], nil
}

// errCancel implements isCancel interface.
type errCancel struct{}

func (_ *errCancel) Error() string {
	return "canceled"
}

func (_ *errCancel) IsCancel() bool {
	return true
}

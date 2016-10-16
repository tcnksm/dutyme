package dutyme

import (
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

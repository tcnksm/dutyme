package dutyme

import (
	"strings"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/pkg/errors"
)

type PagerDuty interface {
	GetUser(email string) (*pagerduty.APIObject, error)
	GetSchedules(name string) ([]pagerduty.Schedule, error)
	Override(scheduleID string, user pagerduty.APIObject, start, end time.Time) (*pagerduty.Override, error)
	DeleteOverride(scheduleID, overrideID string) error
}

type PDClient struct {
	*pagerduty.Client
}

// NewPDClient creates new PagerDuty client
func NewPDClient(token string) (PagerDuty, error) {
	if len(token) == 0 {
		return nil, errors.New("missing Pagerduty API token")
	}
	pg := pagerduty.NewClient(token)
	return &PDClient{
		Client: pg,
	}, nil
}

// GetUser gets a user by the given email address and returns
// its APIObject (it represents user).
//
// FindUser assumes one email belongs to one user.
// If it finds more than two, it fails.
func (c *PDClient) GetUser(email string) (*pagerduty.APIObject, error) {
	if len(email) == 0 {
		return nil, errors.New("missing pagerduty account email")
	}

	// TODO(tcnksm): More strict search?
	res, err := c.ListUsers(pagerduty.ListUsersOptions{
		Query: email,
	})

	if err != nil {
		return nil, errors.Wrap(err, "PagerDuty API request failed: ListUsers")
	}

	users := res.Users
	if len(users) == 0 {
		return nil, errors.Errorf("no such user: %s (correct email?)", email)
	}

	// Assumption: One email belongs to only one user.
	// Please file a ticket if you find exception.
	if len(users) != 1 {
		names := make([]string, 0, len(users))
		for _, u := range users {
			names = append(names, u.Name)
		}
		return nil, errors.Errorf(
			"more than 2 users are found: %s", strings.Join(names, ","))
	}

	return &users[0].APIObject, nil
}

// GetSchecules finds Pagerduty schedules by querying the given name.
// If any or found nothing, returns error.
func (c *PDClient) GetSchedules(name string) ([]pagerduty.Schedule, error) {
	if len(name) == 0 {
		return nil, errors.New("missing schedule name")
	}

	// TODO(tcnksm): More strict search?
	res, err := c.ListSchedules(pagerduty.ListSchedulesOptions{
		Query: name,
	})
	if err != nil {
		return nil, errors.Wrap(err, "PagerDuty API request failed: ListSchedules")
	}

	schedules := res.Schedules
	if len(schedules) == 0 {
		return nil, errors.Errorf("no such schedule: %s ", name)
	}

	return schedules, nil
}

func (c *PDClient) Override(scheduleID string, user pagerduty.APIObject, start, end time.Time) (*pagerduty.Override, error) {
	if len(scheduleID) == 0 {
		return nil, errors.New("misssing scheduleID")
	}

	if user.ID == "" {
		return nil, errors.New("missing user.ID")
	}

	if start.IsZero() || end.IsZero() {
		return nil, errors.New("start and end time should be non-zero value")
	}

	// TODO(tcnksm): Handle when user is already persion in charge.
	override, err := c.CreateOverride(scheduleID, pagerduty.Override{
		Start: start.String(),
		End:   end.String(),
		User:  user,
	})

	if err != nil {
		return nil, errors.Wrap(err, "PagerDuty API request failed: CreateOverride")
	}

	return override, nil
}

func (c *PDClient) DeleteOverride(scheduleID, overrideID string) error {
	if len(scheduleID) == 0 {
		return errors.New("missing schedule ID")
	}

	if len(overrideID) == 0 {
		return errors.New("missing override ID")
	}

	if err := c.DeleteOverride(scheduleID, overrideID); err != nil {
		return errors.Wrap(err, "PagerDuty API request failed: DeleteOverride")
	}

	return nil
}

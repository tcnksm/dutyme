package dutyme

import (
	"strings"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/pkg/errors"
)

type PagerDuty interface {
	GetUser(email string) (*User, error)
	GetSchedules(name string) ([]pagerduty.Schedule, error)
	GetOverrides(scheduleID string, since, until time.Time) ([]pagerduty.Override, error)
	Override(scheduleID string, user *User, start, end time.Time) (*pagerduty.Override, error)
	DeleteOverride(scheduleID, overrideID string) error
}

// User represents pagerduty user
type User struct {
	Email string
	Obj   *pagerduty.APIObject
}

// PDClient is actual pagerduty client which implements PagerDuty interface.
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
func (c *PDClient) GetUser(email string) (*User, error) {
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

	return &User{
		Email: email,
		Obj:   &users[0].APIObject,
	}, nil
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

func (c *PDClient) GetOverrides(scheduleID string, since, until time.Time) ([]pagerduty.Override, error) {
	if len(scheduleID) == 0 {
		return nil, errors.New("misssing scheduleID")
	}

	overrides, err := c.ListOverrides(scheduleID, pagerduty.ListOverridesOptions{
		Since:    since.String(),
		Until:    until.String(),
		Editable: true,
		Overflow: true,
	})

	if err != nil {
		return nil, errors.Wrap(err, "PagerDuty API request failed: ListOverrides")
	}

	if len(overrides) == 0 {
		return nil, &errNotFound{"no overrides are found"}
	}

	return overrides, nil
}

func (c *PDClient) Override(scheduleID string, user *User, start, end time.Time) (*pagerduty.Override, error) {
	if len(scheduleID) == 0 {
		return nil, errors.New("misssing scheduleID")
	}

	if user == nil {
		return nil, errors.New("missing user")
	}

	if start.IsZero() || end.IsZero() {
		return nil, errors.New("start and end time should be non-zero value")
	}

	// TODO(tcnksm): Handle when user is already persion in charge.
	override, err := c.CreateOverride(scheduleID, pagerduty.Override{
		Start: start.String(),
		End:   end.String(),
		User:  *user.Obj,
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

	if err := c.Client.DeleteOverride(scheduleID, overrideID); err != nil {
		return errors.Wrap(err, "PagerDuty API request failed: DeleteOverride")
	}

	return nil
}

// IsNotFound returns true if err inplements notfound interface
// and NotFound returns true.
func isNotFound(err error) bool {
	i, ok := err.(notfound)
	return ok && i.NotFound()
}

type notfound interface {
	NotFound() bool
}

// errNotFound is special error when API requests found no resources.
type errNotFound struct {
	msg string
}

func (e *errNotFound) Error() string {
	return e.msg
}

func (e *errNotFound) NotFound() bool {
	return true
}

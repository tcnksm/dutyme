package dutyme

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/pkg/errors"
)

const (
	// TestToken is read-only test token provided by PagerDuty team
	// https://v2.developer.pagerduty.com/v2/page/api-reference
	TestToken = "w_8PcNuhHa-y3xYdmc1x"
)

const (
	// PagerDuty test token only has read-only priviledge.
	//
	// To test delete and remove functionality, needs real token.
	// You can set it via following env vars.
	EnvTestToken      = "TEST_PG_TOKEN"
	EnvTestEmail      = "TEST_PG_EMAIL"
	EnvTestScheduleID = "TEST_PG_SCID"
)

const (
	testEmail  = "taichi.nakashima@dutyme.com"
	testUserID = "PXPGF42"

	testScheduleID1   = "PI7DH85"
	testScheduleName1 = "Dutyme primary"

	testScheduleID2   = "PI9DH21"
	testScheduleName2 = "Dutyme secondary"

	testOverrideID = "PEYSGVF"
)

type testPDClient struct {
}

func (c *testPDClient) GetUser(email string) (*User, error) {
	if email != testEmail {
		return nil, errors.Errorf("user %s doesn't exist", email)
	}
	return &User{
		Email: testEmail,
		Obj: &pagerduty.APIObject{
			ID:      testUserID,
			Type:    "user",
			Summary: "Taichi Nakashima",
			Self:    fmt.Sprintf("https://api.pagerduty.com/users/%s", testUserID),
			HTMLURL: fmt.Sprintf("https://subdomain.pagerduty.com/users/%s", testUserID),
		},
	}, nil
}

func (c *testPDClient) GetSchedules(name string) ([]pagerduty.Schedule, error) {
	schedules := []pagerduty.Schedule{
		{
			APIObject: pagerduty.APIObject{
				ID: testScheduleID1,
			},
			Name: testScheduleName1,
		},
		{
			APIObject: pagerduty.APIObject{
				ID: testScheduleID2,
			},
			Name: testScheduleName2,
		},
	}

	if !strings.Contains(name, "Dutyme") {
		return nil, errors.Errorf("schedule %s doesn't exist", name)
	}

	if name == testScheduleName1 {
		return schedules[:1], nil
	}

	if name == testScheduleName2 {
		return schedules[1:], nil
	}

	return schedules, nil
}
func (c *testPDClient) Override(scheduleID string, user *User, start, end time.Time) (*pagerduty.Override, error) {
	if end.Before(start) {
		return nil, errors.New("end time must be after start time")
	}

	return &pagerduty.Override{
		ID: testOverrideID,
	}, nil
}

func (c *testPDClient) DeleteOverride(scheduleID, overrideID string) error {
	// Currently not used
	return nil
}

func (c *testPDClient) GetOverrides(scheduleID string, since, until time.Time) ([]pagerduty.Override, error) {
	// Currently not used
	return nil, nil
}

func testNewClient(t *testing.T, token string) PagerDuty {
	if len(token) == 0 {
		return &testPDClient{}
	}

	client, err := NewPDClient(token)
	if err != nil {
		t.Fatal("NewClient failed:", err)
	}
	return client
}

func TestFGetUser(t *testing.T) {
	client := testNewClient(t, TestToken)

	email := "cunningham@pagerduty.com"
	user, err := client.GetUser(email)
	if err != nil {
		t.Fatal("GetUser failed:", err)
	}

	if got, want := user.Obj.ID, "PGJ36Z3"; got != want {
		t.Fatalf("GetUser: user.ID = %s;  want %s", got, want)
	}
}

func TestGetSchedules(t *testing.T) {
	client := testNewClient(t, TestToken)

	name := "BoothDuty"
	schedules, err := client.GetSchedules(name)
	if err != nil {
		t.Fatal("GetSchedules failed:", err)
	}

	if got, want := len(schedules), 2; got != want {
		t.Fatalf("GetSchedules number = %d; want %d", got, want)
	}
}

func TestCreateOverride(t *testing.T) {
	token := os.Getenv(EnvTestToken)
	email := os.Getenv(EnvTestEmail)
	scheduleID := os.Getenv(EnvTestScheduleID)

	if len(token) == 0 || len(email) == 0 || len(scheduleID) == 0 {
		t.Skipf(
			"For TestCreateOverride you need real token and email.\nSet them via env vars: %q, %q, %q",
			EnvTestToken, EnvTestEmail, EnvTestScheduleID)
	}

	client := testNewClient(t, token)

	user, err := client.GetUser(email)
	if err != nil {
		t.Fatal("FindUser failed:", err)
	}

	start := time.Now()
	end := start.Add(1 * time.Hour)
	override, err := client.Override(scheduleID, user, start, end)
	if err != nil {
		t.Fatal("Override failed:", err)
	}

	var skipDefer bool
	defer func() {
		if skipDefer {
			return
		}

		if err := client.DeleteOverride(scheduleID, override.ID); err != nil {
			t.Fatal("Delete Override failed:", err)
		}
	}()

	since := time.Now()
	until := since.Add(1 * time.Hour)
	overrides, err := client.GetOverrides(scheduleID, since, until)
	if err != nil {
		t.Fatal("GetOverrides failed:", err)
	}

	if got, want := len(overrides), 1; got != want {
		t.Fatalf("GetOverrides number = %d, want %d", got, want)
	}

	if err := client.DeleteOverride(scheduleID, override.ID); err != nil {
		t.Fatal("Delete Override failed:", err)
	}
	skipDefer = true

	_, err = client.GetOverrides(scheduleID, since, until)
	if !isNotFound(err) {
		t.Fatalf("expect %s to be NotFound error", err)
	}
}

package dutyme

import (
	"os"
	"testing"
	"time"
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

func testNewClient(t *testing.T) PagerDuty {
	client, err := NewPDClient(TestToken)
	if err != nil {
		t.Fatal("NewClient failed:", err)
	}

	return client
}

func TestFGetUser(t *testing.T) {
	client := testNewClient(t)

	email := "cunningham@pagerduty.com"
	user, err := client.GetUser(email)
	if err != nil {
		t.Fatal("GetUser failed:", err)
	}

	if got, want := user.ID, "PGJ36Z3"; got != want {
		t.Fatalf("GetUser: user.ID = %s;  want %s", got, want)
	}
}

func TestGetSchedules(t *testing.T) {
	client := testNewClient(t)

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
			"For TestCreateOverride you must set env vars: %q, %q, %q",
			EnvTestToken, EnvTestEmail, EnvTestScheduleID)
	}
	t.Skip("TODO: Drop this!")

	client, err := NewPDClient(token)
	if err != nil {
		t.Fatal("NewClient failed:", err)
	}

	user, err := client.GetUser(email)
	if err != nil {
		t.Fatal("FindUser failed:", err)
	}

	start := time.Now()
	end := start.Add(1 * time.Hour)
	override, err := client.Override(scheduleID, *user, start, end)
	if err != nil {
		t.Fatal("Override failed:", err)
	}

	if err := client.DeleteOverride(scheduleID, override.ID); err != nil {
		t.Fatal("Delete Override failed:", err)
	}
}

package dutyme

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/tcnksm/go-input"
)

func TestDutyme_GetUser(t *testing.T) {
	email := "cunningham@pagerduty.com"
	d := Dutyme{
		UI: &input.UI{
			Writer: ioutil.Discard,
			Reader: bytes.NewBufferString(email),
		},
		PD: testNewClient(t),
	}

	defaultEmail := ""
	user, err := d.GetUser(defaultEmail)
	if err != nil {
		t.Fatal("GetUser failed:", err)
	}

	if got, want := user.Obj.ID, "PGJ36Z3"; got != want {
		t.Fatalf("GetUser: user.ID = %s;  want %s", got, want)
	}
}

func TestDutyme_GetUser_default(t *testing.T) {
	d := Dutyme{
		UI: &input.UI{
			Writer: ioutil.Discard,
			Reader: bytes.NewBufferString("\n"),
		},
		PD: testNewClient(t),
	}

	defaultEmail := "cunningham@pagerduty.com"
	user, err := d.GetUser(defaultEmail)
	if err != nil {
		t.Fatal("GetUser failed:", err)
	}

	if got, want := user.Obj.ID, "PGJ36Z3"; got != want {
		t.Fatalf("GetUser: user.ID = %s;  want %s", got, want)
	}
}

func TestDutyme_GetSchedule(t *testing.T) {
	d := Dutyme{
		UI: &input.UI{
			Writer: ioutil.Discard,
			Reader: bytes.NewBufferString("BoothDuty\n1\n"),
		},
		PD: testNewClient(t),
	}

	defaultQuery := ""
	name, id, err := d.GetSchedule(defaultQuery)
	if err != nil {
		t.Fatal("GetSchedule failed:", err)
	}

	if want := "PKM7ZY1"; id != want {
		t.Fatalf("GetSchedule ID = %q, want %q", id, want)
	}

	if want := "BoothDuty Primary"; name != want {
		t.Fatalf("GetSchedule name = %q, want %q", name, want)
	}

}

func TestDutyme_GetSchedule_withoutAsking(t *testing.T) {
	d := Dutyme{
		UI: &input.UI{
			Writer: ioutil.Discard,
			Reader: bytes.NewBufferString("BoothDuty Primary\n"),
		},
		PD: testNewClient(t),
	}

	defaultQuery := ""
	name, id, err := d.GetSchedule(defaultQuery)
	if err != nil {
		t.Fatal("GetSchedule failed:", err)
	}

	if want := "PKM7ZY1"; id != want {
		t.Fatalf("GetSchedule ID = %q, want %q", id, want)
	}

	if want := "BoothDuty Primary"; name != want {
		t.Fatalf("GetSchedule name = %q, want %q", name, want)
	}

}

func TestGetOverride(t *testing.T) {
	token := os.Getenv(EnvTestToken)
	email := os.Getenv(EnvTestEmail)
	scheduleID := os.Getenv(EnvTestScheduleID)

	if len(token) == 0 || len(email) == 0 || len(scheduleID) == 0 {
		t.Skipf(
			"For TestCreateOverride you must set env vars: %q, %q, %q",
			EnvTestToken, EnvTestEmail, EnvTestScheduleID)
	}

	shiftDuration := 5 * time.Hour
	client, err := NewPDClient(token)
	if err != nil {
		t.Fatal("NewClient failed:", err)
	}

	d := Dutyme{
		UI: &input.UI{
			Writer: ioutil.Discard,
			Reader: bytes.NewBufferString("1\n"),
		},
		PD: client,
	}

	user, err := client.GetUser(email)
	if err != nil {
		t.Fatal("GetUser failed:", err)
	}

	// Create 2 overrides
	start := time.Now()
	start.Add(shiftDuration)
	end := start.Add(5 * time.Minute)
	override1, err := client.Override(scheduleID, user, start, end)
	if err != nil {
		t.Fatal("Override failed:", err)
	}

	defer func() {
		if err := client.DeleteOverride(scheduleID, override1.ID); err != nil {
			t.Fatal("Delete Override failed:", err)
		}
	}()

	start = time.Now().Add(1 * time.Hour)
	start.Add(shiftDuration)
	end = start.Add(5 * time.Minute)
	override2, err := client.Override(scheduleID, user, start, end)
	if err != nil {
		t.Fatal("Override failed:", err)
	}

	defer func() {
		if err := client.DeleteOverride(scheduleID, override2.ID); err != nil {
			t.Fatal("Delete Override failed:", err)
		}
	}()

	since := time.Now()
	since.Add(shiftDuration)
	until := since.Add(3 * time.Hour)
	overrideID, err := d.GetOverride(scheduleID, user, since, until)
	if err != nil {
		t.Fatal("GetOverride failed:", err)
	}

	if overrideID != override1.ID {
		t.Fatalf("expect %s to be eq %s", overrideID, override1.ID)
	}
}

func TestGetOverride_withoutAsking(t *testing.T) {
	token := os.Getenv(EnvTestToken)
	email := os.Getenv(EnvTestEmail)
	scheduleID := os.Getenv(EnvTestScheduleID)

	if len(token) == 0 || len(email) == 0 || len(scheduleID) == 0 {
		t.Skipf(
			"For TestCreateOverride you must set env vars: %q, %q, %q",
			EnvTestToken, EnvTestEmail, EnvTestScheduleID)
	}

	shiftDuration := 3 * time.Hour
	client, err := NewPDClient(token)
	if err != nil {
		t.Fatal("NewClient failed:", err)
	}

	d := Dutyme{
		UI: &input.UI{
			Writer: ioutil.Discard,
			Reader: bytes.NewBufferString("\n"),
		},
		PD: client,
	}

	user, err := client.GetUser(email)
	if err != nil {
		t.Fatal("GetUser failed:", err)
	}

	// Create 1 overrides
	start := time.Now()
	start.Add(shiftDuration)
	end := start.Add(1 * time.Hour)
	override, err := client.Override(scheduleID, user, start, end)
	if err != nil {
		t.Fatal("Override failed:", err)
	}

	defer func() {
		if err := client.DeleteOverride(scheduleID, override.ID); err != nil {
			t.Fatal("Delete Override failed:", err)
		}
	}()

	since := time.Now()
	since.Add(shiftDuration)
	until := since.Add(1 * time.Hour)
	overrideID, err := d.GetOverride(scheduleID, user, since, until)
	if err != nil {
		t.Fatal("GetOverride failed:", err)
	}

	if overrideID != override.ID {
		t.Fatalf("expect %q to be eq %q", overrideID, override.ID)
	}
}

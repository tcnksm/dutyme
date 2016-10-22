package dutyme

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/tcnksm/go-input"
)

func testNewDutyme(t *testing.T, token, inputStr string) Dutyme {
	return Dutyme{
		UI: &input.UI{
			Writer: ioutil.Discard,
			Reader: bytes.NewBufferString(inputStr),
		},
		PD: testNewClient(t, token),
	}
}

func TestDutyme_GetUser(t *testing.T) {
	d := testNewDutyme(t, "", testEmail)
	user, err := d.GetUser("")
	if err != nil {
		t.Fatal("GetUser failed:", err)
	}

	if got, want := user.Obj.ID, testUserID; got != want {
		t.Fatalf("GetUser: user.ID = %s;  want %s", got, want)
	}
}

func TestDutyme_GetUser_default(t *testing.T) {
	d := testNewDutyme(t, "", "\n")
	user, err := d.GetUser(testEmail)
	if err != nil {
		t.Fatal("GetUser failed:", err)
	}

	if got, want := user.Obj.ID, testUserID; got != want {
		t.Fatalf("GetUser: user.ID = %s;  want %s", got, want)
	}
}

func TestDutyme_GetSchedule(t *testing.T) {
	d := testNewDutyme(t, "", "Dutyme\n1\n")
	name, id, err := d.GetSchedule("")
	if err != nil {
		t.Fatal("GetSchedule failed:", err)
	}

	if want := testScheduleID1; id != want {
		t.Fatalf("GetSchedule ID = %q, want %q", id, want)
	}

	if want := testScheduleName1; name != want {
		t.Fatalf("GetSchedule name = %q, want %q", name, want)
	}
}

func TestDutyme_GetSchedule_withoutAsking(t *testing.T) {
	d := testNewDutyme(t, "", testScheduleName2)
	name, id, err := d.GetSchedule("")
	if err != nil {
		t.Fatal("GetSchedule failed:", err)
	}

	if want := testScheduleID2; id != want {
		t.Fatalf("GetSchedule ID = %q, want %q", id, want)
	}

	if want := testScheduleName2; name != want {
		t.Fatalf("GetSchedule name = %q, want %q", name, want)
	}
}

func TestDutyme_Override(t *testing.T) {
	d := testNewDutyme(t, "", "Y\n")
	start := time.Now()
	end := start.Add(1 * time.Hour)
	_, err := d.Override(testScheduleID1, &User{}, start, end, false)
	if err != nil {
		t.Fatal("Override failed:", err)
	}
}
func TestDutyme_Override_force(t *testing.T) {
	d := testNewDutyme(t, "", "")
	start := time.Now()
	end := start.Add(1 * time.Hour)
	_, err := d.Override(testScheduleID1, &User{}, start, end, true)
	if err != nil {
		t.Fatal("Override failed:", err)
	}
}

func TestDutyme_Override_cancel(t *testing.T) {
	d := testNewDutyme(t, "", "n\n")
	start := time.Now()
	end := start.Add(1 * time.Hour)
	_, err := d.Override(testScheduleID1, &User{}, start, end, false)

	c, ok := err.(*errCancel)
	if !(ok && c.IsCancel()) {
		t.Fatal("Override must be canceled")
	}
}

func TestGetOverride(t *testing.T) {
	token := os.Getenv(EnvTestToken)
	email := os.Getenv(EnvTestEmail)
	scheduleID := os.Getenv(EnvTestScheduleID)

	if len(token) == 0 || len(email) == 0 || len(scheduleID) == 0 {
		t.Skipf(
			"For TestGetOverride you need real token and email.\nSet them via env vars: %q, %q, %q",
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
			"For TestGetOverride you need real token and email.\nSet them via env vars: %q, %q, %q",
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

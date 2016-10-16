package dutyme

import (
	"bytes"
	"io/ioutil"
	"testing"

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
	t.Skip("TODO: need to fix when query includes space")
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

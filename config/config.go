package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/tcnksm/dutyme/dutyme"
)

type Config struct {
	Token string `json:"token,omitempty"`

	User *dutyme.User `json:"user,omitempty"`

	ScheduleName string `json:"schedule_name,omitempty"`
	ScheduleID   string `json:"schedule_id,omitempty"`

	OverrideID         string `json:"override_id,omitempty"`
	OverrideScheduleID string `json:"override_schedule_Id,omitempty"`
}

func (c *Config) IsOverrideExist() bool {
	// TODO(tcnksm): Check by time
	return c.OverrideID != "" && c.ScheduleID != ""
}

func (c *Config) WriteFile(path string, indent bool) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return errors.Wrap(err, "faield to get abs path")
	}

	f, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer f.Close()

	return c.write(f, indent)
}

func (c *Config) write(wr io.Writer, indent bool) error {
	encoder := json.NewEncoder(wr)

	if indent {
		encoder.SetIndent("", "  ")
	}

	if err := encoder.Encode(c); err != nil {
		return errors.Wrap(err, "failed to encode json")
	}

	return nil
}

func ParseFile(path string) (*Config, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrap(err, "faield to get abs path")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}
	defer f.Close()

	return parse(f)
}

func parse(rd io.Reader) (*Config, error) {
	decoder := json.NewDecoder(rd)
	var config Config
	if err := decoder.Decode(&config); err != nil {
		return nil, errors.Wrap(err, "failed to decode json file")
	}

	return &config, nil
}

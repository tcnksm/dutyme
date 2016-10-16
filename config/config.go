package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/pkg/errors"
)

type Config struct {
	Token string `json:"token,omitempty"`

	Email string              `json:"email,omitempty"`
	User  pagerduty.APIObject `json:"user,omitempty"`

	ScheduleName string `json:"schedule_name,omitempty"`
	ScheduleID   string `json:"schedule_id,omitempty"`

	OverrideID         string `json:"override_id,omitempty"`
	OverrideScheduleID string `json:"override_schedule_Id,omitempty"`
}

func WriteFile(path string, config *Config) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return errors.Wrap(err, "faield to get abs path")
	}

	f, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer f.Close()

	return write(f, config)
}

func write(wr io.Writer, config *Config) error {
	encoder := json.NewEncoder(wr)
	if err := encoder.Encode(config); err != nil {
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

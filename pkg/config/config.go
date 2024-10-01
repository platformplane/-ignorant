package config

import (
	"bytes"
	"os"
	"path/filepath"
	"time"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Vulnerabilities   []Vulnerability    `json:"vulnerabilities"`
	Misconfigurations []Misconfiguration `json:"misconfigurations"`
}

type Vulnerability struct {
	ID        string `yaml:"id" validate:"required"`
	Statement string `yaml:"statement" validate:"required"`

	Paths []string `yaml:"paths"`

	ExpiredAt *time.Time `yaml:"expired_at" validate:"required"`
}

type Misconfiguration struct {
	ID        string `yaml:"id" validate:"required"`
	Statement string `yaml:"statement" validate:"required"`

	Paths []string `yaml:"paths"`

	ExpiredAt *time.Time `yaml:"expired_at"`
}

func Parse(path string) (*Config, error) {
	names := []string{
		path,

		filepath.Join(path, "securityignore.yaml"),
		filepath.Join(path, "securityignore.yml"),

		filepath.Join(path, ".securityignore.yaml"),
		filepath.Join(path, ".securityignore.yml"),
	}

	var err error
	var data []byte

	for _, name := range names {
		data, err = os.ReadFile(name)

		if err == nil {
			break
		}
	}

	cfg := new(Config)

	if len(data) == 0 {
		return cfg, nil
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)

	if err := decoder.Decode(cfg); err != nil {
		return nil, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

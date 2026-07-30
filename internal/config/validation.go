package config

import (
	"errors"
	"fmt"
)

var (
	ErrMissingAPIVersion = errors.New("missing apiVersion")
	ErrMissingKind       = errors.New("missing kind")
	ErrEmptyProfiles     = errors.New("profiles must not be empty")
)

type ValidationResult struct {
	Valid  bool
	Errors []string
}

func Validate(path string) (ValidationResult, error) {
	cfg, err := Load(path)
	if err != nil {
		return ValidationResult{}, err
	}

	var result ValidationResult
	if cfg.APIVersion == "" {
		result.Errors = append(result.Errors, ErrMissingAPIVersion.Error())
	}
	if cfg.Kind == "" {
		result.Errors = append(result.Errors, ErrMissingKind.Error())
	}
	if len(cfg.Profiles) == 0 {
		result.Errors = append(result.Errors, ErrEmptyProfiles.Error())
	}

	for name, profile := range cfg.Profiles {
		if len(profile.Providers) == 0 {
			result.Errors = append(result.Errors, fmt.Sprintf("profile %q has no providers", name))
		}
	}

	result.Valid = len(result.Errors) == 0
	return result, nil
}

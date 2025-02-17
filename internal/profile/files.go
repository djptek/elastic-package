// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package profile

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// NewConfig is a generic function type to return a new Managed config
type NewConfig = func(profileName string, profilePath string) (*simpleFile, error)

// simpleFile defines a file that's managed by the profile system
// and doesn't require any  rendering
type simpleFile struct {
	name string
	path string
	body string
}

const profileStackPath = "stack"

// configFilesDiffer checks to see if a local configItem differs from the one it knows.
func (cfg simpleFile) configFilesDiffer() (bool, error) {
	changes, err := os.ReadFile(cfg.path)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrapf(err, "error reading %s", cfg.path)
	}
	if string(changes) == cfg.body {
		return false, nil
	}
	return true, nil
}

// writeConfig writes the config item
func (cfg simpleFile) writeConfig() error {
	err := os.MkdirAll(filepath.Dir(cfg.path), 0755)
	if err != nil {
		return errors.Wrapf(err, "creating parent directories for file failed (path: %s)", cfg.path)
	}
	err = os.WriteFile(cfg.path, []byte(cfg.body), 0644)
	if err != nil {
		return errors.Wrapf(err, "writing file failed (path: %s)", cfg.path)
	}
	return nil
}

// readConfig reads the config item, overwriting whatever exists in the fileBody.
func (cfg *simpleFile) readConfig() error {
	body, err := os.ReadFile(cfg.path)
	if err != nil {
		return errors.Wrapf(err, "reading file failed (path: %s)", cfg.path)
	}
	cfg.body = string(body)
	return nil
}

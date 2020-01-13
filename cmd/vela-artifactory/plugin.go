// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/sirupsen/logrus"
)

// Plugin represents the configuration loaded for the plugin.
type Plugin struct {
	// config arguments loaded for the plugin
	Config *Config
	// copy arguments loaded for the plugin
	Copy *Copy
	// delete arguments loaded for the plugin
	Delete *Delete
	// set-prop arguments loaded for the plugin
	SetProp *SetProp
	// upload arguments loaded for the plugin
	Upload *Upload
}

// Exec formats and runs the commands for managing artifacts in Artifactory.
func (p *Plugin) Exec() error {
	return nil
}

// Validate verifies the plugin is properly configured.
func (p *Plugin) Validate() error {
	logrus.Debug("validating plugin configuration")

	// validate config configuration
	err := p.Config.Validate()
	if err != nil {
		return err
	}

	// validate copy configuration
	err = p.Copy.Validate()
	if err != nil {
		return err
	}

	// validate delete configuration
	err = p.Delete.Validate()
	if err != nil {
		return err
	}

	// validate set-prop configuration
	err = p.SetProp.Validate()
	if err != nil {
		return err
	}

	// validate upload configuration
	err = p.Upload.Validate()
	if err != nil {
		return err
	}

	return nil
}

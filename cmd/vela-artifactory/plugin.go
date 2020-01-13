// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

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

	switch p.Config.Action {
	case copyAction:
		// validate copy configuration
		return p.Copy.Validate()
	case deleteAction:
		// validate delete configuration
		return p.Delete.Validate()
	case setPropAction:
		// validate set-prop configuration
		return p.SetProp.Validate()
	case uploadAction:
		// validate upload configuration
		return p.Upload.Validate()
	default:
		return fmt.Errorf("invalid action provided: %s (Valid actions: %s, %s, %s, %s)", p.Config.Action, copyAction, deleteAction, setPropAction, uploadAction)
	}
}

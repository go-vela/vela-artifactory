// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

var (
	// ErrInvalidAction defines the error type when the
	// Action provided to the Plugin is unsupported.
	ErrInvalidAction = errors.New("invalid action provided")
)

// Plugin represents the configuration loaded for the plugin.
type Plugin struct {
	// Config stores arguments loaded for the plugin
	Config *Config
	// Copy arguments loaded for the plugin
	Copy *Copy
	// Delete arguments loaded for the plugin
	Delete *Delete
	// DockerPromote arguments loaded for the plugin
	DockerPromote *DockerPromote
	// SetProp arguments loaded for the plugin
	SetProp *SetProp
	// Upload arguments loaded for the plugin
	Upload *Upload
}

// Exec formats and runs the commands for managing artifacts in Artifactory.
func (p *Plugin) Exec() error {
	logrus.Debug("running plugin with provided configuration")

	// create new Artifactory client from config configuration
	cli, err := p.Config.New()
	if err != nil {
		return err
	}

	// execute action specific configuration
	switch p.Config.Action {
	case copyAction:
		// execute copy action
		return p.Copy.Exec(*cli)
	case deleteAction:
		// execute delete action
		return p.Delete.Exec(*cli)
	case dockerPromoteAction:
		// execute docker-promote action
		return p.DockerPromote.Exec(*cli)
	case setPropAction:
		// execute set-prop action
		return p.SetProp.Exec(*cli)
	case uploadAction:
		// execute upload action
		return p.Upload.Exec(*cli)
	default:
		return fmt.Errorf(
			"%w: %s (Valid actions: %s, %s, %s, %s, %s)",
			ErrInvalidAction,
			p.Config.Action,
			copyAction,
			deleteAction,
			dockerPromoteAction,
			setPropAction,
			uploadAction,
		)
	}
}

// Validate verifies the plugin is properly configured.
func (p *Plugin) Validate() error {
	logrus.Debug("validating plugin configuration")

	// validate config configuration
	err := p.Config.Validate()
	if err != nil {
		return err
	}

	// validate action specific configuration
	switch p.Config.Action {
	case copyAction:
		// validate copy configuration
		return p.Copy.Validate()
	case deleteAction:
		// validate delete configuration
		return p.Delete.Validate()
	case dockerPromoteAction:
		// validate docker-promote configuration
		return p.DockerPromote.Validate()
	case setPropAction:
		// validate set-prop configuration
		return p.SetProp.Validate()
	case uploadAction:
		// validate upload configuration
		return p.Upload.Validate()
	default:
		return fmt.Errorf(
			"%w: %s (Valid actions: %s, %s, %s, %s)",
			ErrInvalidAction,
			p.Config.Action,
			copyAction,
			deleteAction,
			setPropAction,
			uploadAction,
		)
	}
}

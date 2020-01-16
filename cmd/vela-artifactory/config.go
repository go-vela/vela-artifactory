// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/utils/log"

	"github.com/sirupsen/logrus"
)

// Config represents the plugin configuration for Artifactory config information.
type Config struct {
	// Action to perform against the Artifactory instance
	Action string
	// API key for communication with the Artifactory instance
	APIKey string
	// password for communication with the Artifactory instance
	Password string
	// full url to Artifactory instance
	URL string
	// user name for communication with the Artifactory instance
	Username string
}

// New creates an Artifactory client for managing artifacts.
func (c *Config) New() (*artifactory.ArtifactoryServicesManager, error) {
	// create new Artifactory details
	details := auth.NewArtifactoryDetails()

	// check if provided URL doesn't have a "/" suffix
	if !strings.HasSuffix(c.URL, "/") {
		// add a "/" to the end of the provided URL
		c.URL = c.URL + "/"
	}

	// set URL for Artifactory details
	details.SetUrl(c.URL)
	// set user name for Artifactory details
	details.SetUser(c.Username)

	// check if API key is provided
	if len(c.APIKey) > 0 {
		// set API key for Artifactory details
		details.SetApiKey(c.APIKey)
	}

	// check if password is provided
	if len(c.Password) > 0 {
		// set password for Artifactory details
		details.SetPassword(c.Password)
	}

	// set logger for Artifactory client
	log.SetLogger(
		// create new logger for Artifactory client
		log.NewLogger(log.INFO, os.Stdout),
	)

	// create new Artifactory config from details
	config, err := artifactory.NewConfigBuilder().
		SetArtDetails(details).
		Build()
	if err != nil {
		return nil, err
	}

	// create new Artifactory client from config and details
	client, err := artifactory.New(&details, config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Validate verifies the Config is properly configured.
func (c *Config) Validate() error {
	logrus.Trace("validating config plugin configuration")

	if len(c.Action) == 0 {
		return fmt.Errorf("no config action provided")
	}

	if len(c.URL) == 0 {
		return fmt.Errorf("no config url provided")
	}

	if len(c.Username) == 0 {
		return fmt.Errorf("no config username provided")
	}

	if len(c.Password) == 0 && len(c.APIKey) == 0 {
		return fmt.Errorf("no config api-key or password provided")
	}

	return nil
}

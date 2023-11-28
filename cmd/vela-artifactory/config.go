// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/httpclient"
	"github.com/jfrog/jfrog-client-go/utils/log"

	"github.com/sirupsen/logrus"
)

// Config represents the plugin configuration for Artifactory config information.
type Config struct {
	// Action to perform against the Artifactory instance
	Action string
	// Token for communication with the Artifactory instance
	Token string
	// API key for communication with the Artifactory instance
	APIKey string
	// DryRun enables pretending to perform the action against the Artifactory instance
	DryRun bool
	// Password for communication with the Artifactory instance
	Password string
	// URL points to the Artifactory instance
	URL string
	// Username for communication with the Artifactory instance
	Username string
}

// New creates an Artifactory client for managing artifacts.
func (c *Config) New() (*artifactory.ArtifactoryServicesManager, error) {
	logrus.Trace("creating new Artifactory client from plugin configuration")

	// create new Artifactory details
	details := auth.NewArtifactoryDetails()

	// check if provided URL doesn't have a "/" suffix
	if !strings.HasSuffix(c.URL, "/") {
		// add a "/" to the end of the provided URL
		c.URL = c.URL + "/"
	}

	// set URL for Artifactory details
	details.SetUrl(c.URL)
	// check if username is provided
	if len(c.Username) > 0 {
		// set user name for Artifactory details
		details.SetUser(c.Username)
	}

	// check if Access/Identity token is provided
	if len(c.Token) > 0 && !httpclient.IsApiKey(c.Token) {
		// set Access/Identity token for Artifactory details
		details.SetAccessToken(c.APIKey)
	} else if len(c.APIKey) > 0 && httpclient.IsApiKey(c.APIKey) { // check if API key is provided
		// set API key for Artifactory details
		details.SetApiKey(c.APIKey)
	} else if len(c.APIKey) > 0 && !httpclient.IsApiKey(c.APIKey) {
		// set Access/Identity token for Artifactory details
		details.SetAccessToken(c.APIKey)
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
	config, err := config.NewConfigBuilder().
		SetServiceDetails(details).
		SetDryRun(c.DryRun).
		Build()
	if err != nil {
		return nil, err
	}

	// create new Artifactory client from config and details
	client, err := artifactory.New(config)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

// Validate verifies the Config is properly configured.
func (c *Config) Validate() error {
	logrus.Trace("validating config plugin configuration")

	// verify action is provided
	if len(c.Action) == 0 {
		return fmt.Errorf("no config action provided")
	}

	// verify URL is provided
	if len(c.URL) == 0 {
		return fmt.Errorf("no config url provided")
	}

	// check if dry run is disabled
	if !c.DryRun {
		// verify username is provided if APIKey is empty
		if len(c.Username) == 0 && len(c.APIKey) == 0 {
			return fmt.Errorf("no config username provided")
		}

		// verify API key or password are provided
		if len(c.APIKey) == 0 && len(c.Password) == 0 {
			return fmt.Errorf("no config api-key or password provided")
		}
	}

	return nil
}

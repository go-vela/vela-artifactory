// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Config represents the plugin configuration for Artifactory config information.
type Config struct {
	// API key for communication with the Artifactory instance
	APIKey string
	// password for communication with the Artifactory instance
	Password string
	// full url to Artifactory instance
	URL string
	// user name for communication with the Artifactory instance
	Username string
}

// Validate verifies the Config is properly configured.
func (c Config) Validate() error {
	logrus.Trace("validating config plugin configuration")

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

// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"time"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"

	"github.com/sirupsen/logrus"
)

const copyAction = "copy"

// Copy represents the plugin configuration for copy information.
type Copy struct {
	// Flat is a flag that enables removing source file directory hierarchy
	Flat bool
	// Recursive is a flag that enables copying sub-directories from source
	Recursive bool
	// Path is the source path to artifact(s) to copy
	Path string
	// Target is the path to copy artifact(s) to
	Target string
}

// Exec formats and runs the commands for copying artifacts in Artifactory.
func (c *Copy) Exec(cli artifactory.ArtifactoryServicesManager) error {
	logrus.Trace("running copy with provided configuration")

	// create new copy parameters
	p := services.NewMoveCopyParams()

	// add copy configuration to copy parameters
	p.CommonParams = &utils.CommonParams{
		Pattern:   c.Path,
		Recursive: c.Recursive,
		Target:    c.Target,
	}
	p.Flat = c.Flat

	retries := 3
	var retryErr error

	// send API call to copy artifacts in Artifactory
	// retry a couple times to avoid intermittent authentication failures from Artifactory
	for i := 0; i < retries; i++ {
		backoff := i*2 + 1

		// send API call to copy artifacts in Artifactory
		_, _, retryErr = cli.Copy(p)
		if retryErr != nil {
			logrus.Errorf("Error copying files using pattern %s on attempt %d: %s",
				p.CommonParams.Pattern,
				i+1, retryErr.Error())

			if i == retries-1 {
				return retryErr
			}

			logrus.Infof("Retrying in %d seconds...", backoff)

			time.Sleep(time.Duration(backoff) * time.Second)
		} else {
			break
		}
	}

	return nil
}

// Validate verifies the Copy is properly configured.
func (c *Copy) Validate() error {
	logrus.Trace("validating copy plugin configuration")

	// verify path is provided
	if len(c.Path) == 0 {
		return fmt.Errorf("no copy path provided")
	}

	// verify target is provided
	if len(c.Target) == 0 {
		return fmt.Errorf("no copy target provided")
	}

	return nil
}

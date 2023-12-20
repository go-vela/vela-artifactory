// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
)

const uploadAction = "upload"

// Upload represents the plugin configuration for upload information.
type Upload struct {
	// enables uploading artifacts to exact target path (excludes source file hierarchy)
	Flat bool
	// enables including directories and files from sources
	IncludeDirs bool
	// enables uploading sub-directories for the sources
	Recursive bool
	// enables reading the sources as a regular expression
	Regexp bool
	// target path to upload artifact(s) to
	Path string
	// list of files to upload
	Sources []string
}

// Exec formats and runs the commands for uploading artifacts in Artifactory.
func (u *Upload) Exec(cli artifactory.ArtifactoryServicesManager) error {
	logrus.Trace("running upload with provided configuration")

	// iterate through all sources
	for _, source := range u.Sources {
		// create new upload parameters
		p := services.NewUploadParams()

		// add upload configuration to upload parameters
		p.CommonParams = &utils.CommonParams{
			IncludeDirs: u.IncludeDirs,
			Pattern:     source,
			Recursive:   u.Recursive,
			Regexp:      u.Regexp,
			Target:      u.Path,
		}
		p.Flat = u.Flat

		retries := 3
		var retryErr error
		totalFailed := 0
		totalSucceeded := 0
		maxSucceeded := 0

		// send API call to upload files in Artifactory
		// retry a couple times to avoid intermittent from Artifactory
		for i := 0; i < retries; i++ {
			backoff := i*2 + 1

			// send API call to upload artifacts in Artifactory
			totalSucceeded, totalFailed, retryErr = cli.UploadFiles(p)

			if maxSucceeded < totalSucceeded {
				maxSucceeded = totalSucceeded
			}

			if totalFailed > 0 || retryErr != nil {
				logrus.Errorf("Error uploading files to target path %s on attempt %d: %s", p.CommonParams.Target, i+1, retryErr.Error())

				if i == retries-1 {
					if maxSucceeded > 0 {
						logrus.Warn("Successfully uploaded some files. Therefore conflicts may have occurred during re-upload attempts. Check Artifactory.")
					}

					if retryErr == nil {
						retryErr = fmt.Errorf("failed to upload %d files", totalFailed)
					}

					return retryErr
				}

				logrus.Infof("Retrying in %d seconds...", backoff)

				time.Sleep(time.Duration(backoff) * time.Second)
			} else {
				break
			}
		}
	}

	return nil
}

// Validate verifies the Upload is properly configured.
func (u *Upload) Validate() error {
	logrus.Trace("validating upload plugin configuration")

	// verify path is provided
	if len(u.Path) == 0 {
		return fmt.Errorf("no upload path provided")
	}

	// verify sources are provided
	if len(u.Sources) == 0 {
		return fmt.Errorf("no upload sources provided")
	}

	return nil
}

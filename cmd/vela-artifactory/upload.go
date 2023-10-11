// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"

	"github.com/sirupsen/logrus"

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
func (u *Upload) Exec(cli ArtifactoryServicesManager) error {
	logrus.Trace("running upload with provided configuration")

	// iterate through all sources
	for _, source := range u.Sources {
		// create new upload parameters
		p := services.NewUploadParams()

		// add upload configuration to upload parameters
		p.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{
			IncludeDirs: u.IncludeDirs,
			Pattern:     source,
			Recursive:   u.Recursive,
			Regexp:      u.Regexp,
			Target:      u.Path,
		}
		p.Flat = u.Flat
		p.Retries = 3

		// send API call to upload artifacts in Artifactory
		_, _, totalFailed, err := cli.UploadFiles(p)
		if totalFailed > 0 || err != nil {
			return err
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

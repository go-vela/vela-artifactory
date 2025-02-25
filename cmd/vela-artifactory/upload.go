// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"strings"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/sirupsen/logrus"
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
	// build props
	BuildProps string
	// list of files to upload
	Sources []string
}

// Exec formats and runs the commands for uploading artifacts in Artifactory.
func (u *Upload) Exec(cli artifactory.ArtifactoryServicesManager) error {
	logrus.Trace("running upload with provided configuration")

	// very simple check that doesn't account for:
	// - regex in sources, in which case it's possible that one source could be multiple files
	// - defining a singular source twice
	if len(u.Sources) > 1 && !strings.HasSuffix(u.Path, "/") {
		logrus.Warn("when uploading multiple sources, path should be a directory")
	}

	// iterate through all sources
	for _, source := range u.Sources {
		// create new upload parameters
		p := services.NewUploadParams()

		// apply build props
		p.BuildProps = u.BuildProps

		// add upload configuration to upload parameters
		p.CommonParams = &utils.CommonParams{
			IncludeDirs: u.IncludeDirs,
			Pattern:     source,
			Recursive:   u.Recursive,
			Regexp:      u.Regexp,
			Target:      u.Path,
		}

		// upload to exact target path
		p.Flat = u.Flat

		// send API call to upload artifacts in Artifactory
		_, totalFailed, err := cli.UploadFiles(artifactory.UploadServiceOptions{FailFast: true}, p)
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

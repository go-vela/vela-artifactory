// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

const uploadAction = "upload"

// Upload represents the plugin configuration for upload information.
type Upload struct {
	// file path to load arguments from
	ArgsFile string
	// enables pretending to upload the artifact(s)
	DryRun bool
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

// Validate verifies the Upload is properly configured.
func (u *Upload) Validate() error {
	logrus.Trace("validating upload plugin configuration")

	if len(u.Path) == 0 {
		return fmt.Errorf("no upload path provided")
	}

	if len(u.Sources) == 0 {
		return fmt.Errorf("no upload sources provided")
	}

	return nil
}

// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Delete represents the plugin configuration for delete information.
type Delete struct {
	// file path to load arguments from
	ArgsFile string
	// enables pretending to remove the artifact(s) in the path
	DryRun bool
	// enables removing sub-directories for the artifact(s) in the path
	Recursive bool
	// target path to artifact(s) to remove
	Path string
}

// Validate verifies the Delete is properly configured.
func (d *Delete) Validate() error {
	logrus.Trace("validating delete plugin configuration")

	if len(d.Path) == 0 {
		return fmt.Errorf("no delete path provided")
	}

	return nil
}

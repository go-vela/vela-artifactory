// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

const copyAction = "copy"

// Copy represents the plugin configuration for copy information.
type Copy struct {
	// enables removing source file directory hierarchy
	Flat bool
	// enables copying sub-directories from source
	Recursive bool
	// source path to artifact(s) to copy
	Path string
	// target path to copy artifact(s) to
	Target string
}

// Validate verifies the Copy is properly configured.
func (c *Copy) Validate() error {
	logrus.Trace("validating copy plugin configuration")

	if len(c.Path) == 0 {
		return fmt.Errorf("no copy path provided")
	}

	if len(c.Target) == 0 {
		return fmt.Errorf("no copy target provided")
	}

	return nil
}

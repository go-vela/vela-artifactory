// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Copy represents the plugin configuration for copy information.
type Copy struct {
	// enables removing source file directory hierarchy
	Flat bool
	// enables copying sub-directories from source
	Recursive bool
	// file source to copy from
	Source string
	// file target to copy source file to
	Target string
}

// Validate verifies the Copy is properly configured.
func (c *Copy) Validate() error {
	logrus.Trace("validating copy plugin configuration")

	if len(c.Source) == 0 {
		return fmt.Errorf("no copy source provided")
	}

	if len(c.Target) == 0 {
		return fmt.Errorf("no copy target provided")
	}

	return nil
}

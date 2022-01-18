// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package fixtures

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler returns an http.Handler that is capable of handling a variety of mock
// Artifactory requests and returning mock responses.
func Handler() http.Handler {
	gin.SetMode(gin.TestMode)

	e := gin.New()
	e.GET("/the-repo-key/project", getRepo)
	e.GET("/the-repo-key/project/filename.tar2", getFile)

	return e
}

func getRepo(c *gin.Context) {
	c.String(200, "")
}

func getFile(c *gin.Context) {
	c.String(200, "")
}

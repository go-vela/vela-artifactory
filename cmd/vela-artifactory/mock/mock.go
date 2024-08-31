// SPDX-License-Identifier: Apache-2.0

package mock

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	InvalidArtifactoryServerURL = "http://localhost:8081/artifactory"
	Username                    = "octocat"
	APIKey                      = "superSecretAPIKey"
	Password                    = "superSecretPassword"
	Token                       = "superSecretToken"
)

// Handlers returns an http.Handler group that is capable of handling
// Artifactory API requests and returning mock responses.
func Handlers() *gin.Engine {
	gin.SetMode(gin.TestMode)

	e := gin.New()

	e.GET("/api/system/version", getVersion)
	e.POST("/api/search/aql", search)
	e.POST("/api/copy", copyArtifact)
	e.DELETE("/", deleteArtifact)
	e.GET("/api/docker/:registry/v2/_catalog", getRepositories)
	e.GET("/api/docker/:registry/v2/docker-dev/tags/list", getTags)
	e.POST("/api/docker/:registry/v2/promote", promoteImage)
	e.PUT("/api/storage", setProp)
	e.PUT("/foo/bar", uploadFiles)

	return e
}

func getVersion(c *gin.Context) {
	c.String(200, loadFixture("mock/fixtures/version.json"))
}

func search(c *gin.Context) {
	c.String(200, loadFixture("mock/fixtures/search.json"))
}

func copyArtifact(c *gin.Context) {
	c.JSON(200, "Copy ended successfully")
}

func deleteArtifact(c *gin.Context) {
	c.JSON(204, "Delete ended successfully")
}

func setProp(c *gin.Context) {
	c.JSON(204, "Property set successfully")
}

func uploadFiles(c *gin.Context) {
	c.JSON(200, map[string]interface{}{
		"checksums": map[string]interface{}{"checksum": "abcxyz123"},
	})
}

func getRepositories(c *gin.Context) {
	registry := c.Param("registry")

	if strings.Contains(registry, "not-found") {
		c.JSON(404, fmt.Sprintf("Registry %s does not exist", registry))
		return
	}

	c.String(200, loadFixture("mock/fixtures/repositories.json"))
}

func getTags(c *gin.Context) {
	registry := c.Param("registry")

	if strings.Contains(registry, "not-found") {
		c.JSON(404, fmt.Sprintf("Registry %s does not exist", registry))
		return
	}

	c.String(200, loadFixture("mock/fixtures/tags.json"))
}

func promoteImage(c *gin.Context) {
	registry := c.Param("registry")

	if strings.Contains(registry, "not-found") {
		c.JSON(404, fmt.Sprintf("Registry %s does not exist", registry))
		return
	}

	c.JSON(200, "Promotion ended successfully")
}

func loadFixture(file string) string {
	data, _ := os.ReadFile(file)
	return string(data)
}

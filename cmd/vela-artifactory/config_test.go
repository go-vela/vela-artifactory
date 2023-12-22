// SPDX-License-Identifier: Apache-2.0

package main

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/vela-artifactory/cmd/vela-artifactory/mock"
)

func TestArtifactory_Config_New(t *testing.T) {
	// setup types
	c := &Config{
		Action:   "copy",
		APIKey:   mock.APIKey,
		DryRun:   false,
		URL:      mock.InvalidArtifactoryServerURL,
		Username: mock.Username,
		Password: mock.Password,
		Client: &Client{
			Retries:            3,
			RetryWaitMilliSecs: 1,
		},
	}

	got, err := c.New()
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if got == nil {
		t.Errorf("New is nil")
	}
}

func TestArtifactory_Config_Validate(t *testing.T) {
	// setup types
	c := &Config{
		Action:   "copy",
		APIKey:   mock.APIKey,
		DryRun:   false,
		URL:      mock.InvalidArtifactoryServerURL,
		Username: mock.Username,
		Password: mock.Password,
	}

	err := c.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_Config_Validate_NoAction(t *testing.T) {
	// setup types
	c := &Config{
		APIKey:   mock.APIKey,
		DryRun:   false,
		URL:      mock.InvalidArtifactoryServerURL,
		Username: mock.Username,
		Password: mock.Password,
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Config_Validate_NoAPIKeyOrPassword(t *testing.T) {
	// setup types
	c := &Config{
		Action:   "copy",
		DryRun:   false,
		URL:      mock.InvalidArtifactoryServerURL,
		Username: mock.Username,
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Config_Validate_NoAPIKey(t *testing.T) {
	// setup types
	c := &Config{
		Action:   "copy",
		DryRun:   false,
		Password: mock.Password,
		URL:      mock.InvalidArtifactoryServerURL,
		Username: mock.Username,
	}

	err := c.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_Config_Validate_NoPassword(t *testing.T) {
	// setup types
	c := &Config{
		Action:   "copy",
		DryRun:   false,
		APIKey:   mock.APIKey,
		URL:      mock.InvalidArtifactoryServerURL,
		Username: mock.Username,
	}

	err := c.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_Config_Validate_NoUrl(t *testing.T) {
	// setup types
	c := &Config{
		Action:   "copy",
		DryRun:   false,
		APIKey:   mock.APIKey,
		Username: mock.Username,
		Password: mock.Password,
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Config_Validate_NoUsername(t *testing.T) {
	// setup types
	c := &Config{
		Action:   "copy",
		DryRun:   false,
		Password: mock.Password,
		URL:      mock.InvalidArtifactoryServerURL,
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Config_Validate_NoAuth(t *testing.T) {
	// setup types
	c := &Config{
		Action: "copy",
		DryRun: false,
		URL:    mock.InvalidArtifactoryServerURL,
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Config_500_Retries(t *testing.T) {
	searchAttempts := 0
	copyAttempts := 0

	httpRetries := 5
	httpRetryWaitMilliSecs := 2

	failureResponseCode := 500

	e := gin.New()

	search := func(c *gin.Context) {
		searchAttempts++

		data, _ := os.ReadFile("mock/fixtures/search.json")

		if searchAttempts == httpRetries+1 {
			c.String(200, string(data))
		} else {
			c.Status(failureResponseCode)
		}
	}

	copyArtifact := func(c *gin.Context) {
		copyAttempts++

		c.Status(failureResponseCode)
	}

	e.POST("/api/search/aql", search)
	e.POST("/api/copy", copyArtifact)

	ss := httptest.NewServer(e)

	// setup types
	p := &Plugin{
		Config: &Config{
			Action:   "copy",
			Token:    mock.Token,
			APIKey:   mock.APIKey,
			DryRun:   false,
			URL:      ss.URL,
			Username: mock.Username,
			Password: mock.Password,
			Client: &Client{
				Retries:            httpRetries,
				RetryWaitMilliSecs: httpRetryWaitMilliSecs,
			},
		},
		Copy: &Copy{
			Flat:      false,
			Recursive: false,
			Path:      "foo/bar",
			Target:    "bar/foo",
		},
		Delete:  &Delete{},
		SetProp: &SetProp{},
		Upload:  &Upload{},
	}

	err := p.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}

	if searchAttempts != httpRetries+1 {
		t.Errorf("client should have attempted to search %d times, only attempted %d times", httpRetries+1, searchAttempts)
	}

	if copyAttempts != httpRetries+1 {
		t.Errorf("client should have attempted to copy %d times, only attempted %d times", httpRetries+1, copyAttempts)
	}
}

func TestArtifactory_Config_403_Retries(t *testing.T) {
	searchAttempts := 0
	copyAttempts := 0

	httpRetries := 3
	httpRetryWaitMilliSecs := 1

	failureResponseCode := 403

	e := gin.New()

	search := func(c *gin.Context) {
		searchAttempts++

		data, _ := os.ReadFile("mock/fixtures/search.json")

		if searchAttempts == httpRetries+1 {
			c.String(200, string(data))
		} else {
			c.Status(failureResponseCode)
		}
	}

	copyArtifact := func(c *gin.Context) {
		copyAttempts++

		c.Status(failureResponseCode)
	}

	e.POST("/api/search/aql", search)
	e.POST("/api/copy", copyArtifact)

	ss := httptest.NewServer(e)

	// setup types
	p := &Plugin{
		Config: &Config{
			Action:   "copy",
			Token:    mock.Token,
			APIKey:   mock.APIKey,
			DryRun:   false,
			URL:      ss.URL,
			Username: mock.Username,
			Password: mock.Password,
			Client: &Client{
				Retries:            httpRetries,
				RetryWaitMilliSecs: httpRetryWaitMilliSecs,
			},
		},
		Copy: &Copy{
			Flat:      false,
			Recursive: false,
			Path:      "foo/bar",
			Target:    "bar/foo",
		},
		Delete:  &Delete{},
		SetProp: &SetProp{},
		Upload:  &Upload{},
	}

	err := p.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}

	if searchAttempts != httpRetries+1 {
		t.Errorf("client should have attempted to search %d times, only attempted %d times", httpRetries+1, searchAttempts)
	}

	if copyAttempts != httpRetries+1 {
		t.Errorf("client should have attempted to copy %d times, only attempted %d times", httpRetries+1, copyAttempts)
	}
}

func TestArtifactory_Config_401_Retries(t *testing.T) {
	searchAttempts := 0
	copyAttempts := 0

	httpRetries := 3
	httpRetryWaitMilliSecs := 1

	failureResponseCode := 401

	e := gin.New()

	search := func(c *gin.Context) {
		searchAttempts++

		data, _ := os.ReadFile("mock/fixtures/search.json")

		if searchAttempts == httpRetries+1 {
			c.String(200, string(data))
		} else {
			c.Status(failureResponseCode)
		}
	}

	copyArtifact := func(c *gin.Context) {
		copyAttempts++

		c.Status(failureResponseCode)
	}

	e.POST("/api/search/aql", search)
	e.POST("/api/copy", copyArtifact)

	ss := httptest.NewServer(e)

	// setup types
	p := &Plugin{
		Config: &Config{
			Action:   "copy",
			Token:    mock.Token,
			APIKey:   mock.APIKey,
			DryRun:   false,
			URL:      ss.URL,
			Username: mock.Username,
			Password: mock.Password,
			Client: &Client{
				Retries:            httpRetries,
				RetryWaitMilliSecs: httpRetryWaitMilliSecs,
			},
		},
		Copy: &Copy{
			Flat:      false,
			Recursive: false,
			Path:      "foo/bar",
			Target:    "bar/foo",
		},
		Delete:  &Delete{},
		SetProp: &SetProp{},
		Upload:  &Upload{},
	}

	err := p.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}

	if searchAttempts != 1 {
		t.Errorf("client should have attempted to search %d times, only attempted %d times", httpRetries+1, searchAttempts)
	}

	if copyAttempts != 0 {
		t.Errorf("client should have attempted to copy %d times, only attempted %d times", httpRetries+1, copyAttempts)
	}
}

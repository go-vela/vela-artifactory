// SPDX-License-Identifier: Apache-2.0

package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/auth/cert"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/http/httpclient"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/sirupsen/logrus"
)

// Config represents the plugin configuration for Artifactory config information.
type Config struct {
	// Action to perform against the Artifactory instance
	Action string
	// Token for communication with the Artifactory instance
	Token string
	// API key for communication with the Artifactory instance
	APIKey string
	// DryRun enables pretending to perform the action against the Artifactory instance
	DryRun bool
	// Username for communication with the Artifactory instance
	Username string
	// Password for communication with the Artifactory instance
	Password string
	// URL points to the Artifactory instance
	URL string
	// Client represents the HTTP client configurations for interacting with Artifactory
	*Client
}

// Client represents the HTTP client configurations for interacting with Artifactory.
type Client struct {
	// Retries is the number of times to retry communication with Artifactory
	Retries int
	// RetryWaitMilliSecs is the number of milliseconds to wait between retries
	RetryWaitMilliSecs int
	// CertPath is the path to the certificate for communication with Artifactory
	CertPath string
	// CertKeyPath is the path to the certificate key for communication with Artifactory
	CertKeyPath string
	// InsecureTLS enables insecure TLS communication with Artifactory
	InsecureTLS bool
}

// New creates an Artifactory client for managing artifacts.
func (c *Config) New() (*artifactory.ArtifactoryServicesManager, error) {
	logrus.Trace("creating new Artifactory client from plugin configuration")

	// create new Artifactory details
	details := auth.NewArtifactoryDetails()

	// check if provided URL doesn't have a "/" suffix
	if !strings.HasSuffix(c.URL, "/") {
		// add a "/" to the end of the provided URL
		c.URL = c.URL + "/"
	}

	// set URL for Artifactory details
	details.SetUrl(c.URL)
	// check if username is provided
	if len(c.Username) > 0 {
		// set user name for Artifactory details
		details.SetUser(c.Username)
	}

	// check if Access/Identity token is provided
	if len(c.Token) > 0 && !httpclient.IsApiKey(c.Token) {
		// set Access/Identity token for Artifactory details
		details.SetAccessToken(c.APIKey)
	} else if len(c.APIKey) > 0 && httpclient.IsApiKey(c.APIKey) { // check if API key is provided
		// set API key for Artifactory details
		details.SetApiKey(c.APIKey)
	} else if len(c.APIKey) > 0 && !httpclient.IsApiKey(c.APIKey) {
		// set Access/Identity token for Artifactory details
		details.SetAccessToken(c.APIKey)
	}

	// check if password is provided
	if len(c.Password) > 0 {
		// set password for Artifactory details
		details.SetPassword(c.Password)
	}

	// set logger for Artifactory client
	log.SetLogger(
		// create new logger for Artifactory client
		log.NewLogger(log.INFO, os.Stdout),
	)

	// set default HTTP retries and retry wait in milliseconds
	if c.Retries < 0 {
		logrus.Warn("invalid retry count provided, defaulting to 3")

		c.Retries = 3
	}

	if c.RetryWaitMilliSecs < 0 {
		logrus.Warn("invalid retry wait in milliseconds provided, defaulting to 500")

		c.RetryWaitMilliSecs = 500
	}

	// create a retryable http client using https://github.com/hashicorp/go-retryablehttp
	// which is mostly a wrapper around https://golang.org/pkg/net/http/#Client with a customized
	// roundtrip transport that can be modified to retry certain http responses
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = c.Retries
	retryClient.RetryWaitMin = time.Millisecond * time.Duration(c.RetryWaitMilliSecs)

	transport := cleanhttp.DefaultPooledTransport()

	//nolint:gosec // disable warning for InsecureSkipVerify
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: c.InsecureTLS}

	// read certs if provided a path to a certificate and key
	if c.CertPath != "" && c.CertKeyPath != "" {
		certificate, err := cert.LoadCertificate(c.CertPath, c.CertKeyPath)
		if err != nil {
			return nil, err
		}

		transport.TLSClientConfig.Certificates = []tls.Certificate{certificate}
	}

	retryClient.HTTPClient.Transport = transport

	// apply a custom retry policy that has been modified to retry on 403 responses
	// using an incrementing backoff
	retryClient.CheckRetry = RetryPolicy

	// create new Artifactory config from details
	config, err := config.NewConfigBuilder().
		SetServiceDetails(details).
		SetDryRun(c.DryRun).
		// disable the default jfrog client retry policy
		SetHttpRetryWaitMilliSecs(0).
		SetHttpRetries(0).
		// enable the custom retry policy using retryablehttp
		SetHttpClient(retryClient.StandardClient()).
		Build()
	if err != nil {
		return nil, err
	}

	// create new Artifactory client from config and details
	client, err := artifactory.New(config)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

// Validate verifies the Config is properly configured.
func (c *Config) Validate() error {
	logrus.Trace("validating config plugin configuration")

	// verify action is provided
	if len(c.Action) == 0 {
		return fmt.Errorf("no config action provided")
	}

	// verify URL is provided
	if len(c.URL) == 0 {
		return fmt.Errorf("no config url provided")
	}

	// check if dry run is disabled
	if !c.DryRun {
		// verify username is provided if APIKey is empty
		if len(c.Username) == 0 && len(c.APIKey) == 0 {
			return fmt.Errorf("no config username provided")
		}

		// verify API key or password are provided
		if len(c.APIKey) == 0 && len(c.Password) == 0 {
			return fmt.Errorf("no config api-key or password provided")
		}
	}

	return nil
}

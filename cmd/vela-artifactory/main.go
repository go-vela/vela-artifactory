// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/mail"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	_ "github.com/joho/godotenv/autoload"

	"github.com/go-vela/vela-artifactory/version"
)

//nolint:funlen // ignore function length due to comments and flags
func main() {
	// capture application version information
	v := version.New()

	// serialize the version information as pretty JSON
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logrus.Fatal(err)
	}

	// output the version information to stdout
	fmt.Fprintf(os.Stdout, "%s\n", string(bytes))

	// create new CLI application
	app := &cli.Command{

		// Plugin Information

		Name:      "vela-artifactory",
		Usage:     "Vela Artifactory plugin for managing artifacts",
		Copyright: "Copyright 2019 Target Brands, Inc. All rights reserved.",
		Authors: []any{
			&mail.Address{
				Name:    "Vela Admins",
				Address: "vela@target.com",
			},
		},

		// Plugin Metadata

		Action:  run,
		Version: v.Semantic(),

		// Plugin Flags

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "log.level",
				Value: "info",
				Usage: "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_LOG_LEVEL"),
					cli.EnvVar("ARTIFACTORY_LOG_LEVEL"),
					cli.File("/vela/parameters/artifactory/log_level"),
					cli.File("/vela/secrets/artifactory/log_level"),
				),
			},
			&cli.StringFlag{
				Name:  "path",
				Usage: "source/target path to artifact(s) for action",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_PATH"),
					cli.EnvVar("ARTIFACTORY_PATH"),
					cli.File("/vela/parameters/artifactory/path"),
					cli.File("/vela/secrets/artifactory/path"),
				),
			},
			&cli.BoolFlag{
				Name:  "recursive",
				Usage: "enables operating on sub-directories for the source/target path",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_RECURSIVE"),
					cli.EnvVar("ARTIFACTORY_RECURSIVE"),
					cli.File("/vela/parameters/artifactory/recursive"),
					cli.File("/vela/secrets/artifactory/recursive"),
				),
			},

			// Config Flags

			&cli.StringFlag{
				Name:  "config.action",
				Usage: "action to perform against the Artifactory instance",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_ACTION"),
					cli.EnvVar("ARTIFACTORY_ACTION"),
					cli.File("/vela/parameters/artifactory/action"),
					cli.File("/vela/secrets/artifactory/action"),
				),
			},
			&cli.BoolFlag{
				Name:  "config.dry_run",
				Usage: "enables pretending to perform the action",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_DRY_RUN"),
					cli.EnvVar("ARTIFACTORY_DRY_RUN"),
					cli.File("/vela/parameters/artifactory/dry_run"),
					cli.File("/vela/secrets/artifactory/dry_run"),
				),
			},
			&cli.StringFlag{
				Name:  "config.api_key",
				Usage: "API key for communication with the Artifactory instance",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_API_KEY"),
					cli.EnvVar("ARTIFACTORY_API_KEY"),
					cli.File("/vela/parameters/artifactory/api_key"),
					cli.File("/vela/secrets/artifactory/api_key"),
				),
			},
			&cli.StringFlag{
				Name:  "config.token",
				Usage: "Access/Identity token for communication with the Artifactory instance",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_TOKEN"),
					cli.EnvVar("ARTIFACTORY_TOKEN"),
					cli.File("/vela/parameters/artifactory/token"),
					cli.File("/vela/secrets/artifactory/token"),
				),
			},
			&cli.StringFlag{
				Name:  "config.password",
				Usage: "password for communication with the Artifactory instance",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_PASSWORD"),
					cli.EnvVar("ARTIFACTORY_PASSWORD"),
					cli.File("/vela/parameters/artifactory/password"),
					cli.File("/vela/secrets/artifactory/password"),
				),
			},
			&cli.StringFlag{
				Name:  "config.url",
				Usage: "Artifactory instance to communicate with",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_URL"),
					cli.EnvVar("ARTIFACTORY_URL"),
					cli.File("/vela/parameters/artifactory/url"),
					cli.File("/vela/secrets/artifactory/url"),
				),
			},
			&cli.StringFlag{
				Name:  "config.username",
				Usage: "user name for communication with the Artifactory instance",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_USERNAME"),
					cli.EnvVar("ARTIFACTORY_USERNAME"),
					cli.File("/vela/parameters/artifactory/username"),
					cli.File("/vela/secrets/artifactory/username"),
				),
			},

			// Client Flags

			&cli.IntFlag{
				Name:  "client.retries",
				Value: 3,
				Usage: "number of times to retry failed http attempts",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_HTTP_CLIENT_RETRIES"),
					cli.EnvVar("ARTIFACTORY_HTTP_CLIENT_RETRIES"),
					cli.File("/vela/parameters/artifactory/http_client_retries"),
					cli.File("/vela/secrets/artifactory/http_client_retries"),
				),
			},
			&cli.IntFlag{
				Name:  "client.retry_wait",
				Value: 500,
				Usage: "amount of milliseconds to wait between failed http attempts",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_HTTP_CLIENT_RETRY_WAIT_MILLISECONDS"),
					cli.EnvVar("ARTIFACTORY_HTTP_CLIENT_RETRY_WAIT_MILLISECONDS"),
					cli.File("/vela/parameters/artifactory/http_client_retry_wait_milliseconds"),
					cli.File("/vela/secrets/artifactory/http_client_retry_wait_milliseconds"),
				),
			},
			&cli.StringFlag{
				Name:  "client.cert",
				Usage: "file path to the client certificate to use for TLS communication",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_HTTP_CLIENT_CERT"),
					cli.EnvVar("ARTIFACTORY_HTTP_CLIENT_CERT"),
					cli.File("/vela/parameters/artifactory/http_client_cert"),
					cli.File("/vela/secrets/artifactory/http_client_cert"),
				),
			},
			&cli.StringFlag{
				Name:  "client.cert_key",
				Usage: "file path to the client certificate key to use for TLS communication",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_HTTP_CLIENT_CERT_KEY"),
					cli.EnvVar("ARTIFACTORY_HTTP_CLIENT_CERT_KEY"),
					cli.File("/vela/parameters/artifactory/http_client_cert_key"),
					cli.File("/vela/secrets/artifactory/http_client_cert_key"),
				),
			},
			&cli.BoolFlag{
				Name:  "client.insecure_tls",
				Usage: "enables skipping TLS verification when communicating with the Artifactory instance",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_HTTP_CLIENT_INSECURE_TLS"),
					cli.EnvVar("ARTIFACTORY_HTTP_CLIENT_INSECURE_TLS"),
					cli.File("/vela/parameters/artifactory/http_client_insecure_tls"),
					cli.File("/vela/secrets/artifactory/http_client_insecure_tls"),
				),
			},

			// Copy Flags

			&cli.BoolFlag{
				Name:  "copy.flat",
				Usage: "enables removing source directory hierarchy",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_FLAT"),
					cli.EnvVar("ARTIFACTORY_FLAT"),
					cli.File("/vela/parameters/artifactory/flat"),
					cli.File("/vela/secrets/artifactory/flat"),
				),
			},
			&cli.StringFlag{
				Name:  "copy.target",
				Usage: "target path to copy artifact(s) to",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_TARGET"),
					cli.EnvVar("ARTIFACTORY_TARGET"),
					cli.File("/vela/parameters/artifactory/target"),
					cli.File("/vela/secrets/artifactory/target"),
				),
			},

			// Docker Promote Flags

			&cli.StringFlag{
				Name:  "docker_promote.source_repo",
				Usage: "source Docker repository in Artifactory for the move or copy",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_SOURCE_REPO"),
					cli.EnvVar("ARTIFACTORY_SOURCE_REPO"),
					cli.File("/vela/parameters/artifactory/source_repo"),
					cli.File("/vela/secrets/artifactory/source_repo"),
				),
			},
			&cli.StringFlag{
				Name:  "docker_promote.target_repo",
				Usage: "destination Docker repository in Artifactory for the move or copy",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_TARGET_REPO"),
					cli.EnvVar("ARTIFACTORY_TARGET_REPO"),
					cli.File("/vela/parameters/artifactory/target_repo"),
					cli.File("/vela/secrets/artifactory/target_repo"),
				),
			},
			&cli.StringFlag{
				Name:  "docker_promote.docker_registry",
				Usage: "source Docker registry to promote an image from",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_DOCKER_REGISTRY"),
					cli.EnvVar("ARTIFACTORY_DOCKER_REGISTRY"),
					cli.File("/vela/parameters/artifactory/docker_registry"),
					cli.File("/vela/secrets/artifactory/docker_registry"),
				),
			},
			&cli.StringFlag{
				Name:  "docker_promote.target_docker_registry",
				Usage: "target Docker registry to promote an image to (uses 'docker_registry' if empty)",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_TARGET_DOCKER_REGISTRY"),
					cli.EnvVar("ARTIFACTORY_TARGET_DOCKER_REGISTRY"),
					cli.File("/vela/parameters/artifactory/target_docker_registry"),
					cli.File("/vela/secrets/artifactory/target_docker_registry"),
				),
			},
			&cli.StringFlag{
				Name:  "docker_promote.source_tag",
				Usage: "tag name of image to promote (promotes all tags if empty)",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_TAG"),
					cli.EnvVar("ARTIFACTORY_TAG"),
					cli.File("/vela/parameters/artifactory/tag"),
					cli.File("/vela/secrets/artifactory/tag"),
				),
			},
			&cli.StringSliceFlag{
				Name:  "docker_promote.target_tags",
				Usage: "target tag to assign the image after promotion",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_TARGET_TAGS"),
					cli.EnvVar("ARTIFACTORY_TARGET_TAGS"),
					cli.File("/vela/parameters/artifactory/target_tags"),
					cli.File("/vela/secrets/artifactory/target_tags"),
				),
			},
			&cli.BoolFlag{
				Name:  "docker_promote.copy",
				Value: true,
				Usage: "set to copy instead of moving the image",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_COPY"),
					cli.EnvVar("ARTIFACTORY_COPY"),
					cli.File("/vela/parameters/artifactory/copy"),
					cli.File("/vela/secrets/artifactory/copy"),
				),
			},
			&cli.BoolFlag{
				Name:  "docker_promote.props",
				Usage: "property to be set on the artifact when it is being promoted",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_PROMOTE_PROPS"),
					cli.EnvVar("ARTIFACTORY_PROMOTE_PROPS"),
					cli.File("/vela/parameters/artifactory/promote_props"),
					cli.File("/vela/secrets/artifactory/promote_props"),
				),
			},

			// Set Prop Flags

			&cli.StringFlag{
				Name:  "set_prop.props",
				Usage: "properties to set on the artifact(s)",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_PROPS"),
					cli.EnvVar("ARTIFACTORY_PROPS"),
					cli.File("/vela/parameters/artifactory/props"),
					cli.File("/vela/secrets/artifactory/props"),
				),
			},

			// Upload Flags

			&cli.BoolFlag{
				Name:  "upload.flat",
				Usage: "enables uploading artifacts to exact target path (excludes source file hierarchy)",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_FLAT"),
					cli.EnvVar("ARTIFACTORY_FLAT"),
					cli.File("/vela/parameters/artifactory/flat"),
					cli.File("/vela/secrets/artifactory/flat"),
				),
			},
			&cli.BoolFlag{
				Name:  "upload.include_dirs",
				Usage: "enables including directories from sources",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_INCLUDE_DIRS"),
					cli.EnvVar("ARTIFACTORY_INCLUDE_DIRS"),
					cli.File("/vela/parameters/artifactory/include_dirs"),
					cli.File("/vela/secrets/artifactory/include_dirs"),
				),
			},
			&cli.BoolFlag{
				Name:  "upload.regexp",
				Usage: "enables reading the sources as a regular expression",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_REGEXP"),
					cli.EnvVar("ARTIFACTORY_REGEXP"),
					cli.File("/vela/parameters/artifactory/regexp"),
					cli.File("/vela/secrets/artifactory/regexp"),
				),
			},
			&cli.StringSliceFlag{
				Name:  "upload.sources",
				Usage: "list of artifact(s) to upload",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_SOURCES"),
					cli.EnvVar("ARTIFACTORY_SOURCES"),
					cli.File("/vela/parameters/artifactory/sources"),
					cli.File("/vela/secrets/artifactory/sources"),
				),
			},
			&cli.StringFlag{
				Name:  "upload.build_props",
				Usage: "build props to apply",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("PARAMETER_BUILD_PROPS"),
					cli.EnvVar("ARTIFACTORY_BUILD_PROPS"),
					cli.File("/vela/parameters/artifactory/build_props"),
					cli.File("/vela/secrets/artifactory/build_props"),
				),
			},
		},
	}

	err = app.Run(context.Background(), os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

// run executes the plugin based off the configuration provided.
func run(_ context.Context, c *cli.Command) error {
	// set the log level for the plugin
	switch c.String("log.level") {
	case "t", "trace", "Trace", "TRACE":
		logrus.SetLevel(logrus.TraceLevel)
	case "d", "debug", "Debug", "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "w", "warn", "Warn", "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "e", "error", "Error", "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "f", "fatal", "Fatal", "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	case "p", "panic", "Panic", "PANIC":
		logrus.SetLevel(logrus.PanicLevel)
	case "i", "info", "Info", "INFO":
		fallthrough
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.WithFields(logrus.Fields{
		"code":     "https://github.com/go-vela/vela-artifactory",
		"docs":     "https://go-vela.github.io/docs/plugins/registry/pipeline/artifactory",
		"registry": "https://hub.docker.com/r/target/vela-artifactory",
	}).Info("Vela Artifactory Plugin")

	// When a user wants to configure the plugin using the
	// /vela/parameters/artifactory path there is a high probability
	// a user will do something like `echo 'VALUE' > /vela/parameters/artifactory/param`
	// which unintentionally introduced a new line into the parameter.
	// If we don't sanitize those paths the client library doesn't raise an error,
	// and instead behaves unexpectedly. In some cases with multiple file uploads
	// it will create a file with the last part of the path segment which takes on the value
	// of the last uploaded file.
	sanitizedPath := strings.TrimSpace(c.String("path"))
	sanitizedCopyTarget := strings.TrimSpace(c.String("copy.target"))

	// create the plugin
	p := &Plugin{
		// config configuration
		Config: &Config{
			Action:   c.String("config.action"),
			Token:    c.String("config.token"),
			APIKey:   c.String("config.api_key"),
			DryRun:   c.Bool("config.dry_run"),
			Password: c.String("config.password"),
			URL:      c.String("config.url"),
			Username: c.String("config.username"),
			// http client configuration
			Client: &Client{
				Retries:            c.Int("client.retries"),
				RetryWaitMilliSecs: c.Int("client.retry_wait"),
				CertPath:           c.String("client.cert"),
				CertKeyPath:        c.String("client.cert_key"),
				InsecureTLS:        c.Bool("client.insecure_tls"),
			},
		},
		// copy configuration
		Copy: &Copy{
			Flat:      c.Bool("copy.flat"),
			Path:      sanitizedPath,
			Recursive: c.Bool("recursive"),
			Target:    sanitizedCopyTarget,
		},
		// delete configuration
		Delete: &Delete{
			Path:      sanitizedPath,
			Recursive: c.Bool("recursive"),
		},
		// docker-promote configuration
		DockerPromote: &DockerPromote{
			SourceRepo:           c.String("docker_promote.source_repo"),
			TargetRepo:           c.String("docker_promote.target_repo"),
			DockerRegistry:       c.String("docker_promote.docker_registry"),
			TargetDockerRegistry: c.String("docker_promote.target_docker_registry"),
			SourceTag:            c.String("docker_promote.source_tag"),
			TargetTags:           c.StringSlice("docker_promote.target_tags"),
			Copy:                 c.Bool("docker_promote.copy"),
			PromoteProperty:      c.Bool("docker_promote.props"),
		},
		// set-prop configuration
		SetProp: &SetProp{
			Path:     sanitizedPath,
			RawProps: c.String("set_prop.props"),
		},
		// upload configuration
		Upload: &Upload{
			Flat:        c.Bool("upload.flat"),
			IncludeDirs: c.Bool("upload.include_dirs"),
			Recursive:   c.Bool("recursive"),
			Regexp:      c.Bool("upload.regexp"),
			Path:        sanitizedPath,
			Sources:     c.StringSlice("upload.sources"),
			BuildProps:  c.String("upload.build_props"),
		},
	}

	// validate the plugin
	err := p.Validate()
	if err != nil {
		return err
	}

	// execute the plugin
	return p.Exec()
}

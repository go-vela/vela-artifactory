// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

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
	app := cli.NewApp()

	// Plugin Information

	app.Name = "vela-artifactory"
	app.HelpName = "vela-artifactory"
	app.Usage = "Vela Artifactory plugin for managing artifacts"
	app.Copyright = "Copyright 2019 Target Brands, Inc. All rights reserved."
	app.Authors = []*cli.Author{
		{
			Name:  "Vela Admins",
			Email: "vela@target.com",
		},
	}

	// Plugin Metadata

	app.Action = run
	app.Compiled = time.Now()
	app.Version = v.Semantic()

	// Plugin Flags

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_LOG_LEVEL", "ARTIFACTORY_LOG_LEVEL"},
			FilePath: "/vela/parameters/artifactory/log_level,/vela/secrets/artifactory/log_level",
			Name:     "log.level",
			Usage:    "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:    "info",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_PATH", "ARTIFACTORY_PATH"},
			FilePath: "/vela/parameters/artifactory/path,/vela/secrets/artifactory/path",
			Name:     "path",
			Usage:    "source/target path to artifact(s) for action",
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_RECURSIVE", "ARTIFACTORY_RECURSIVE"},
			FilePath: "/vela/parameters/artifactory/recursive,/vela/secrets/artifactory/recursive",
			Name:     "recursive",
			Usage:    "enables operating on sub-directories for the source/target path",
		},

		// Config Flags

		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_ACTION", "ARTIFACTORY_ACTION"},
			FilePath: "/vela/parameters/artifactory/action,/vela/secrets/artifactory/action",
			Name:     "config.action",
			Usage:    "action to perform against the Artifactory instance",
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_DRY_RUN", "ARTIFACTORY_DRY_RUN"},
			FilePath: "/vela/parameters/artifactory/dry_run,/vela/secrets/artifactory/dry_run",
			Name:     "config.dry_run",
			Usage:    "enables pretending to perform the action",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_API_KEY", "ARTIFACTORY_API_KEY"},
			FilePath: "/vela/parameters/artifactory/api_key,/vela/secrets/artifactory/api_key,/vela/secrets/managed-auth/api_key",
			Name:     "config.api_key",
			Usage:    "API key for communication with the Artifactory instance",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_TOKEN", "ARTIFACTORY_TOKEN"},
			FilePath: "/vela/parameters/artifactory/token,/vela/secrets/artifactory/token,/vela/secrets/managed-auth/token",
			Name:     "config.token",
			Usage:    "Access/Identity token for communication with the Artifactory instance",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_PASSWORD", "ARTIFACTORY_PASSWORD"},
			FilePath: "/vela/parameters/artifactory/password,/vela/secrets/artifactory/password,/vela/secrets/managed-auth/password",
			Name:     "config.password",
			Usage:    "password for communication with the Artifactory instance",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_URL", "ARTIFACTORY_URL"},
			FilePath: "/vela/parameters/artifactory/url,/vela/secrets/artifactory/url",
			Name:     "config.url",
			Usage:    "Artifactory instance to communicate with",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_USERNAME", "ARTIFACTORY_USERNAME"},
			FilePath: "/vela/parameters/artifactory/username,/vela/secrets/artifactory/username,/vela/secrets/managed-auth/username",
			Name:     "config.username",
			Usage:    "user name for communication with the Artifactory instance",
		},

		// Client Flags

		&cli.IntFlag{
			EnvVars:  []string{"PARAMETER_HTTP_CLIENT_RETRIES", "ARTIFACTORY_HTTP_CLIENT_RETRIES"},
			FilePath: "/vela/parameters/artifactory/client/retries,/vela/secrets/artifactory/client/retries",
			Name:     "client.retries",
			Usage:    "number of times to retry failed http attempts",
			Value:    3,
		},
		&cli.IntFlag{
			EnvVars:  []string{"PARAMETER_HTTP_CLIENT_RETRY_WAIT_MILLISECONDS", "ARTIFACTORY_HTTP_CLIENT_RETRY_WAIT_MILLISECONDS"},
			FilePath: "/vela/parameters/artifactory/client/retry_wait,/vela/secrets/artifactory/client/retry_wait",
			Name:     "client.retry_wait",
			Usage:    "amount of milliseconds to wait between failed http attempts",
			Value:    500,
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_HTTP_CLIENT_CERT", "ARTIFACTORY_HTTP_CLIENT_CERT"},
			FilePath: "/vela/parameters/artifactory/client/cert,/vela/secrets/artifactory/client/cert",
			Name:     "client.cert",
			Usage:    "file path to the client certificate to use for TLS communication",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_HTTP_CLIENT_CERT_KEY", "ARTIFACTORY_HTTP_CLIENT_CERT_KEY"},
			FilePath: "/vela/parameters/artifactory/client/cert_key,/vela/secrets/artifactory/client/cert_key",
			Name:     "client.cert_key",
			Usage:    "file path to the client certificate key to use for TLS communication",
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_HTTP_CLIENT_INSECURE_TLS", "ARTIFACTORY_HTTP_CLIENT_INSECURE_TLS"},
			FilePath: "/vela/parameters/artifactory/client/insecure_tls,/vela/secrets/artifactory/client/insecure_tls",
			Name:     "client.insecure_tls",
			Usage:    "enables skipping TLS verification when communicating with the Artifactory instance",
		},

		// Copy Flags

		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_FLAT", "ARTIFACTORY_FLAT"},
			FilePath: "/vela/parameters/artifactory/flat,/vela/secrets/artifactory/flat",
			Name:     "copy.flat",
			Usage:    "enables removing source directory hierarchy",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_TARGET", "ARTIFACTORY_TARGET"},
			FilePath: "/vela/parameters/artifactory/target,/vela/secrets/artifactory/target",
			Name:     "copy.target",
			Usage:    "target path to copy artifact(s) to",
		},

		// Docker Promote Flags

		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_SOURCE_REPO", "ARTIFACTORY_SOURCE_REPO"},
			FilePath: "/vela/parameters/artifactory/source_repo,/vela/secrets/artifactory/source_repo",
			Name:     "docker_promote.source_repo",
			Usage:    "source Docker repository in Artifactory for the move or copy",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_TARGET_REPO", "ARTIFACTORY_TARGET_REPO"},
			FilePath: "/vela/parameters/artifactory/target_repo,/vela/secrets/artifactory/target_repo",
			Name:     "docker_promote.target_repo",
			Usage:    "destination Docker repository in Artifactory for the move or copy",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_DOCKER_REGISTRY", "ARTIFACTORY_DOCKER_REGISTRY"},
			FilePath: "/vela/parameters/artifactory/docker_registry,/vela/secrets/artifactory/docker_registry",
			Name:     "docker_promote.docker_registry",
			Usage:    "source Docker registry to promote an image from",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_TARGET_DOCKER_REGISTRY", "ARTIFACTORY_TARGET_DOCKER_REGISTRY"},
			FilePath: "/vela/parameters/artifactory/target_docker_registry,/vela/secrets/artifactory/target_docker_registry",
			Name:     "docker_promote.target_docker_registry",
			Usage:    "target Docker registry to promote an image to (uses 'docker_registry' if empty)",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_TAG", "ARTIFACTORY_TAG"},
			FilePath: "/vela/parameters/artifactory/tag,/vela/secrets/artifactory/tag",
			Name:     "docker_promote.source_tag",
			Usage:    "tag name of image to promote (promotes all tags if empty)",
		},
		&cli.StringSliceFlag{
			EnvVars:  []string{"PARAMETER_TARGET_TAGS", "ARTIFACTORY_TARGET_TAGS"},
			FilePath: "/vela/parameters/artifactory/target_tags,/vela/secrets/artifactory/target_tags",
			Name:     "docker_promote.target_tags",
			Usage:    "target tag to assign the image after promotion",
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_COPY", "ARTIFACTORY_COPY"},
			FilePath: "/vela/parameters/artifactory/copy,/vela/secrets/artifactory/copy",
			Name:     "docker_promote.copy",
			Usage:    "set to copy instead of moving the image",
			Value:    true,
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_PROMOTE_PROPS", "ARTIFACTORY_PROMOTE_PROPS"},
			FilePath: "/vela/parameters/artifactory/promote,/vela/secrets/artifactory/promote",
			Name:     "docker_promote.props",
			Usage:    "property to be set on the artifact when it is being promoted",
		},

		// Set Prop Flags

		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_PROPS", "ARTIFACTORY_PROPS"},
			FilePath: "/vela/parameters/artifactory/props,/vela/secrets/artifactory/props",
			Name:     "set_prop.props",
			Usage:    "properties to set on the artifact(s)",
		},

		// Upload Flags

		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_FLAT", "ARTIFACTORY_FLAT"},
			FilePath: "/vela/parameters/artifactory/flat,/vela/secrets/artifactory/flat",
			Name:     "upload.flat",
			Usage:    "enables uploading artifacts to exact target path (excludes source file hierarchy)",
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_INCLUDE_DIRS", "ARTIFACTORY_INCLUDE_DIRS"},
			FilePath: "/vela/parameters/artifactory/include_dirs,/vela/secrets/artifactory/include_dirs",
			Name:     "upload.include_dirs",
			Usage:    "enables including directories from sources",
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_REGEXP", "ARTIFACTORY_REGEXP"},
			FilePath: "/vela/parameters/artifactory/regexp,/vela/secrets/artifactory/regexp",
			Name:     "upload.regexp",
			Usage:    "enables reading the sources as a regular expression",
		},
		&cli.StringSliceFlag{
			EnvVars:  []string{"PARAMETER_SOURCES", "ARTIFACTORY_SOURCES"},
			FilePath: "/vela/parameters/artifactory/sources,/vela/secrets/artifactory/sources",
			Name:     "upload.sources",
			Usage:    "list of artifact(s) to upload",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_BUILD_PROPS", "ARTIFACTORY_BUILD_PROPS"},
			FilePath: "/vela/parameters/artifactory/build_props,/vela/secrets/artifactory/build_props",
			Name:     "upload.build_props",
			Usage:    "build props to apply",
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

// run executes the plugin based off the configuration provided.
func run(c *cli.Context) error {
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

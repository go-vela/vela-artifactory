// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := cli.NewApp()

	// Plugin Information

	app.Name = "vela-artifactory"
	app.HelpName = "vela-artifactory"
	app.Usage = "Vela Artifactory plugin for managing artifacts"
	app.Copyright = "Copyright (c) 2020 Target Brands, Inc. All rights reserved."
	app.Authors = []*cli.Author{
		{
			Name:  "Vela Admins",
			Email: "vela@target.com",
		},
	}

	// Plugin Metadata

	app.Compiled = time.Now()
	app.Action = run

	// Plugin Flags

	app.Flags = []cli.Flag{

		&cli.StringFlag{
			EnvVars: []string{"PARAMETER_LOG_LEVEL", "VELA_LOG_LEVEL", "ARTIFACTORY_LOG_LEVEL"},
			Name:    "log.level",
			Usage:   "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:   "info",
		},
		&cli.StringFlag{
			EnvVars: []string{"PARAMETER_PATH", "ARTIFACTORY_PATH"},
			Name:    "path",
			Usage:   "source/target path to artifact(s) for action",
		},

		// Config Flags

		&cli.StringFlag{
			EnvVars: []string{"PARAMETER_ACTION", "CONFIG_ACTION", "ARTIFACTORY_ACTION"},
			Name:    "config.action",
			Usage:   "action to perform against the Artifactory instance",
		},
		&cli.BoolFlag{
			EnvVars: []string{"PARAMETER_DRY_RUN", "CONFIG_DRY_RUN", "ARTIFACTORY_DRY_RUN"},
			Name:    "config.dry_run",
			Usage:   "enables pretending to perform the action",
		},
		&cli.StringFlag{
			EnvVars: []string{"PARAMETER_API_KEY", "CONFIG_API_KEY", "ARTIFACTORY_API_KEY"},
			Name:    "config.api_key",
			Usage:   "API key for communication with the Artifactory instance",
		},
		&cli.StringFlag{
			EnvVars: []string{"PARAMETER_PASSWORD", "CONFIG_PASSWORD", "ARTIFACTORY_PASSWORD"},
			Name:    "config.password",
			Usage:   "password for communication with the Artifactory instance",
		},
		&cli.StringFlag{
			EnvVars: []string{"PARAMETER_URL", "CONFIG_URL", "ARTIFACTORY_URL"},
			Name:    "config.url",
			Usage:   "Artifactory instance to communicate with",
		},
		&cli.StringFlag{
			EnvVars: []string{"PARAMETER_USERNAME", "CONFIG_USERNAME", "ARTIFACTORY_USERNAME"},
			Name:    "config.username",
			Usage:   "user name for communication with the Artifactory instance",
		},

		// Copy Flags

		&cli.BoolFlag{
			EnvVars: []string{"PARAMETER_FLAT", "COPY_FLAT"},
			Name:    "copy.flat",
			Usage:   "enables removing source directory hierarchy",
		},
		&cli.BoolFlag{
			EnvVars: []string{"PARAMETER_RECURSIVE", "COPY_RECURSIVE"},
			Name:    "copy.recursive",
			Usage:   "enables copying sub-directories for the artifact(s)",
		},
		&cli.StringFlag{
			EnvVars: []string{"PARAMETER_TARGET", "COPY_TARGET"},
			Name:    "copy.target",
			Usage:   "target path to copy artifact(s) to",
		},

		// Delete Flags

		&cli.BoolFlag{
			EnvVars: []string{"PARAMETER_RECURSIVE", "DELETE_RECURSIVE"},
			Name:    "delete.recursive",
			Usage:   "enables removing sub-directories for the artifact(s)",
		},

		// Docker promote Flags

		&cli.StringFlag{
			EnvVars: []string{"PARAMETER_TARGET_REPO", "DOCKER_PROMOTE_TARGET_REPO"},
			Name:    "docker_promote.target_repo",
			Usage:   "target repository is the repository for the move or copy",
		},
		&cli.StringFlag{
			EnvVars: []string{"PARAMETER_DOCKER_REGISTRY", "DOCKER_PROMOTE_DOCKER_REGISTRY"},
			Name:    "docker_promote.docker_registry",
			Usage:   "docker registry is the registry to promote",
		},
		&cli.StringFlag{
			EnvVars: []string{"PARAMETER_TARGET_DOCKER_REGISTRY", "DOCKER_PROMOTE_TARGET_DOCKER_REGISTRY"},
			Name:    "docker_promote.target_docker_registry",
			Usage:   "target docker registry is an optional target registry, if null, will use the same name as 'docker_registry'",
		},
		&cli.StringFlag{
			EnvVars: []string{"PARAMETER_TAG", "DOCKER_PROMOTE_TAG"},
			Name:    "docker_promote.tag",
			Usage:   "tag name to promote if null the entire docker repository will be promoted.",
		},
		&cli.StringSliceFlag{
			EnvVars: []string{"PARAMETER_TARGET_TAGS", "DOCKER_PROMOTE_TARGET_TAGS"},
			Name:    "docker_promote.target_tags",
			Usage:   "target tag to assign the image after promotion",
		},
		&cli.BoolFlag{
			EnvVars: []string{"PARAMETER_COPY", "DOCKER_PROMOTE_COPY"},
			Name:    "docker_promote.copy",
			Usage:   "set to copy instead of moving the image",
		},
		&cli.BoolFlag{
			EnvVars: []string{"PARAMETER_PROPS_PROMOTE", "SET_PROP_PROMOTE"},
			Name:    "docker_prop.promote",
			Usage:   "propety to be set on the artifact when it is being promoted",
		},

		// Set Prop Flags

		&cli.StringFlag{
			EnvVars: []string{"PARAMETER_PROPS", "SET_PROP_PROPS"},
			Name:    "set_prop.props",
			Usage:   "properties to set on the artifact(s)",
		},

		// Upload Flags

		&cli.BoolFlag{
			EnvVars: []string{"PARAMETER_FLAT", "UPLOAD_FLAT"},
			Name:    "upload.flat",
			Usage:   "enables uploading artifacts to exact target path (excludes source file hierarchy)",
		},
		&cli.BoolFlag{
			EnvVars: []string{"PARAMETER_INCLUDE_DIRS", "UPLOAD_INCLUDE_DIRS"},
			Name:    "upload.include_dirs",
			Usage:   "enables including directories from sources",
		},
		&cli.BoolFlag{
			EnvVars: []string{"PARAMETER_REGEXP", "UPLOAD_REGEXP"},
			Name:    "upload.regexp",
			Usage:   "enables reading the sources as a regular expression",
		},
		&cli.BoolFlag{
			EnvVars: []string{"PARAMETER_RECURSIVE", "UPLOAD_RECURSIVE"},
			Name:    "upload.recursive",
			Usage:   "enables uploading sub-directories for the sources",
		},
		&cli.StringSliceFlag{
			EnvVars: []string{"PARAMETER_SOURCES", "UPLOAD_SOURCES"},
			Name:    "upload.sources",
			Usage:   "list of artifact(s) to upload",
		},
	}

	err := app.Run(os.Args)
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
		"docs":     "https://go-vela.github.io/docs/plugins/registry/artifactory",
		"registry": "https://hub.docker.com/r/target/vela-artifactory",
	}).Info("Vela Artifactory Plugin")

	// create the plugin
	p := &Plugin{
		// config configuration
		Config: &Config{
			Action:   c.String("config.action"),
			APIKey:   c.String("config.api_key"),
			DryRun:   c.Bool("config.dry_run"),
			Password: c.String("config.password"),
			URL:      c.String("config.url"),
			Username: c.String("config.username"),
		},
		// copy configuration
		Copy: &Copy{
			Flat:      c.Bool("copy.flat"),
			Path:      c.String("path"),
			Recursive: c.Bool("copy.recursive"),
			Target:    c.String("copy.target"),
		},
		// delete configuration
		Delete: &Delete{
			Path:      c.String("path"),
			Recursive: c.Bool("delete.recursive"),
		},
		// docker-promote configuration
		DockerPromote: &DockerPromote{
			TargetRepo:           c.String("docker_promote.target_repo"),
			DockerRegistry:       c.String("docker_promote.docker_registry"),
			TargetDockerRegistry: c.String("docker_promote.target_docker_registry"),
			Tag:                  c.String("docker_promote.tag"),
			TargetTags:           c.StringSlice("docker_promote.target_tags"),
			Copy:                 c.Bool("docker_promote.copy"),
			PromoteProperty:      c.Bool("docker_prop.promote"),
		},
		// set-prop configuration
		SetProp: &SetProp{
			Path:     c.String("path"),
			RawProps: c.String("set_prop.props"),
		},
		// upload configuration
		Upload: &Upload{
			Flat:        c.Bool("upload.flat"),
			IncludeDirs: c.Bool("upload.include_dirs"),
			Recursive:   c.Bool("upload.recursive"),
			Regexp:      c.Bool("upload.regexp"),
			Path:        c.String("path"),
			Sources:     c.StringSlice("upload.sources"),
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

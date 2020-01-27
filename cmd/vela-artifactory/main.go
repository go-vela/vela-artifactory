// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := cli.NewApp()

	// Plugin Information

	app.Name = "vela-artifactory"
	app.HelpName = "vela-artifactory"
	app.Usage = "Vela Artifactory plugin for managing artifacts"
	app.Copyright = "Copyright (c) 2020 Target Brands, Inc. All rights reserved."
	app.Authors = []cli.Author{
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

		cli.StringFlag{
			EnvVar: "PARAMETER_LOG_LEVEL,VELA_LOG_LEVEL,ARTIFACTORY_LOG_LEVEL",
			Name:   "log.level",
			Usage:  "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:  "info",
		},
		cli.StringFlag{
			EnvVar: "PARAMETER_PATH,ARTIFACTORY_PATH",
			Name:   "path",
			Usage:  "source/target path to artifact(s) for action",
		},

		// Config Flags

		cli.StringFlag{
			EnvVar: "PARAMETER_ACTION,CONFIG_ACTION,ARTIFACTORY_ACTION",
			Name:   "config.action",
			Usage:  "action to perform against the Artifactory instance",
		},
		cli.BoolFlag{
			EnvVar: "PARAMETER_DRY_RUN,CONFIG_DRY_RUN,ARTIFACTORY_DRY_RUN",
			Name:   "config.dry_run",
			Usage:  "enables pretending to perform the action",
		},
		cli.StringFlag{
			EnvVar: "PARAMETER_API_KEY,CONFIG_API_KEY,ARTIFACTORY_API_KEY",
			Name:   "config.api_key",
			Usage:  "API key for communication with the Artifactory instance",
		},
		cli.StringFlag{
			EnvVar: "PARAMETER_PASSWORD,CONFIG_PASSWORD,ARTIFACTORY_PASSWORD",
			Name:   "config.password",
			Usage:  "password for communication with the Artifactory instance",
		},
		cli.StringFlag{
			EnvVar: "PARAMETER_URL,CONFIG_URL,ARTIFACTORY_URL",
			Name:   "config.url",
			Usage:  "Artifactory instance to communicate with",
		},
		cli.StringFlag{
			EnvVar: "PARAMETER_USERNAME,CONFIG_USERNAME,ARTIFACTORY_USERNAME",
			Name:   "config.username",
			Usage:  "user name for communication with the Artifactory instance",
		},

		// Copy Flags

		cli.BoolFlag{
			EnvVar: "PARAMETER_FLAT,COPY_FLAT",
			Name:   "copy.flat",
			Usage:  "enables removing source directory hierarchy",
		},
		cli.BoolFlag{
			EnvVar: "PARAMETER_RECURSIVE,COPY_RECURSIVE",
			Name:   "copy.recursive",
			Usage:  "enables copying sub-directories for the artifact(s)",
		},
		cli.StringFlag{
			EnvVar: "PARAMETER_TARGET,COPY_TARGET",
			Name:   "copy.target",
			Usage:  "target path to copy artifact(s) to",
		},

		// Delete Flags

		cli.BoolFlag{
			EnvVar: "PARAMETER_RECURSIVE,DELETE_RECURSIVE",
			Name:   "delete.recursive",
			Usage:  "enables removing sub-directories for the artifact(s)",
		},

		// Set Prop Flags

		cli.StringFlag{
			EnvVar: "PARAMETER_PROPS,SET_PROP_PROPS",
			Name:   "set_prop.props",
			Usage:  "properties to set on the artifact(s)",
		},

		// Upload Flags

		cli.BoolFlag{
			EnvVar: "PARAMETER_FLAT,UPLOAD_FLAT",
			Name:   "upload.flat",
			Usage:  "enables uploading artifacts to exact target path (excludes source file hierarchy)",
		},
		cli.BoolFlag{
			EnvVar: "PARAMETER_INCLUDE_DIRS,UPLOAD_INCLUDE_DIRS",
			Name:   "upload.include_dirs",
			Usage:  "enables including directories from sources",
		},
		cli.BoolFlag{
			EnvVar: "PARAMETER_REGEXP,UPLOAD_REGEXP",
			Name:   "upload.regexp",
			Usage:  "enables reading the sources as a regular expression",
		},
		cli.BoolFlag{
			EnvVar: "PARAMETER_RECURSIVE,UPLOAD_RECURSIVE",
			Name:   "upload.recursive",
			Usage:  "enables uploading sub-directories for the sources",
		},
		cli.StringSliceFlag{
			EnvVar: "PARAMETER_SOURCES,UPLOAD_SOURCES",
			Name:   "upload.sources",
			Usage:  "list of artifact(s) to upload",
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
		"code": "https://github.com/go-vela/vela-artifactory",
		"docs": "https://go-vela.github.io/docs/plugins/registry/artifactory",
		"time": time.Now(),
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

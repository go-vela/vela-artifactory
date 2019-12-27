// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"
	"os"

	"github.com/go-vela/vela-artifactory/artifactory"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var revision string // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "artifactory plugin"
	app.Usage = "artifactory plugin"
	app.Action = run
	app.Version = revision
	app.Flags = []cli.Flag{

		//
		// plugin args
		//

		cli.StringFlag{
			Name:   "actions",
			Usage:  "Actions to perform against artifactory",
			EnvVar: "PARAMETER_ACTIONS",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "Enable debug logging",
			EnvVar: "PARAMETER_DEBUG",
		},
		cli.StringFlag{
			Name:   "password",
			Usage:  "Artifactory server password",
			EnvVar: "PARAMETER_PASSWORD,ARTIFACTORY_PASSWORD",
		},
		cli.StringFlag{
			Name:   "apikey",
			Usage:  "Artifactory API Key",
			EnvVar: "PARAMETER_APIKEY,ARTIFACTORY_APIKEY",
		},
		cli.StringFlag{
			Name:   "url",
			Usage:  "Artifactory server URL",
			EnvVar: "PARAMETER_URL",
		},
		cli.StringFlag{
			Name:   "username",
			Usage:  "Artifactory server username",
			EnvVar: "PARAMETER_USERNAME,ARTIFACTORY_USERNAME",
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	logrus.WithFields(logrus.Fields{
		"Revision": revision,
	}).Info("Artifactory Plugin Version")

	if c.Bool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	actions, err := unmarshalActions(c.String("actions"))

	if err != nil {
		return err
	}

	plugin := Plugin{
		Actions: actions,
		Config: artifactory.Config{
			APIKey:   c.String("apikey"),
			Debug:    c.Bool("debug"),
			Password: c.String("password"),
			URL:      c.String("url"),
			Username: c.String("username"),
		},
	}

	err = plugin.Exec()

	return err
}

func unmarshalActions(rawJSON string) ([]Action, error) {
	logrus.WithFields(logrus.Fields{
		"raw-actions": rawJSON,
	}).Debug()

	bytes := []byte(rawJSON)
	var actions []Action
	err := json.Unmarshal(bytes, &actions)

	if err != nil {
		return nil, err
	}

	var rawActions []json.RawMessage
	err = json.Unmarshal(bytes, &rawActions)

	if err != nil {
		return nil, err
	}

	for idx := range actions {
		actions[idx].RawArguments = rawActions[idx]
	}

	logrus.WithFields(logrus.Fields{
		"actions": actions,
	}).Debug()

	return actions, nil
}

// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	goarty "github.com/target/go-arty/artifactory"
)

const dockerPromoteAction = "docker-promote"

// DockerPromote represents the plugin configuration for setting a Docker Promotion.
type DockerPromote struct {
	// target repo is the target repository for the move or copy
	TargetRepo string
	// docker registry is the name of the registry promote
	DockerRegistry string
	// target docker registry is an optional target registry, if null, will use the same name as 'DockerRegistry'
	TargetDockerRegistry string
	// tag is an optional tag name to promote, if null - the entire docker repository will be promoted. Available from v4.10.
	Tag string
	// target tag is an optional target tag to assign the image after promotion, if null - will use the same tag
	TargetTags []string
	// An optional value to set whether to copy instead of move. Default: true
	Copy bool
	// An optional value to set an item property to add a promoted date.
	PromoteProperty bool
}

// Exec formats and runs the commands for uploading artifacts in Artifactory.
func (p *DockerPromote) Exec(c *Config) error {
	logrus.Trace("running docker-promote with provided configuration")

	// create go-arty client for interacting with Docker promotion
	// https://github.com/target/go-arty
	client, err := goarty.NewClient(c.URL, nil)
	if err != nil {
		return err
	}

	// set basic authentication on client
	client.Authentication.SetBasicAuth(c.Username, c.Password)

	// set the auth token if user passed token
	if len(c.APIKey) != 0 {
		client.Authentication.SetTokenAuth(c.APIKey)
	}

	var payloads []*goarty.ImagePromotion

	for _, t := range p.TargetTags {
		payload := &goarty.ImagePromotion{
			TargetRepo:             goarty.String(p.TargetRepo),
			DockerRepository:       goarty.String(p.DockerRegistry),
			TargetDockerRepository: goarty.String(p.TargetDockerRegistry),
			Tag:                    goarty.String(p.Tag),
			TargetTag:              goarty.String(t),
			Copy:                   goarty.Bool(p.Copy),
		}

		payloads = append(payloads, payload)

		pretty, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return err
		}

		logrus.Tracef("Created payload for target tag %s: %s", t, string(pretty))
	}

	for _, payload := range payloads {
		logrus.Debugf("Promoting target tag %s", payload.GetTargetTag())

		_, _, err := client.Docker.PromoteImage(p.TargetRepo, payload)
		if err != nil {
			return err
		}

		if p.PromoteProperty {
			promotedImagePath := fmt.Sprintf("%s/%s", payload.GetTargetRepo(), payload.GetTargetTag())

			properties := make(map[string][]string)
			properties["promoted_on"] = append(properties["promoted_on"], time.Now().UTC().Format(time.RFC3339))

			_, err = client.Storage.SetItemProperties(payload.GetDockerRepository(), promotedImagePath, properties)
			if err != nil {
				return err
			}
		}

		logrus.Infof("Promotion ended successfully for target tag %s", payload.GetTargetTag())
	}

	return nil
}

// Validate verifies the Promote is properly configured.
func (p *DockerPromote) Validate() error {
	// verify a target repo is provided
	if len(p.TargetRepo) == 0 {
		return fmt.Errorf("no target repository provided")
	}

	// verify a docker registry is provided
	if len(p.DockerRegistry) == 0 {
		return fmt.Errorf("no docker repository provided")
	}

	return nil
}

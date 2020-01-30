// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"
	"fmt"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/sirupsen/logrus"
	goarty "github.com/target/go-arty/artifactory"
)

const dockerPromoteAction = "docker-promote"

// DockerPromote represents the plugin configuration for setting a Docker Promotion.
type DockerPromote struct {
	// target repo is the target repository for the move or copy
	TargetRepo string
	// docker repo it the docker repository name to promote
	DockerRepo string
	// target docker repo An optional docker repository name, if null, will use the same name as 'dockerRepository'
	TargetDockerRepo string
	// tag is an optional tag name to promote, if null - the entire docker repository will be promoted. Available from v4.10.
	Tag string
	// target tag is an optional target tag to assign the image after promotion, if null - will use the same tag
	TargetTags []string
	// An optional value to set whether to copy instead of move. Default: true
	Copy bool
}

// Exec formats and runs the commands for uploading artifacts in Artifactory.
func (p *DockerPromote) Exec(cli *artifactory.ArtifactoryServicesManager) error {
	logrus.Trace("running docker-promote with provided configuration")

	// get client information
	username := cli.GetConfig().GetArtDetails().GetUser()
	password := cli.GetConfig().GetArtDetails().GetPassword()
	url := cli.GetConfig().GetArtDetails().GetUrl()

	// create go-arty client for interacting with Docker promotion
	// https://github.com/target/go-arty
	client, _ := goarty.NewClient(url, nil)
	client.Authentication.SetBasicAuth(username, password)

	var payloads []*goarty.ImagePromotion

	for _, t := range p.TargetTags {
		payload := &goarty.ImagePromotion{
			TargetRepo:             goarty.String(p.TargetRepo),
			DockerRepository:       goarty.String(p.DockerRepo),
			TargetDockerRepository: goarty.String(p.TargetDockerRepo),
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
		logrus.Debugf("Promoting target tag %s", *payload.TargetTag)

		_, _, err := client.Docker.PromoteImage(p.TargetRepo, payload)
		if err != nil {
			return err
		}

		logrus.Infof("Promotion ended successfully for target tag %s", *payload.TargetTag)
	}

	return nil
}

// Validate verifies the Promote is properly configured.
func (p *DockerPromote) Validate() error {
	// verify a target tags are provided
	if len(p.TargetRepo) == 0 {
		return fmt.Errorf("no target repository provided")
	}

	// verify a target tags are provided
	if len(p.DockerRepo) == 0 {
		return fmt.Errorf("no docker repository provided")
	}

	return nil
}

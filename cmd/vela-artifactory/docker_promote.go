// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/sirupsen/logrus"
)

const dockerPromoteAction = "docker-promote"

// DockerPromote represents the plugin configuration for setting a Docker Promotion.
type DockerPromote struct {
	// Docker repository in Artifactory for the move or copy
	TargetRepo string
	// DockerRegistry is the source Docker registry to promote an image from
	DockerRegistry string
	// TargetDockerRegistry is the target Docker registry to promote an image to (uses 'DockerRegistry' if empty)
	TargetDockerRegistry string
	// Tag is the name of image to promote (promotes all tags if empty)
	Tag string
	// TargetTags are the target tags to assign to the image after promotion
	TargetTags []string
	// Copy is a flag to set to copy instead of moving the image (default: true)
	Copy bool
	// PromoteProperty is an optional value to set an item property to add a promoted date.
	PromoteProperty bool
}

// Exec formats and runs the commands for uploading artifacts in Artifactory.
func (p *DockerPromote) Exec(cli artifactory.ArtifactoryServicesManager) error {
	logrus.Trace("running docker-promote with provided configuration")

	var payloads []*services.DockerPromoteParams

	for _, t := range p.TargetTags {
		params := services.NewDockerPromoteParams(
			p.TargetRepo,
			p.DockerRegistry,
			p.TargetDockerRegistry,
		)

		params.SourceTag = p.Tag
		params.TargetTag = t
		params.Copy = p.Copy

		payloads = append(payloads, &params)

		pretty, err := json.MarshalIndent(params, "", "  ")
		if err != nil {
			return err
		}

		logrus.Tracef("Created payload for target tag %s: %s", t, string(pretty))
	}

	for _, payload := range payloads {
		logrus.Debugf("Promoting target tag %s", payload.GetTargetTag())

		err := cli.PromoteDocker(*payload)
		if err != nil {
			return err
		}

		if p.PromoteProperty {
			var promotedImagePath string

			if payload.TargetRepo != "" {
				promotedImagePath = fmt.Sprintf(
					"%s/%s",
					payload.TargetRepo,
					payload.TargetTag,
				)
			} else {
				promotedImagePath = fmt.Sprintf(
					"%s/%s",
					payload.SourceRepo,
					payload.TargetTag,
				)
			}

			ts := time.Now().UTC().Format(time.RFC3339)
			properties := fmt.Sprintf("promoted_on=%s", ts)

			searchParams := services.NewSearchParams()
			searchParams.Recursive = true
			searchParams.IncludeDirs = true
			searchParams.Pattern = promotedImagePath

			reader, err := cli.SearchFiles(searchParams)
			if err != nil {
				return err
			}

			defer reader.Close()

			propsParams := services.NewPropsParams()
			propsParams.Props = properties
			propsParams.Reader = reader

			_, err = cli.SetProps(propsParams)

			if err != nil {
				return err
			}
		}

		logrus.Infof("Promotion ended successfully for tag %s for target tag %s",
			payload.SourceTag,
			payload.TargetTag)
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

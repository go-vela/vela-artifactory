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

	if len(p.TargetTags) == 0 {
		logrus.Trace("no tags to promote")
	}

	for _, t := range p.TargetTags {
		params := services.NewDockerPromoteParams(
			"",
			"",
			"",
		)
		params.SourceRepo = p.TargetRepo
		params.TargetRepo = p.TargetRepo

		params.SourceDockerImage = p.DockerRegistry
		params.TargetDockerImage = p.TargetDockerRegistry

		params.SourceTag = p.Tag
		params.TargetTag = t

		params.Copy = p.Copy

		payloads = append(payloads, &params)

		pretty, err := json.MarshalIndent(params, "", "  ")
		if err != nil {
			return err
		}

		logrus.Tracef("created payload for target tag %s: %s", t, string(pretty))
	}

	for _, payload := range payloads {
		logrus.Infof("Promoting tag %s to target %s", payload.GetSourceTag(), payload.GetTargetTag())

		err := cli.PromoteDocker(*payload)
		if err != nil {
			return err
		}

		if p.PromoteProperty {
			promotedImagePath := fmt.Sprintf(
				"%s/%s/%s",
				payload.TargetRepo,
				payload.TargetDockerImage,
				payload.TargetTag,
			)

			logrus.Infof("Setting promote properties for %s", promotedImagePath)

			searchParams := services.NewSearchParams()
			searchParams.Recursive = true
			searchParams.IncludeDirs = true
			searchParams.Pattern = promotedImagePath

			logrus.Tracef("searching files using pattern %s", promotedImagePath)

			reader, err := cli.SearchFiles(searchParams)
			if err != nil {
				return err
			}

			defer reader.Close()

			propsParams := services.NewPropsParams()
			ts := time.Now().UTC().Format(time.RFC3339)
			properties := fmt.Sprintf("promoted_on=%s", ts)
			propsParams.Props = properties
			propsParams.Reader = reader

			logrus.Tracef("assigning property [%s] to %d matched images", properties, len(reader.GetFilesPaths()))

			totalSuccess, err := cli.SetProps(propsParams)
			if err != nil {
				return err
			}

			if totalSuccess > 0 {
				logrus.Infof("Successfully assigned property [%s] to %d images", properties, totalSuccess)
			} else {
				logrus.Info("Promote properties not assigned to any images.")
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

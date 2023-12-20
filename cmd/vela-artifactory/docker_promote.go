// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/sirupsen/logrus"
)

const dockerPromoteAction = "docker-promote"

// DockerPromote represents the plugin configuration for setting a Docker Promotion.
type DockerPromote struct {
	// SourceRepo is the Docker repository in Artifactory to use as the source for move or copy
	SourceRepo string
	// TargetRepo is the Docker repository in Artifactory to use as the destination for move or copy
	TargetRepo string
	// DockerRegistry is the source Docker registry to promote an image from
	DockerRegistry string
	// TargetDockerRegistry is the target Docker registry to promote an image to (uses 'DockerRegistry' if empty)
	TargetDockerRegistry string
	// SourceTag is the name of image to promote (promotes all tags if empty)
	SourceTag string
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
		// avoid assigning parameters via constructor to ensure promote endpoint is constructed properly
		params := services.NewDockerPromoteParams(
			"",
			"",
			"",
		)

		params.SourceRepo = p.SourceRepo
		params.TargetRepo = p.TargetRepo

		// assign source repo to target repo when not provided, to ensure backwards compatibility
		// with existing configurations that assume source repo is not required
		if len(params.SourceRepo) == 0 {
			params.SourceRepo = params.TargetRepo
		}

		params.SourceDockerImage = p.DockerRegistry
		params.TargetDockerImage = p.TargetDockerRegistry

		// use DockerRegistry as TargetDockerRegistry when a target is not set
		if len(params.TargetDockerImage) == 0 {
			params.TargetDockerImage = params.SourceDockerImage
		}

		params.SourceTag = p.SourceTag
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

		retries := 3
		var retryErr error

		// send API call to promote Docker image in Artifactory
		// retry a couple times to avoid intermittent from Artifactory
		for i := 0; i < retries; i++ {
			backoff := i*2 + 1

			retryErr = cli.PromoteDocker(*payload)
			if retryErr != nil {
				logrus.Errorf("Error promoting tag %s to target %s on attempt %d: %s",
					payload.GetSourceTag(), payload.GetTargetTag(), i+1, retryErr.Error())

				if i == retries-1 {
					return retryErr
				}

				logrus.Infof("Retrying in %d seconds...", backoff)

				time.Sleep(time.Duration(backoff) * time.Second)
			} else {
				break
			}
		}

		// recursively assign promoted_on property based on plugin configuration
		if p.PromoteProperty {
			logrus.Infof("Setting promote properties for %s/%s/%s",
				payload.TargetRepo,
				payload.TargetDockerImage,
				payload.TargetTag)

			// setup base property params
			ts := time.Now().UTC().Format(time.RFC3339)
			promotedOnProperty := fmt.Sprintf("promoted_on=%s", ts)
			propsParams := services.NewPropsParams()
			propsParams.Props = promotedOnProperty

			// setup base search params
			searchParams := services.NewSearchParams()
			searchParams.Recursive = true
			searchParams.IncludeDirs = true

			// this is required to set properties on the image's folder artifact
			imageFolderSearchPattern := fmt.Sprintf(
				"%s/%s/%s",
				payload.TargetRepo,
				payload.TargetDockerImage,
				payload.TargetTag,
			)

			logrus.Tracef("searching files using pattern %s", imageFolderSearchPattern)

			searchParams.Pattern = imageFolderSearchPattern

			var imageFolderReader *content.ContentReader

			// send API call to search for files in Artifactory
			// retry a couple times to avoid intermittent from Artifactory
			for i := 0; i < retries; i++ {
				backoff := i*2 + 1

				imageFolderReader, retryErr = cli.SearchFiles(searchParams)
				if retryErr != nil {
					logrus.Errorf("Error searching for files using pattern %s on attempt %d: %s",
						imageFolderSearchPattern, i+1, retryErr.Error())

					if i == retries-1 {
						return retryErr
					}

					logrus.Infof("Retrying in %d seconds...", backoff)

					time.Sleep(time.Duration(backoff) * time.Second)
				} else {
					break
				}
			}

			// close the reader when done
			defer imageFolderReader.Close()

			// assign the files found to be used in SetProps
			propsParams.Reader = imageFolderReader

			logrus.Tracef("assigning property [%s] to %d matched images", promotedOnProperty, len(imageFolderReader.GetFilesPaths()))

			imageFolderSuccess := 0

			// send API call to set properties on a Docker image folder in Artifactory
			// retry a couple times to avoid intermittent from Artifactory
			for i := 0; i < retries; i++ {
				backoff := i*2 + 1

				imageFolderSuccess, retryErr = cli.SetProps(propsParams)
				if retryErr != nil {
					logrus.Errorf("Error setting properties for folder on tag %s to target %s on attempt %d: %s",
						payload.GetSourceTag(), payload.GetTargetTag(), i+1, retryErr.Error())

					if i == retries-1 {
						return retryErr
					}

					logrus.Infof("Retrying in %d seconds...", backoff)

					time.Sleep(time.Duration(backoff) * time.Second)
				} else {
					break
				}
			}

			// setup base search params
			searchParams = services.NewSearchParams()
			searchParams.Recursive = true
			searchParams.IncludeDirs = true

			// this is required to set properties on the image's folder artifact
			imageContentsSearchPattern := fmt.Sprintf(
				"%s/%s/%s/*",
				payload.TargetRepo,
				payload.TargetDockerImage,
				payload.TargetTag,
			)

			logrus.Tracef("searching files using pattern %s", imageContentsSearchPattern)

			searchParams.Pattern = imageContentsSearchPattern

			var imageContentsReader *content.ContentReader

			// send API call to search for files in Artifactory
			// use a retry backoff to avoid 403 errors
			for i := 0; i < retries; i++ {
				backoff := i*2 + 1

				imageContentsReader, retryErr = cli.SearchFiles(searchParams)
				if retryErr != nil {
					logrus.Errorf("Error searching for files using pattern %s on attempt %d: %s",
						imageContentsSearchPattern, i+1, retryErr.Error())

					if i == retries-1 {
						return retryErr
					}

					logrus.Infof("Retrying in %d seconds...", backoff)

					time.Sleep(time.Duration(backoff) * time.Second)
				} else {
					break
				}
			}

			defer imageContentsReader.Close()

			// assign the files found to be used in SetProps
			propsParams.Reader = imageContentsReader

			logrus.Tracef("assigning property [%s] to %d matched images", promotedOnProperty, len(imageContentsReader.GetFilesPaths()))

			imageContentsSuccess := 0

			// send API call to set properties on a Docker image contents in Artifactory
			// use a retry backoff to avoid 403 errors
			for i := 0; i < retries; i++ {
				backoff := i*2 + 1

				imageContentsSuccess, retryErr = cli.SetProps(propsParams)
				if retryErr != nil {
					logrus.Errorf("Error setting properties for contents on tag %s to target %s on attempt %d: %s",
						payload.GetSourceTag(), payload.GetTargetTag(), i+1, retryErr.Error())

					if i == retries-1 {
						return retryErr
					}

					logrus.Infof("Retrying in %d seconds...", backoff)

					time.Sleep(time.Duration(backoff) * time.Second)
				} else {
					break
				}
			}

			totalSuccess := imageFolderSuccess + imageContentsSuccess
			if totalSuccess > 0 {
				logrus.Infof("Successfully assigned property [%s] to %d image files.", promotedOnProperty, totalSuccess)
			} else {
				logrus.Warn("Promote properties not assigned to any image files.")
			}
		}

		logrus.Infof("Promotion ended successfully for tag %s promoted to target tag %s",
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

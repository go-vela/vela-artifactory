# Copyright (c) 2020 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

build: binary-build

run: build docker-build docker-run

test: build docker-build docker-example

#################################
######      Go clean       ######
#################################

clean:

	@go mod tidy
	@go vet ./...
	@go fmt ./...
	@echo "I'm kind of the only name in clean energy right now"

#################################
######    Build Binary     ######
#################################

binary-build:

	GOOS=linux CGO_ENABLED=0 go build -o release/vela-artifactory github.com/go-vela/vela-artifactory/cmd/vela-artifactory

#################################
######    Docker Build     ######
#################################

docker-build:

	docker build --no-cache -t vela-artifactory:local .

#################################
######     Docker Run      ######
#################################

docker-run:

	docker run --rm \
		-e ARTIFACTORY_API_KEY \
		-e ARTIFACTORY_PASSWORD \
		-e ARTIFACTORY_USERNAME \
		-e PARAMETER_ACTION \
		-e PARAMETER_DRY_RUN \
		-e PARAMETER_FLAT \
		-e PARAMETER_INCLUDE_DIRS \
		-e PARAMETER_PROPS \
		-e PARAMETER_RECURSIVE \
		-e PARAMETER_REGEXP \
		-e PARAMETER_SOURCES \
		-e PARAMETER_TARGET \
		-e PARAMETER_URL \
		vela-artifactory:local

docker-example:

	docker run --rm \
		-e PARAMETER_ACTION=upload \
		-e PARAMETER_DRY_RUN=true \
		-e PARAMETER_FLAT=false \
		-e PARAMETER_INCLUDE_DIRS=false \
		-e PARAMETER_PATH \
		-e PARAMETER_RECURSIVE=false \
		-e PARAMETER_REGEXP=false \
		-e PARAMETER_SOURCES=LICENSE \
		-e PARAMETER_URL \
		vela-artifactory:local

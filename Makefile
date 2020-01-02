# Copyright (c) 2019 Target Brands, Inc. All rights reserved.
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
		-e PARAMETER_ACTIONS \
		-e PARAMETER_DEBUG \
		-e PARAMETER_PASSWORD \
		-e PARAMETER_APIKEY \
		-e PARAMETER_URL \
		-e PARAMETER_USERNAME \
		-v $(pwd)/README.md:/README.md \
		vela-artifactory:local

docker-example:

	docker run --rm \
		-e PARAMETER_ACTIONS \
		-e PARAMETER_DEBUG \
		-e PARAMETER_PASSWORD \
		-e PARAMETER_APIKEY \
		-e PARAMETER_URL \
		-e PARAMETER_USERNAME \
		-v $(pwd)/README.md:/README.md \
		vela-artifactory:local

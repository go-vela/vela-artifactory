# Copyright (c) 2021 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

################################################################################
##    docker build --no-cache --target binary -t vela-artifactory:binary .    ##
################################################################################

FROM alpine as binary

ARG JFROG_VERSION=1.33.2

ADD https://releases.jfrog.io/artifactory/jfrog-cli/v1/${JFROG_VERSION}/jfrog-cli-linux-amd64/jfrog /bin/jfrog

RUN chmod a+x /bin/jfrog

##############################################################################
##    docker build --no-cache --target certs -t vela-artifactory:certs .    ##
##############################################################################

FROM alpine as certs

RUN apk add --update --no-cache ca-certificates

###############################################################
##    docker build --no-cache -t vela-artifactory:local .    ##
###############################################################

FROM scratch

COPY --from=binary /bin/jfrog /bin/jfrog

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/vela-artifactory /bin/vela-artifactory

ENTRYPOINT ["/bin/vela-artifactory"]

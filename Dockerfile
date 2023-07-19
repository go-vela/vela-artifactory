# Copyright (c) 2022 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

################################################################################
##    docker build --no-cache --target binary -t vela-artifactory:binary .    ##
################################################################################

FROM alpine@sha256:82d1e9d7ed48a7523bdebc18cf6290bdb97b82302a8a9c27d4fe885949ea94d1 as binary

ARG JFROG_VERSION=1.33.2

ADD https://releases.jfrog.io/artifactory/jfrog-cli/v1/${JFROG_VERSION}/jfrog-cli-linux-amd64/jfrog /bin/jfrog

RUN chmod a+x /bin/jfrog

##############################################################################
##    docker build --no-cache --target certs -t vela-artifactory:certs .    ##
##############################################################################

FROM alpine@sha256:82d1e9d7ed48a7523bdebc18cf6290bdb97b82302a8a9c27d4fe885949ea94d1 as certs

RUN apk add --update --no-cache ca-certificates

###############################################################
##    docker build --no-cache -t vela-artifactory:local .    ##
###############################################################

FROM scratch

COPY --from=binary /bin/jfrog /bin/jfrog

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/vela-artifactory /bin/vela-artifactory

ENTRYPOINT ["/bin/vela-artifactory"]

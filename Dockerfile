# SPDX-License-Identifier: Apache-2.0

################################################################################
##    docker build --no-cache --target binary -t vela-artifactory:binary .    ##
################################################################################

FROM alpine:3.24.4@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf as binary

ARG JFROG_VERSION=1.54.1

ADD https://releases.jfrog.io/artifactory/jfrog-cli/v1/${JFROG_VERSION}/jfrog-cli-linux-amd64/jfrog /bin/jfrog

RUN chmod a+x /bin/jfrog
RUN chmod -R 777 /tmp

##############################################################################
##    docker build --no-cache --target certs -t vela-artifactory:certs .    ##
##############################################################################

FROM alpine:3.24.4@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf as certs

RUN apk add --update --no-cache ca-certificates

###############################################################
##    docker build --no-cache -t vela-artifactory:local .    ##
###############################################################

FROM scratch

COPY --from=binary /bin/jfrog /bin/jfrog
COPY --from=binary /tmp /tmp

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/vela-artifactory /bin/vela-artifactory

ENTRYPOINT ["/bin/vela-artifactory"]

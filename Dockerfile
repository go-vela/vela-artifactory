# SPDX-License-Identifier: Apache-2.0

################################################################################
##    docker build --no-cache --target binary -t vela-artifactory:binary .    ##
################################################################################

FROM alpine:3.19.1@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b as binary

ARG JFROG_VERSION=1.33.2

ADD https://releases.jfrog.io/artifactory/jfrog-cli/v1/${JFROG_VERSION}/jfrog-cli-linux-amd64/jfrog /bin/jfrog

RUN chmod a+x /bin/jfrog
RUN chmod -R 777 /tmp

##############################################################################
##    docker build --no-cache --target certs -t vela-artifactory:certs .    ##
##############################################################################

FROM alpine:3.19.1@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b as certs

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

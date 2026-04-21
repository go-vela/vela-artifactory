# SPDX-License-Identifier: Apache-2.0

################################################################################
##    docker build --no-cache --target binary -t vela-artifactory:binary .    ##
################################################################################

FROM alpine:3.23.4@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11 as binary

ARG JFROG_VERSION=1.54.1

ADD https://releases.jfrog.io/artifactory/jfrog-cli/v1/${JFROG_VERSION}/jfrog-cli-linux-amd64/jfrog /bin/jfrog

RUN chmod a+x /bin/jfrog
RUN chmod -R 777 /tmp

##############################################################################
##    docker build --no-cache --target certs -t vela-artifactory:certs .    ##
##############################################################################

FROM alpine:3.23.4@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11 as certs

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

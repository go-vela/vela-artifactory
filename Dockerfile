# Copyright (c) 2019 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

FROM scratch

ADD https://api.bintray.com/content/jfrog/jfrog-cli-go/\$latest/jfrog-cli-linux-amd64/jfrog?bt_package=jfrog-cli-linux-amd64 /bin/jfrog

COPY release/vela-artifactory /bin/vela-artifactory

ENTRYPOINT ["/bin/vela-artifactory"]

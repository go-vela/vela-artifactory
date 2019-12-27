// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package artifactory

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
)

type (
	// SearchArgs are arguments for searching files
	SearchArgs struct {
		Path      string
		Recursive bool
	}
)

func (a Artifactory) search(args SearchArgs) ([]utils.ResultItem, error) {
	searchParams := services.NewSearchParams()
	searchParams.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{
		Pattern:   args.Path,
		Recursive: args.Recursive,
	}

	return a.client.SearchFiles(searchParams)
}

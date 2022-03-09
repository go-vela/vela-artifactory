// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
)

type ArtifactoryServicesManager interface {
	UploadFiles(...services.UploadParams) ([]utils.FileInfo, int, int, error)
	Copy(services.MoveCopyParams) (int, int, error)
	GetPathsToDelete(services.DeleteParams) ([]utils.ResultItem, error)
	DeleteFiles([]utils.ResultItem) (int, error)
}

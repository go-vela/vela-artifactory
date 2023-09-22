// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
)

func TestArtifactory_Upload_Exec_Error(t *testing.T) {
	// setup types
	config := &Config{
		Action:   "upload",
		APIKey:   "superSecretAPIKey",
		DryRun:   false,
		Password: "superSecretPassword",
		URL:      "http://localhost:8081/artifactory",
		Username: "octocat",
	}

	cli, err := config.New()
	if err != nil {
		t.Errorf("Unable to create Artifactory client: %v", err)
	}

	u := &Upload{
		Flat:        true,
		IncludeDirs: false,
		Recursive:   true,
		Regexp:      false,
		Path:        "foo/bar",
		Sources:     []string{"baz.txt"},
	}

	err = u.Exec(cli)
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestArtifactory_Upload_Validate(t *testing.T) {
	// setup types
	u := &Upload{
		Flat:        true,
		IncludeDirs: false,
		Recursive:   true,
		Regexp:      false,
		Path:        "foo/bar",
		Sources:     []string{"baz.txt"},
	}

	err := u.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_Upload_Validate_NoPath(t *testing.T) {
	// setup types
	u := &Upload{
		Flat:        true,
		IncludeDirs: false,
		Recursive:   true,
		Regexp:      false,
		Sources:     []string{"baz.txt"},
	}

	err := u.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Upload_Validate_NoSources(t *testing.T) {
	// setup types
	u := &Upload{
		Flat:        true,
		IncludeDirs: false,
		Recursive:   true,
		Regexp:      false,
		Path:        "foo/bar",
	}

	err := u.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Upload_Failure_Returns_Error(t *testing.T) {
	cli := &MockArtifactoryService{Fail: true}

	u := &Upload{
		Flat:        true,
		IncludeDirs: false,
		Recursive:   true,
		Regexp:      false,
		Path:        "foo/bar",
		Sources:     []string{"baz.txt"},
	}

	err := u.Exec(cli)
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestArtifactory_Upload_Success(t *testing.T) {
	cli := &MockArtifactoryService{Fail: false}

	u := &Upload{
		Flat:        true,
		IncludeDirs: false,
		Recursive:   true,
		Regexp:      false,
		Path:        "foo/bar",
		Sources:     []string{"baz.txt"},
	}

	err := u.Exec(cli)
	if err != nil {
		t.Errorf("Exec should have returned err")
	}
}

type MockArtifactoryService struct {
	Fail bool
}

func (a *MockArtifactoryService) UploadFiles(...services.UploadParams) ([]utils.FileInfo, int, int, error) {
	if a.Fail {
		return nil, 0, 1, errors.New("upload failed")
	}

	return nil, 0, 0, nil
}

func (a *MockArtifactoryService) Copy(services.MoveCopyParams) (int, int, error) {
	return 0, 0, nil
}

func (a *MockArtifactoryService) GetPathsToDelete(services.DeleteParams) ([]utils.ResultItem, error) {
	return nil, nil
}

func (a *MockArtifactoryService) DeleteFiles([]utils.ResultItem) (int, error) {
	return 0, nil
}

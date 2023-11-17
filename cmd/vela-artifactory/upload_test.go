// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/go-vela/vela-artifactory/cmd/vela-artifactory/mock"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
)

func TestArtifactory_Plugin_Exec_Upload(t *testing.T) {
	// setup types
	s := httptest.NewServer(mock.Handlers())

	p := &Plugin{
		Config: &Config{
			Action:   "upload",
			Token:    mock.Token,
			APIKey:   mock.APIKey,
			DryRun:   false,
			Password: mock.Password,
			URL:      s.URL,
			Username: mock.Username,
		},
		Copy:    &Copy{},
		Delete:  &Delete{},
		SetProp: &SetProp{},
		Upload: &Upload{
			Flat:        true,
			IncludeDirs: false,
			Recursive:   true,
			Regexp:      false,
			Path:        "foo/bar",
			Sources:     []string{"mock/testdata/baz.txt"},
		},
	}

	err := p.Exec()
	if err != nil {
		t.Errorf("Exec returned err %v", err)
	}
}
func TestArtifactory_Upload_Exec_Error(t *testing.T) {
	// setup types
	config := &Config{
		Action:   "upload",
		APIKey:   mock.APIKey,
		DryRun:   false,
		Password: mock.Password,
		URL:      mock.InvalidArtifactoryServerURL,
		Username: mock.Username,
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

	err = u.Exec(*cli)
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

// func TestArtifactory_Upload_Failure_Returns_Error(t *testing.T) {
// 	cli := &MockArtifactoryService{Fail: true}

// 	u := &Upload{
// 		Flat:        true,
// 		IncludeDirs: false,
// 		Recursive:   true,
// 		Regexp:      false,
// 		Path:        "foo/bar",
// 		Sources:     []string{"baz.txt"},
// 	}

// 	err := u.Exec(*cli)
// 	if err == nil {
// 		t.Errorf("Exec should have returned err")
// 	}
// }

// func TestArtifactory_Upload_Success(t *testing.T) {
// 	cli := &MockArtifactoryService{Fail: false}

// 	u := &Upload{
// 		Flat:        true,
// 		IncludeDirs: false,
// 		Recursive:   true,
// 		Regexp:      false,
// 		Path:        "foo/bar",
// 		Sources:     []string{"baz.txt"},
// 	}

// 	err := u.Exec(*cli)
// 	if err != nil {
// 		t.Errorf("Exec should have returned err")
// 	}
// }

type MockArtifactoryService struct {
	Fail bool
}

func (a *MockArtifactoryService) UploadFiles(...services.UploadParams) (int, int, error) {
	if a.Fail {
		return 0, 1, errors.New("upload failed")
	}

	return 0, 0, nil
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

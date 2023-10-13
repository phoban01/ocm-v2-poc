package github

import (
	"fmt"
	"os"
	"path/filepath"

	v2 "github.com/phoban01/ocm-v2/api/v2"
)

type repository struct {
	owner   string
	repo    string
	url     string
	tmpdir  string
	blobdir string
}

var _ v2.Repository = (*repository)(nil)

func Repository(owner, repo string) (v2.Repository, error) {
	tmpdir, err := os.MkdirTemp("", "ocm-github-*")
	if err != nil {
		return nil, err
	}
	blobdir := filepath.Join(tmpdir, "blobs")
	url := fmt.Sprintf("https://github.com/%s/%s", owner, repo)
	return &repository{
		owner:   owner,
		repo:    repo,
		url:     url,
		tmpdir:  tmpdir,
		blobdir: blobdir,
	}, nil
}

func (r *repository) Context() v2.RepositoryContext {
	return &repositoryContext{
		url: fmt.Sprintf("github://%s", r.url),
	}
}

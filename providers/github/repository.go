package github

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/authn"
	"github.com/phoban01/ocm-v2/api/v2/configr"
)

type repository struct {
	auth     authn.Authenticator
	keychain authn.Keychain
	config   configr.Configuration
	owner    string
	repo     string
	url      string
	tmpdir   string
	blobdir  string
}

var _ v2.Repository = (*repository)(nil)

func Repository(owner, repo string, opts ...Option) (v2.Repository, error) {
	tmpdir, err := os.MkdirTemp("", "ocm-github-*")
	if err != nil {
		return nil, err
	}

	blobdir := filepath.Join(tmpdir, "blobs")

	url := fmt.Sprintf("https://github.com/%s/%s", owner, repo)

	r := &repository{
		owner:   owner,
		repo:    repo,
		url:     url,
		tmpdir:  tmpdir,
		blobdir: blobdir,
	}

	for _, f := range opts {
		f(r)
	}

	if r.auth != nil && r.keychain != nil {
		return nil, errors.New("auth and keychain cannot both be set")
	}

	return r, nil
}

func (r *repository) Context() v2.RepositoryContext {
	return &repositoryContext{
		url: fmt.Sprintf("github://%s", r.url),
	}
}

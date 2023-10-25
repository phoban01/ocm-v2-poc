package oci

import (
	"context"
	"fmt"

	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/authn"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

type repository struct {
	prefix    string
	registry  string
	component string
	auth      *auth.Credential
	client    *auth.Client
}

var _ v2.Repository = (*repository)(nil)

var StoragePrefix = "component-descriptors"

func Repository(registry string, opts ...Option) (v2.Repository, error) {
	r := &repository{
		prefix:   StoragePrefix,
		registry: registry,
	}

	for _, f := range opts {
		if err := f(r); err != nil {
			return nil, err
		}
	}

	client := &auth.Client{
		Client: retry.DefaultClient,
		Cache:  auth.DefaultCache,
		Credential: func(_ context.Context, _ string) (auth.Credential, error) {
			return auth.EmptyCredential, nil
		},
	}

	if r.auth != nil {
		client.Credential = func(_ context.Context, _ string) (auth.Credential, error) {
			return *r.auth, nil
		}
	}

	r.client = client

	return r, nil
}

type Option func(r *repository) error

func WithCredentials(creds authn.Authenticator) func(r *repository) error {
	return func(r *repository) error {
		cfg, err := creds.Authorization()
		if err != nil {
			return err
		}
		r.auth = &auth.Credential{
			Username:     cfg.Username,
			Password:     cfg.Password,
			AccessToken:  cfg.AccessToken,
			RefreshToken: cfg.RefreshToken,
		}
		return nil
	}
}

func (r *repository) Context() v2.RepositoryContext {
	return &repositoryContext{
		url: r.registry,
	}
}

func (r *repository) url() string {
	return fmt.Sprintf("%s/%s/%s", r.registry, r.prefix, r.component)
}

func (r *repository) storage() (*remote.Repository, error) {
	repo, err := remote.NewRepository(r.url())
	if err != nil {
		return nil, err
	}
	repo.Client = r.client
	return repo, nil
}

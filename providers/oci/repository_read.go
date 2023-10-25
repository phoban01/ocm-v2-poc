package oci

import (
	"context"
	"encoding/json"
	"io"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

func (r *repository) ReadBlob(ctx context.Context, digest string) (io.ReadCloser, error) {
	store, err := r.storage()
	if err != nil {
		return nil, err
	}
	_, data, err := store.FetchReference(ctx, digest)
	return data, err
}

func (r *repository) Get(name, version string) (v2.Component, error) {
	ctx := context.TODO()
	r.component = name

	store, err := r.storage()
	if err != nil {
		return nil, err
	}

	_, content, err := store.FetchReference(ctx, version)
	if err != nil {
		return nil, err
	}
	defer content.Close()

	data, err := io.ReadAll(content)
	if err != nil {
		return nil, err
	}

	var manifest ocispec.Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	config, err := store.Fetch(ctx, manifest.Config)
	if err != nil {
		return nil, err
	}
	defer config.Close()

	configData, err := io.ReadAll(config)
	if err != nil {
		return nil, err
	}

	desc := types.Descriptor{}
	if err := json.Unmarshal(configData, &desc); err != nil {
		return nil, err
	}

	return &component{
		repository: r,
		descriptor: desc,
	}, nil
}

func (r *repository) List() ([]v2.Component, error) {
	return nil, nil
}

func (r *repository) Delete() error {
	return nil
}

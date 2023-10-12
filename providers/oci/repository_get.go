package oci

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

func (r *repository) Get(componentName, version string) (v2.Component, error) {
	url := fmt.Sprintf(
		"%s/component-descriptors/%s:%s",
		r.registry,
		componentName,
		version,
	)

	ref, err := name.ParseReference(url)
	if err != nil {
		return nil, err
	}

	img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return nil, err
	}

	manifest, err := img.Manifest()
	if err != nil {
		return nil, err
	}

	cd, err := name.NewDigest(fmt.Sprintf("%s@%s", url, manifest.Layers[0].Digest.String()))
	if err != nil {
		return nil, err
	}

	layer, err := remote.Layer(cd, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return nil, err
	}

	reader, err := layer.Uncompressed()
	if err != nil {
		return nil, err
	}

	layerData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	desc := types.Descriptor{}
	if err := json.Unmarshal(layerData, &desc); err != nil {
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

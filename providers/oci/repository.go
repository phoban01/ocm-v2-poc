package oci

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	gmutate "github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/static"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/google/go-containerregistry/pkg/v1/types"
	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/mutate"
)

type repository struct {
	registry string
}

var _ v2.Repository = (*repository)(nil)

func Repository(registry string) (v2.Repository, error) {
	return &repository{registry: registry}, nil
}

func (s *repository) Context() *v2.RepositoryContext {
	return &v2.RepositoryContext{
		Type: "ocm.repository/oci",
		URL:  fmt.Sprintf("oci://%s", s.registry),
	}
}

func (s *repository) Get(name, version string) (v2.Component, error) {
	return nil, nil
}

func (s *repository) List() ([]v2.Component, error) {
	return nil, nil
}

func (s *repository) Delete() error {
	return nil
}

func (s *repository) Write(component v2.Component) error {
	// update the repository context
	component = mutate.WithRepositoryContext(component, s)

	desc, err := component.Descriptor()
	if err != nil {
		return err
	}

	url := fmt.Sprintf(
		"%s/component-descriptors/%s:%s",
		s.registry,
		desc.Name,
		desc.Version,
	)

	resources, err := component.Resources()
	if err != nil {
		return err
	}

	desc, err = component.Descriptor()
	if err != nil {
		return err
	}

	ref, err := name.ParseReference(url)
	if err != nil {
		return err
	}

	data, err := json.Marshal(desc)
	if err != nil {
		return err
	}

	l := static.NewLayer(data, "ocm.software/vnd.ocm.software.component-descriptor")
	if err != nil {
		return err
	}

	img := gmutate.MediaType(empty.Image, types.OCIManifestSchema1)
	img, err = gmutate.Append(img, gmutate.Addendum{Layer: l})
	if err != nil {
		return err
	}

	resources, err = component.Resources()
	if err != nil {
		return err
	}

	// TODO: ensure access is updated correctly; this pushes the layer
	// but doesn't update the location or component descriptor.
	for _, r := range resources {
		if r.Deferrable() {
			continue
		}

		layer, err := tarball.LayerFromOpener(
			r.Access().Data,
			tarball.WithMediaType("application/tar-gzip"),
		)
		if err != nil {
			return err
		}

		img, err = gmutate.Append(img, gmutate.Addendum{Layer: layer})
		if err != nil {
			return err
		}
	}

	if err := remote.Write(ref, img, remote.WithAuthFromKeychain(authn.DefaultKeychain)); err != nil {
		return err
	}

	return nil
}

type rawManifest struct {
	body      []byte
	mediaType types.MediaType
}

func (r *rawManifest) RawManifest() ([]byte, error) {
	return r.body, nil
}

func (r *rawManifest) MediaType() (types.MediaType, error) {
	return r.mediaType, nil
}

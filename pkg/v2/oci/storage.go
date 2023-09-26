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
	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/mutate"
)

type storage struct {
	registry string
}

var _ v2.Storage = (*storage)(nil)

func Storage(registry string) (v2.Storage, error) {
	return &storage{registry: registry}, nil
}

func (s *storage) Context() *v2.StorageContext {
	return &v2.StorageContext{
		Type: "ocm.storage/oci",
		URL:  fmt.Sprintf("oci://%s", s.registry),
	}
}

func (s *storage) Write(component v2.Component) error {
	desc, err := component.Descriptor()
	if err != nil {
		return err
	}

	url := fmt.Sprintf(
		"%s/component-descriptors/%s:%s",
		s.registry,
		desc.Metadata.Name,
		desc.Version,
	)

	resources, err := component.Resources()
	if err != nil {
		return err
	}

	for _, r := range resources {
		if r.Deferrable() {
			continue
		}

		dig, err := r.Digest()
		if err != nil {
			return err
		}

		access := fmt.Sprintf("%s@sha256:%s", url, dig)

		r = mutate.SetAccess(r, access)
		component = mutate.ReplaceResource(component, r)
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

	for _, r := range resources {
		if r.Deferrable() {
			continue
		}

		layer, err := tarball.LayerFromOpener(r.Blob, tarball.WithMediaType("application/tar-gzip"))
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

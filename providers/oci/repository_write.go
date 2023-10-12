package oci

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	gmutate "github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/static"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	gtypes "github.com/google/go-containerregistry/pkg/v1/types"
	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/mutate"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

func (r *repository) Write(component v2.Component) error {
	desc, err := component.Descriptor()
	if err != nil {
		return err
	}

	url := fmt.Sprintf(
		"%s/component-descriptors/%s:%s",
		r.registry,
		desc.ObjectMeta.Name,
		desc.ObjectMeta.Version,
	)

	ref, err := name.ParseReference(url)
	if err != nil {
		return err
	}

	resources, err := component.Resources()
	if err != nil {
		return err
	}

	pusher, err := remote.NewPusher(remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return err
	}

	processedItems := make([]v2.Resource, len(resources))
	layersToAdd := make([]v1.Layer, 0)

	for i, item := range resources {
		if item.Deferrable() {
			processedItems[i] = item
			continue
		}

		acc, err := item.Access()
		if err != nil {
			return err
		}

		layer, err := tarball.LayerFromOpener(
			acc.Data,
			tarball.WithMediaType("application/tar-gzip"),
			tarball.WithCompressedCaching,
		)
		if err != nil {
			return err
		}

		if err := pusher.Upload(context.Background(), ref.Context(), layer); err != nil {
			return err
		}

		layersToAdd = append(layersToAdd, layer)

		dig, err := layer.Digest()
		if err != nil {
			return err
		}

		item = mutate.WithAccess(item, &localAccess{
			mediaType: "application/tar-gzip",
			digest: types.Digest{
				HashAlgorithm:          "sha256",
				NormalisationAlgorithm: "json/v1",
				Value:                  dig.String(),
			},
		})

		processedItems[i] = item
	}

	component = mutate.WithResources(component, processedItems...)

	desc, err = component.Descriptor()
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

	img := gmutate.MediaType(empty.Image, gtypes.OCIManifestSchema1)
	img, err = gmutate.Append(img, gmutate.Addendum{Layer: l})
	if err != nil {
		return err
	}

	for _, l := range layersToAdd {
		img, err = gmutate.Append(img, gmutate.Addendum{Layer: l})
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

func (r *repository) ReadBlob(digest string) (v2.Access, error) {
	return nil, nil
}

func (r *repository) WriteBlob(v2.Access) error {
	return nil
}

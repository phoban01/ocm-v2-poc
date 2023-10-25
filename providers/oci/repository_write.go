package oci

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	digest "github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/mutate"
	"github.com/phoban01/ocm-v2/api/v2/types"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/registry/remote"
)

// digest needs to come from access and should be a ocispecDigest
// we should probably store the length on the access also.
func (r *repository) WriteBlob(ctx context.Context, acc v2.Access) (v2.Access, error) {
	dig, err := acc.Digest()
	if err != nil {
		return nil, err
	}

	data, err := acc.Data()
	if err != nil {
		return nil, err
	}
	defer data.Close()

	size, err := acc.Length()
	if err != nil {
		return nil, err
	}

	desc := ocispec.Descriptor{
		MediaType: "application/tar-gzip",
		Digest:    digest.NewDigestFromEncoded(digest.SHA256, dig.Value),
		Size:      size,
	}

	store, err := r.storage()
	if err != nil {
		return nil, err
	}

	if err := store.Push(ctx, desc, data); err != nil {
		return nil, err
	}

	return &localBlob{
		mediaType: "application/tar-gzip",
		length:    size,
		digest: types.Digest{
			HashAlgorithm:          "sha256",
			NormalisationAlgorithm: "json/v1",
			Value:                  dig.Value,
		},
	}, nil
}

func (r *repository) Write(ctx context.Context, component v2.Component) error {
	desc, err := component.Descriptor()
	if err != nil {
		return err
	}

	r.component = desc.Name

	resources, err := component.Resources()
	if err != nil {
		return err
	}

	processedItems := make([]v2.Resource, len(resources))

	layers := make([]ocispec.Descriptor, 0)

	for i, item := range resources {
		if item.Deferrable() {
			// need to handle deferrables here
			item, err := r.handleDeferrable(ctx, item)
			if err != nil {
				return err
			}

			processedItems[i] = item
			continue
		}

		acc, err := item.Access()
		if err != nil {
			return err
		}

		size, err := acc.Length()
		if err != nil {
			return err
		}

		dig, err := acc.Digest()
		if err != nil {
			return err
		}

		acc, err = r.WriteBlob(ctx, acc)
		if err != nil {
			return err
		}

		layers = append(layers, ocispec.Descriptor{
			MediaType: acc.MediaType(),
			Digest:    digest.NewDigestFromEncoded(digest.SHA256, dig.Value),
			Size:      size,
			Annotations: map[string]string{
				"ocm.software/resource-name": item.Name(),
				"ocm.software/resource-type": string(item.Type()),
				"ocm.software/access-type":   acc.Type(),
			},
		})

		item = mutate.WithAccess(item, acc)

		processedItems[i] = item
	}

	component = mutate.WithResources(component, processedItems...)

	desc, err = component.Descriptor()
	if err != nil {
		return err
	}

	configBlob, err := json.Marshal(desc)
	if err != nil {
		return err
	}

	configDesc := content.NewDescriptorFromBytes(
		"application/vnd.ocm.software.component.config.v1+json",
		configBlob,
	)

	manifestBlob, err := generateManifest(configDesc, layers...)
	if err != nil {
		return err
	}

	manifestDesc := content.NewDescriptorFromBytes(ocispec.MediaTypeImageManifest, manifestBlob)

	store, err := r.storage()
	if err != nil {
		return err
	}

	if err := store.Push(ctx, configDesc, bytes.NewReader(configBlob)); err != nil {
		return err
	}

	if err := store.PushReference(ctx, manifestDesc, bytes.NewReader(manifestBlob), desc.Version); err != nil {
		return err
	}

	return nil
}

func generateManifest(config ocispec.Descriptor, layers ...ocispec.Descriptor) ([]byte, error) {
	content := ocispec.Manifest{
		Config:    config,
		Layers:    layers,
		Versioned: specs.Versioned{SchemaVersion: 2},
	}
	return json.Marshal(content)
}

func (r *repository) handleDeferrable(ctx context.Context, item v2.Resource) (v2.Resource, error) {
	switch item.Type() {
	case "ociImage":
		version := item.Version()
		acc, err := item.Access()
		if err != nil {
			return nil, err
		}

		srcRepo := strings.TrimSuffix(acc.Reference(), ":"+version)
		targetRef := fmt.Sprintf("%s/%s/%s", r.registry, r.component, srcRepo)
		store, err := remote.NewRepository(targetRef)
		if err != nil {
			return nil, err
		}
		store.Client = r.client

		src, err := remote.NewRepository(srcRepo)
		if err != nil {
			return nil, err
		}

		_, err = oras.Copy(ctx, src, version, store, version, oras.DefaultCopyOptions)
		if err != nil {
			return nil, err
		}

		acc = &accessor{
			repository: r,
			mediaType:  acc.MediaType(),
			ref:        fmt.Sprintf("%s:%s", targetRef, version),
		}

		return mutate.WithAccess(item, acc), nil
	}
	return item, nil
}

package main

import (
	"context"
	"fmt"
	"log"

	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/build"
	"github.com/phoban01/ocm-v2/api/v2/mutate"
	"github.com/phoban01/ocm-v2/api/v2/types"
	"github.com/phoban01/ocm-v2/providers/filesystem"
	"github.com/phoban01/ocm-v2/providers/helm"
	"github.com/phoban01/ocm-v2/providers/oci"
)

func main() {
	config, err := NewFileResource("config", "config.yaml", "application/x-yaml")
	if err != nil {
		log.Fatal(err)
	}

	image, err := NewImageResource("nginx", "docker.io/nginx", "1.25.2")
	if err != nil {
		log.Fatal(err)
	}

	chart, err := NewChartResource("chart", "nginx-stable/nginx-ingress", "0.17.1")
	if err != nil {
		log.Fatal(err)
	}

	resources := []v2.Resource{config, image, chart}

	// create a new component
	cmp := build.New("ocm.software/v2/server", "v1.0.0", "acme.org")

	// add the resources to the component
	cmp = mutate.WithResources(cmp, resources...)

	// setup the repository
	repo, err := oci.Repository("ghcr.io/phoban01")
	if err != nil {
		log.Fatal(err)
	}

	// write the component to the archive
	if err := repo.Write(context.TODO(), cmp); err != nil {
		log.Fatal(err)
	}
}

func NewFileResource(name, path, mediaType string) (v2.Resource, error) {
	meta := types.ObjectMeta{
		Name: name,
		Type: types.Blob,
		Labels: map[string]string{
			"ocm.software/filename": path,
		},
	}

	access, err := filesystem.FromFile(path, filesystem.WithMediaType(mediaType))
	if err != nil {
		return nil, err
	}

	return build.NewResource(meta, access), nil
}

func NewImageResource(name, ref, version string) (v2.Resource, error) {
	meta := types.ObjectMeta{
		Name:    name,
		Type:    types.OCIImage,
		Version: version,
	}

	access, err := oci.FromImage(fmt.Sprintf("%s:%s", ref, version))
	if err != nil {
		return nil, err
	}

	return build.NewResource(meta, access, build.Deferrable(true)), nil
}

func NewChartResource(name, ref, version string) (v2.Resource, error) {
	meta := types.ObjectMeta{
		Name:    name,
		Type:    types.HelmChart,
		Version: version,
	}

	access, err := helm.FromRepository(ref, version)
	if err != nil {
		return nil, err
	}

	return build.NewResource(meta, access), nil
}

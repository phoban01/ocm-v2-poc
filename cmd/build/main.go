package main

import (
	"log"

	"github.com/phoban01/ocm-v2/api/v2/archive"
	"github.com/phoban01/ocm-v2/api/v2/builder"
	"github.com/phoban01/ocm-v2/api/v2/mutate"
	"github.com/phoban01/ocm-v2/api/v2/types"
	"github.com/phoban01/ocm-v2/providers/blob"
	"github.com/phoban01/ocm-v2/providers/oci"
)

func main() {
	// define config metadata
	configMeta := types.ObjectMeta{
		Name: "config",
		Type: types.Blob,
	}

	// get the config blob accessor
	configAccess, err := blob.FromFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// define the config resource
	config := builder.NewResource(configMeta, configAccess)

	// define the image metadata
	imageMeta := types.ObjectMeta{
		Name: "image",
		Type: types.OCIImage,
	}

	// get the oci image
	imageAcc, err := oci.FromImage("docker.io/redis:latest")
	if err != nil {
		log.Fatal(err)
	}

	// define the image resource, (builder.Deferrable means the resource does not need to be read at build time)
	image := builder.NewResource(imageMeta, imageAcc, builder.Deferrable(true))

	// create a new component
	cmp := builder.New("ocm.software/test", "v1.0.0", "acme.org")

	// add the resources to the component
	cmp = mutate.WithResources(cmp, config, image)

	// setup the repository
	repo, err := archive.Repository("transport-archive")
	if err != nil {
		log.Fatal(err)
	}

	// write the component to the archive
	if err := repo.Write(cmp); err != nil {
		log.Fatal(err)
	}
}

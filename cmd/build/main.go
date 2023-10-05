package main

import (
	"log"

	"github.com/phoban01/ocm-v2/pkg/providers/blob"
	"github.com/phoban01/ocm-v2/pkg/providers/oci"
	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/archive"
	"github.com/phoban01/ocm-v2/pkg/v2/builder"
	"github.com/phoban01/ocm-v2/pkg/v2/mutate"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

func main() {
	// create a new component
	cmp := builder.New("ocm.software/test", "v1.0.0", "acme.org")

	// define file access
	config, err := blob.FromFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// define image access
	image, err := oci.FromImage("docker.io/redis:latest")
	if err != nil {
		log.Fatal(err)
	}

	// gather resources
	resources := []v2.Resource{
		builder.NewResource(types.ObjectMeta{
			Name: "config",
			Type: types.Blob,
		}, config),
		builder.NewResource(types.ObjectMeta{
			Name: "image",
			Type: types.OCIImage,
		}, image, builder.Deferrable(true)),
	}

	// add the resources to the component
	cmp = mutate.WithResources(cmp, resources...)

	// setup the repository
	ctf, err := archive.Repository("./test-ctf")
	if err != nil {
		log.Fatal(err)
	}

	// write the component to the repository
	if err := ctf.Write(cmp); err != nil {
		log.Fatal(err)
	}
}

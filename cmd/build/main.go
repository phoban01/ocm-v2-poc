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
	// define config resource
	configFile, err := blob.FromFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	config := builder.NewResource(types.ObjectMeta{
		Name: "config",
		Type: types.Blob,
	}, configFile)

	// define image resource
	imageAcc, err := oci.FromImage("docker.io/redis:latest")
	if err != nil {
		log.Fatal(err)
	}

	image := builder.NewResource(types.ObjectMeta{
		Name: "image",
		Type: types.OCIImage,
	}, imageAcc, builder.Deferrable(true))

	// define helm resource
	// helmAcc, err := helm.FromChart("nginx-stable/nginx-ingress", "0.17.1")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// chart := builder.NewResource(types.ObjectMeta{
	// 	Name: "chart",
	// 	Type: types.HelmChart,
	// }, helmAcc)

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

package main

import (
	"log"
	"os"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/archive"
	"github.com/phoban01/ocm-v2/pkg/v2/blob"
	"github.com/phoban01/ocm-v2/pkg/v2/builder"
	"github.com/phoban01/ocm-v2/pkg/v2/mutate"
	"github.com/phoban01/ocm-v2/pkg/v2/oci"
)

func main() {
	// create a new component
	cmp := builder.New("ocm.software/test", "v1.0.0", "acme.org")

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// create resources
	resources := []v2.Resource{
		blob.FromBytes("data", data),
		blob.FromFile("config", "config.yaml"),
		oci.Resource("web-server", "docker.io/nginx:1.25.2"),
		oci.Resource("redis", "docker.io/redis:latest"),
	}

	// add the resource to the component
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

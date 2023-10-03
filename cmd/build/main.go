package main

import (
	"log"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/archive"
	"github.com/phoban01/ocm-v2/pkg/v2/builder"
	"github.com/phoban01/ocm-v2/pkg/v2/file"
	"github.com/phoban01/ocm-v2/pkg/v2/mutate"
	"github.com/phoban01/ocm-v2/pkg/v2/oci"
)

func main() {
	// create a new component
	cmp := builder.New("ocm.software/test", "v1.0.0", "acme.org")

	// create resources
	resources := []v2.Resource{
		file.Resource("data", "config.yaml"),
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

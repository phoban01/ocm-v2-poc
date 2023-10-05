# OCM V2

This is an experimental re-imagining of how the core OCM library could look.

## **⚠️ this is not a feature complete implementation, it is a draft**

It is heavily inspired by the structure of go-containerregistry. 

The API provides some basic primitives which can be composed to create components.

The following example shows how we can build a component and write to a repository:

```golang
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
```
## Operations

![image](https://github.com/phoban01/ocm-v2/assets/4415593/9e15a2c8-a7e5-4742-89fb-8ee10fb8d091)

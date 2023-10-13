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

	"github.com/phoban01/ocm-v2/api/v2/build"
	"github.com/phoban01/ocm-v2/api/v2/mutate"
	"github.com/phoban01/ocm-v2/api/v2/types"
	"github.com/phoban01/ocm-v2/providers/filesystem"
)

func main() {
	// define metadata for the resource
	meta := types.ObjectMeta{
		Name: "config",
		Type: types.ResourceType("file"),
	}

	// create an access method for a file on disk
	// notice the filesystem provider helper methods to access resources
	// filesystem.ReadFile returns v2.Access
	access, err := filesystem.ReadFile(
		"config.yaml",
		filesystem.WithMediaType("application/x-yaml"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// build the config resource using the metadata and access method
	config := build.NewResource(meta, access)

	// create the component
	cmp := build.New("ocm.software/test", "v1.0.0", "acme.org")

	// add resources to the component using the mutate package
	cmp = mutate.WithResources(cmp, config)

	// setup the repository using the filesystem provider
	repo, err := filesystem.Repository("./transport-archive")
	if err != nil {
		log.Fatal(err)
	}

	// write the component to the repository
	if err := repo.Write(cmp); err != nil {
		log.Fatal(err)
	}
}
```
## Operations

![image](https://github.com/phoban01/ocm-v2/assets/4415593/9e15a2c8-a7e5-4742-89fb-8ee10fb8d091)

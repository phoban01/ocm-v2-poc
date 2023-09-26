# OCM V2

This is an experimental re-imagining of how the core OCM library could look.

It is heavily inspired by the structure of go-containerregistry. 

The API provides some basic primitives which can be composed to create components.

Access methods are provided by satisfying the `Resource` interface.

The following example shows how we can build a component and write to an OCI repository:

```golang
package main

import (
	"log"

	"github.com/phoban01/ocm-v2/pkg/signer"
	v2 "github.com/phoban01/ocm-v2/pkg/v2"
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
	}

	// add the resource to the component
	cmp = mutate.AddResources(cmp, resources...)

	// get a list of the components' resources for signing
	signables, err := cmp.Resources()
	if err != nil {
		log.Fatal(err)
	}

	// create the signer
	sig, err := signer.New("data-sig", "rsa.priv", signables...)
	if err != nil {
		log.Fatal(err)
	}

	// add the signatures the component
	cmp = mutate.AddSignatures(cmp, sig)

	// create a bundle (component archive) to store the component
	store, err := oci.Storage("ghcr.io/phoban01")
	if err != nil {
		log.Fatal(err)
	}

	// update the storage context
	cmp = mutate.AddStorageContext(cmp, store)

	// write the component
	if err := store.Write(cmp); err != nil {
		log.Fatal(err)
	}
}
```

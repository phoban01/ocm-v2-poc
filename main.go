package main

import (
	"log"

	"github.com/phoban01/ocm-v2/pkg/signer"
	"github.com/phoban01/ocm-v2/pkg/v2/builder"
	"github.com/phoban01/ocm-v2/pkg/v2/bundle"
	"github.com/phoban01/ocm-v2/pkg/v2/file"
	"github.com/phoban01/ocm-v2/pkg/v2/mutate"
)

func main() {
	// create a new component
	cmp := builder.New("ocm.software/test", "v1.0.0", "acme.org")

	// create a new resource
	resource := file.New("data", "config.yaml")

	// add the resource to the component
	cmp = mutate.AddResources(cmp, resource)

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
	store, err := bundle.New("./component-bundle")
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

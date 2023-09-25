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
	cmp := builder.New("test", "v1.0.0", "acme.org")

	resource := file.New("data", "config.yaml")

	cmp = mutate.AddResources(cmp, resource)

	signables, err := cmp.Resources()
	if err != nil {
		log.Fatal(err)
	}

	sig, err := signer.New("data-sig", "rsa.priv", signables...)
	if err != nil {
		log.Fatal(err)
	}

	cmp = mutate.AddSignatures(cmp, sig)

	if err := bundle.Write("./component-bundle", cmp); err != nil {
		log.Fatal(err)
	}
}

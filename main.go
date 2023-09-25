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
	c := builder.New("test", "v1.0.0", "acme.org")

	x := file.New("data", "config.yaml")

	c = mutate.AddResources(c, x)

	signables, err := c.Resources()
	if err != nil {
		log.Fatal(err)
	}

	sig, err := signer.New("data-sig", "rsa.priv", signables...)
	if err != nil {
		log.Fatal(err)
	}

	c = mutate.AddSignatures(c, sig)

	if err := bundle.Write("./component-bundle", c); err != nil {
		log.Fatal(err)
	}
}

# OCM V2

This is an experimental re-imagining of how the core OCM library could look.

It is heavily inspired by the structure of go-containerregistry. 

The API provides some primitives which can be composed to create components.

The following example shows how we can build a component and write a bundle to a local filesystem directory:

```golang
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
```

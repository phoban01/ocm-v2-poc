# OCM V2

This is an experimental re-imagining of how the core OCM library could look.

It is heavily inspired by the structure of go-containerregistry. 

The API provides some basic primitives which can be composed to create components.

Access methods are provided by satisfying the `Resource` interface.

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
```

To run:

```shell
go run main.go && cat ./component-bundle/component-descriptor.json | jq
 
```
Results in:

```json
{
  "metadata": {
    "name": "test"
  },
  "provider": {
    "name": "acme.org"
  },
  "version": "v1.0.0",
  "resources": [
    {
      "name": "data",
      "labels": {
        "ocm.software/filename": "config.yaml"
      },
      "digest": "9d39bb26e078077722e7ed4fab078fc926fd6d9701f1adb67b35c2e227657d03",
      "media_type": "application/x-yaml",
      "resource_type": "file"
    }
  ],
  "signatures": [
    {
      "name": "data-sig",
      "digest": "9173640566a08d88ccbf40c8fcc8a054b26c1494db7c74522a9dc1210cc461799aa155d6746994fc77f6ece5b579c86e921d9b0b94e30e13e9e1093c8e37f5af45028323d9b1787f5a36fff1833a2e55a61e15de986df784d7d6348b1dfa285b989ceed8ec34e4fcde1725e5d458b016a3fa6b29c6014cb2fd4157cb74f920f28235ae03762cf734976ad177881f05f85ede0b58dcfc6a9982db656112ad4701335f80eedeeb38944d4f81801f73fd876a69e917e25949f9d5b7a9f3316443817946fc3f21d04801a3fc1442bffede15e9cf1b3c93156d82be1b16835d9da0d11c56844d3841c6952301cfc4f699c8ab66252d7bad85ca309277f4cb2140d575"
    }
  ]
}
```

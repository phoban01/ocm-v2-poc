package main

import (
	"context"
	"log"

	"github.com/phoban01/ocm-v2/providers/filesystem"
	"github.com/phoban01/ocm-v2/providers/oci"
)

func main() {
	repo, err := oci.Repository("ghcr.io/phoban01")
	if err != nil {
		log.Fatal(err)
	}

	cmp, err := repo.Get("ocm.software/v2/server", "v1.0.0")
	if err != nil {
		log.Fatal(err)
	}

	archive, err := filesystem.Repository("local-copy")
	if err != nil {
		log.Fatal(err)
	}

	if err := archive.Write(context.TODO(), cmp); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"log"

	"github.com/phoban01/ocm-v2/providers/filesystem"
	"github.com/phoban01/ocm-v2/providers/oci"
)

func main() {
	repo, err := oci.Repository("ghcr.io/phoban01/mytest")
	if err != nil {
		log.Fatal(err)
	}

	cmp, err := repo.Get("ocm.software/piaras", "v5.0.0")
	if err != nil {
		log.Fatal(err)
	}

	archive, err := filesystem.Repository("local-copy")
	if err != nil {
		log.Fatal(err)
	}

	if err := archive.Write(cmp); err != nil {
		log.Fatal(err)
	}

	resources, err := cmp.Resources()
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range resources {
		fmt.Println(item.Name())
	}
}

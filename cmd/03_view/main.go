package main

import (
	"encoding/json"
	"fmt"
	"log"

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

	desc, err := cmp.Descriptor()
	if err != nil {
		log.Fatal(err)
	}

	out, err := json.MarshalIndent(desc, " ", " ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(out))
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/phoban01/ocm-v2/api/v2/authn"
	"github.com/phoban01/ocm-v2/providers/oci"
)

func main() {
	creds := &authn.Basic{
		Username: "phoban01",
		Password: os.Getenv("GITHUB_TOKEN"),
	}

	repo, err := oci.Repository("ghcr.io/phoban01", oci.WithCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}

	cmp, err := repo.Get("ocm.software/just-in-time", "v2.0.0")
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

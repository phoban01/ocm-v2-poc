package main

import (
	"fmt"
	"io"
	"log"

	"github.com/phoban01/ocm-v2/api/v2/query"
	"github.com/phoban01/ocm-v2/api/v2/types"
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

	config, err := query.GetResourceData(cmp, types.ObjectMeta{
		Name: "config",
		Type: types.Blob,
	})

	buffer := make([]byte, 4096) // Adjust size as needed
	for {
		n, err := config.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println("Error:", err)
			return
		}

		if n == 0 {
			break
		}

		fmt.Print(string(buffer[:n]))
	}
}

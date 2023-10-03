package main

import (
	"fmt"
	"log"

	"github.com/phoban01/ocm-v2/pkg/v2/archive"
)

func main() {
	ctf, err := archive.Repository("test-ctf")
	if err != nil {
		log.Fatal(err)
	}

	cmp, err := ctf.Get("ocm.software/test", "v1.0.0")
	if err != nil {
		log.Fatal(err)
	}

	res, err := cmp.Resources()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("TYPE\tNAME\tDIGEST")
	for _, v := range res {
		dig, err := v.Digest()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\t%s\t%s\n", v.ResourceType(), v.Name(), dig.Value)
	}
}

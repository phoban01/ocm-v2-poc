package main

import (
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/phoban01/ocm-v2/providers/helm"
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

	res, err := cmp.Resources()
	if err != nil {
		log.Fatal(err)
	}

	var table [][]string
	table = append(table, []string{"TYPE", "NAME", "MEDIA TYPE", "DIGEST"})
	for _, v := range res {
		dig, err := v.Digest()
		if err != nil {
			log.Fatal(err)
		}

		acc, err := v.Access()
		if err != nil {
			log.Fatal(err)
		}
		row := []string{string(v.Type()), v.Name(), acc.MediaType(), dig.Value}
		table = append(table, row)
	}

	printTable(table)

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

func printTable(data [][]string) {
	if len(data) == 0 {
		fmt.Println("Table is empty")
		return
	}

	// Calculate the maximum width for each column
	colWidths := make([]int, len(data[0]))
	for _, row := range data {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Print the table header
	for i, cell := range data[0] {
		fmt.Printf("%-*s", colWidths[i]+2, cell) // +2 for padding
	}
	fmt.Println()

	// Print the table data
	for _, row := range data[1:] {
		for i, cell := range row {
			fmt.Printf("%-*s", colWidths[i]+2, cell) // +2 for padding
		}
		fmt.Println()
	}
}

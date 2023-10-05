package main

import (
	"fmt"
	"log"

	"github.com/phoban01/ocm-v2/pkg/providers/myblob"
	"github.com/phoban01/ocm-v2/pkg/providers/myoci"
	"github.com/phoban01/ocm-v2/pkg/v2/archive"
)

func main() {
	myoci.Use()
	myblob.Use()

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

	var table [][]string
	table = append(table, []string{"TYPE", "NAME", "MEDIA TYPE", "DIGEST"})
	for _, v := range res {
		dig, err := v.Digest()
		if err != nil {
			log.Fatal(err)
		}
		table = append(
			table,
			[]string{string(v.Type()), v.Name(), v.Access().MediaType(), dig.Value},
		)
	}

	printTable(table)
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

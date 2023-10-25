package main

import (
	"context"
	"log"

	"github.com/phoban01/ocm-v2/api/v2/build"
	"github.com/phoban01/ocm-v2/api/v2/mutate"
	"github.com/phoban01/ocm-v2/api/v2/types"
	fs "github.com/phoban01/ocm-v2/providers/filesystem"
)

func main() {
	// define metadata for the resource
	meta := types.ObjectMeta{
		Name: "config",
		Type: types.ResourceType("file"),
	}

	// create an access method for a file on disk
	// notice the filesystem provider helper methods to access resources
	// filesystem.ReadFile returns v2.Access
	content, err := fs.FromFile("config.yaml", fs.WithMediaType("application/x-yaml"))
	if err != nil {
		log.Fatal(err)
	}

	// build the config resource using the metadata and access method
	config := build.NewResource(meta, content)

	// create the component
	cfg := build.Component{
		Name:     "ocm.software/test",
		Version:  "v1.0.0",
		Provider: "acme.org",
	}

	// add resources to the component using the mutate package
	cmp := mutate.WithResources(cfg.New(), config)

	// setup the repository using the filesystem provider
	repo, err := fs.Repository("./transport-archive")
	if err != nil {
		log.Fatal(err)
	}

	// write the component to the repository
	if err := repo.Write(context.TODO(), cmp); err != nil {
		log.Fatal(err)
	}
}

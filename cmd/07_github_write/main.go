package main

import (
	"log"

	"github.com/phoban01/ocm-v2/api/v2/build"
	"github.com/phoban01/ocm-v2/api/v2/mutate"
	"github.com/phoban01/ocm-v2/api/v2/types"
	"github.com/phoban01/ocm-v2/providers/filesystem"
	"github.com/phoban01/ocm-v2/providers/github"
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
	access, err := filesystem.ReadFile(
		"config.yaml",
		filesystem.WithMediaType("application/x-yaml"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// build the config resource using the metadata and access method
	config := build.NewResource(meta, access)

	// create the component
	cmp := build.New("ocm.software/test", "v1.0.0", "acme.org")

	// add resources to the component using the mutate package
	cmp = mutate.WithResources(cmp, config)

	// setup the repository using the github provider
	// arguments are owner and repository
	repo, err := github.Repository("phoban01", "ocm-github-repository")
	if err != nil {
		log.Fatal(err)
	}

	// write the component to the repository
	if err := repo.Write(cmp); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"context"
	"log"
	"os"

	"github.com/phoban01/ocm-v2/api/v2/authn"
	"github.com/phoban01/ocm-v2/api/v2/build"
	"github.com/phoban01/ocm-v2/api/v2/configr"
	"github.com/phoban01/ocm-v2/api/v2/mutate"
	"github.com/phoban01/ocm-v2/api/v2/types"
	"github.com/phoban01/ocm-v2/providers/filesystem"
	"github.com/phoban01/ocm-v2/providers/github"
)

func main() {
	// define metadata for the resource
	meta := types.ObjectMeta{
		Name: "kubernetes-config",
		Type: types.ResourceType("file"),
	}

	// create an access method for a file on disk
	// notice the filesystem provider helper methods to access resources
	// filesystem.ReadFile returns v2.Access
	access, err := filesystem.FromFile(
		"config.yaml",
		filesystem.WithMediaType("application/x-yaml"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// build the config resource using the metadata and access method
	resource := build.NewResource(meta, access)

	// create the component
	cmp := build.New("ocm.software/test", "v1.2.0", "acme.org")

	// add resources to the component using the mutate package
	cmp = mutate.WithResources(cmp, resource)

	auth := &authn.Bearer{
		Token: os.Getenv("GITHUB_TOKEN"),
	}

	config := &configr.StaticConfig{
		"commit.message": "[ocm] automatically added via ocm",
	}

	opts := []github.Option{
		github.WithAuth(auth),
		github.WithConfig(config),
	}

	// setup the repository using the github provider
	// arguments are owner and repository
	repo, err := github.Repository("phoban01", "ocm-github-repository", opts...)
	if err != nil {
		log.Fatal(err)
	}

	// write the component to the repository
	if err := repo.Write(context.TODO(), cmp); err != nil {
		log.Fatal(err)
	}
}

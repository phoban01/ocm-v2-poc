package main

import (
	"encoding/json"
	"fmt"
	"log"

	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/build"
	"github.com/phoban01/ocm-v2/api/v2/mutate"
	"github.com/phoban01/ocm-v2/api/v2/types"
	"github.com/phoban01/ocm-v2/providers/filesystem"
	"github.com/phoban01/ocm-v2/providers/helm"
	"github.com/phoban01/ocm-v2/providers/oci"
)

// - one resource referencing a local helm chart,
// - one referencing a helm chart stored in the same OCI registry as the cd
// - one referencing a helm chart stored in another oci registry
// - one stored in some helm chart repo?
// It would be good to see how the CD is created, uploaded, downloaded and afterwards how the content of the resources is fetched.
func main() {
	// chart from file system
	localChartMeta := types.ObjectMeta{
		Name: "istio",
		Type: types.HelmChart,
	}

	localChartAccess, err := filesystem.FromFile(
		"./istiod-1.19.3.tgz",
		filesystem.WithMediaType("application/tar+gzip"),
	)
	if err != nil {
		log.Fatal(err)
	}

	localChart := build.NewResource(localChartMeta, localChartAccess)

	// chart from helm repository
	helmRepoChart, err := helm.ChartFromRepository(
		"nginx",
		"nginx-stable/nginx-ingress",
		"0.17.1",
	)
	if err != nil {
		log.Fatal(err)
	}

	// oci registry
	ociChartSameRegistry, err := helm.ChartFromRepository(
		"podinfo",
		"oci://ghcr.io/phoban01/charts/podinfo",
		"6.5.2",
	)
	if err != nil {
		log.Fatal(err)
	}

	// remote oci registry
	ociChartRemoteRegistry, err := helm.ChartFromRepository(
		"spire",
		"oci://ghcr.io/spiffe/helm-charts/spire",
		"0.13.2",
	)
	if err != nil {
		log.Fatal(err)
	}

	resources := []v2.Resource{
		localChart,
		helmRepoChart,
		ociChartSameRegistry,
		ociChartRemoteRegistry,
	}

	// create a new component
	cmp := build.New("ocm.software/helm-test", "v1.0.0", "acme.org")

	// add the resources to the component
	cmp = mutate.WithResources(cmp, resources...)

	// setup the repository
	repo, err := oci.Repository("ghcr.io/phoban01/ocm-v2-helm")
	if err != nil {
		log.Fatal(err)
	}

	// write the component to the archive
	if err := repo.Write(cmp); err != nil {
		log.Fatal(err)
	}

	// read the component from the archive
	cmp, err = repo.Get("ocm.software/helm-test", "v1.0.0")
	if err != nil {
		log.Fatal(err)
	}

	// get the descriptor
	cd, err := cmp.Descriptor()
	if err != nil {
		log.Fatal(err)
	}

	// serialize the descriptor
	cdJSON, err := json.MarshalIndent(cd, " ", " ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(cdJSON))
}

package helm

import (
	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/provider"
)

type accessor struct {
	chart   string
	version string
	labels  map[string]string
}

var _ v2.Access = (*accessor)(nil)

var MediaType = "helm"

func init() {
	provider.Register(&accessor{})
}

func FromChart(chart, version string) (v2.Access, error) {
	return &accessor{chart: chart, version: version}, nil
}

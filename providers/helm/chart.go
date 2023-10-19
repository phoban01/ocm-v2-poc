package helm

import (
	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/provider"
)

type accessor struct {
	path       string
	chart      string
	version    string
	accessType string
	mediaType  string
	labels     map[string]string
}

var _ v2.Access = (*accessor)(nil)

var (
	LocalAccessType = "localBlob/v1"
	MediaType       = "application/tar+gzip"
	OCIAccessType   = "ociArtifact/v1"
	OCIMediaType    = "application/vnd.cncf.helm.chart.content.v1.tar+gzip"
)

func init() {
	provider.Register(&accessor{})
}

func FromRepository(chart, version string, opts ...AccessOption) (v2.Access, error) {
	a := &accessor{
		chart:   chart,
		version: version,
	}
	for _, f := range opts {
		f(a)
	}
	return a, nil
}

func FromFile(path string) (v2.Access, error) {
	return &accessor{path: path}, nil
}

package helm

import (
	"strings"

	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/build"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

func ChartFromRepository(name, ref, version string) (v2.Resource, error) {
	meta := types.ObjectMeta{
		Name:    name,
		Type:    types.HelmChart,
		Version: version,
	}

	ociChart := strings.HasPrefix(ref, "oci://")
	mediaType := MediaType
	accessType := LocalAccessType
	if ociChart {
		mediaType = OCIMediaType
		accessType = OCIAccessType
	}

	opts := []AccessOption{
		WithMediaType(mediaType),
		WithAccessType(accessType),
	}

	access, err := FromRepository(ref, version, opts...)
	if err != nil {
		return nil, err
	}

	return build.NewResource(meta, access, build.Deferrable(ociChart)), nil
}

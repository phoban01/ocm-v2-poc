package decode

import (
	"errors"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/file"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

func Resource(resource types.Resource) (v2.Resource, error) {
	switch resource.Type {
	case file.Type:
		return file.DecodeResource(resource)
	default:
		return nil, errors.New("unknown resource type")
	}
}

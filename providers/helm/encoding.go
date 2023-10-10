package helm

import (
	"encoding/json"
	"fmt"

	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

func (a *accessor) Decode(ctx v2.RepositoryContext, resource types.Resource) (v2.Access, error) {
	if err := json.Unmarshal(resource.Access, a); err != nil {
		return nil, fmt.Errorf("failed to unmarshal access: %w", err)
	}
	// a.basepath = filepath.Join(strings.TrimPrefix(ctx.Location(), "file://"), "blobs")
	return a, nil
}

func (a *accessor) MarshalJSON() ([]byte, error) {
	result := map[string]string{
		"chart":     a.chart,
		"version":   a.version,
		"type":      string(a.Type()),
		"mediaType": a.MediaType(),
	}
	return json.Marshal(result)
}

func (a *accessor) UnmarshalJSON(data []byte) error {
	obj := make(map[string]string)
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	a.chart = obj["chart"]
	a.version = obj["version"]
	a.mediaType = obj["mediaType"]
	return nil
}

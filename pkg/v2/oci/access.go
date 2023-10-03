package oci

import (
	"encoding/json"
	"io"

	"github.com/google/go-containerregistry/pkg/v1/mutate"
	v2 "github.com/phoban01/ocm-v2/pkg/v2"
)

type artifactAccess struct {
	image *image
}

func (a *artifactAccess) Type() v2.AccessType {
	return v2.AccessType("ociArtifact")
}

func (a *artifactAccess) Data() (io.ReadCloser, error) {
	return mutate.Extract(a.image.img), nil
}

func (a *artifactAccess) MarshalJSON() ([]byte, error) {
	result := map[string]string{
		"imageReference": a.image.ref,
		"type":           string(a.Type()),
	}
	return json.Marshal(result)
}

func (a *artifactAccess) UnmarshalJSON(data []byte) error {
	obj := make(map[string]string)
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	a.image = &image{
		ref: obj["imageReference"],
	}
	return nil
}

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

func (b *artifactAccess) Type() v2.AccessType {
	return v2.AccessType("ociArtifact")
}

func (b *artifactAccess) Data() (io.ReadCloser, error) {
	return mutate.Extract(b.image.img), nil
}

func (d *artifactAccess) MarshalJSON() ([]byte, error) {
	result := map[string]string{
		"imageReference": d.image.ref,
		"type":           string(d.Type()),
	}
	return json.Marshal(result)
}

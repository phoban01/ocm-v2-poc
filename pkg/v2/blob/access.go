package blob

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
)

type access struct {
	blob *blob
}

var _ v2.Access = (*access)(nil)

func (a *access) Type() v2.AccessType {
	return v2.AccessType("localBlob/v1")
}

func (a *access) Data() (io.ReadCloser, error) {
	if a.blob.data != nil {
		return io.NopCloser(bytes.NewReader(a.blob.data)), nil
	}
	return os.Open(a.blob.path)
}

func (a *access) MarshalJSON() ([]byte, error) {
	result := map[string]string{
		"localReference": a.blob.path,
		"type":           string(a.Type()),
	}
	return json.Marshal(result)
}

func (a *access) UnmarshalJSON(data []byte) error {
	obj := make(map[string]string)
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	a.blob = &blob{
		path: obj["localReference"],
	}
	return nil
}

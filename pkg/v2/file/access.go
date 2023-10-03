package file

import (
	"encoding/json"
	"io"
	"os"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
)

type access struct {
	file *file
}

var _ v2.Access = (*access)(nil)

func (a *access) Type() v2.AccessType {
	return v2.AccessType("localBlob/v1")
}

func (a *access) Data() (io.ReadCloser, error) {
	return os.Open(a.file.path)
}

func (a *access) MarshalJSON() ([]byte, error) {
	result := map[string]string{
		"localReference": a.file.path,
		"type":           string(a.Type()),
	}
	return json.Marshal(result)
}

func (a *access) UnmarshalJSON(data []byte) error {
	obj := make(map[string]string)
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	a.file = &file{
		path: obj["localReference"],
	}
	return nil
}

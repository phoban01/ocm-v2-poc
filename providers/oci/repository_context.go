package oci

import (
	"encoding/json"

	v2 "github.com/phoban01/ocm-v2/api/v2"
)

type repositoryContext struct {
	url    string
	prefix string
}

const RepositoryType = "ociRegistry"

var _ v2.RepositoryContext = (*repositoryContext)(nil)

func (c *repositoryContext) Type() string {
	return RepositoryType
}

func (c *repositoryContext) Location() string {
	return c.url
}

func (c *repositoryContext) MarshalJSON() ([]byte, error) {
	data := map[string]string{
		"type":     RepositoryType,
		"location": c.url,
	}
	return json.Marshal(data)
}

func (c *repositoryContext) UnmarshalJSON(data []byte) error {
	var response map[string]string
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}
	c.url = response["location"]
	return nil
}

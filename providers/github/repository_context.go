package github

import (
	"encoding/json"

	v2 "github.com/phoban01/ocm-v2/api/v2"
)

type repositoryContext struct {
	url string
}

const RepositoryType = "filesystem"

var _ v2.RepositoryContext = (*repositoryContext)(nil)

func (c *repositoryContext) Type() string {
	return RepositoryType
}

func (c *repositoryContext) Location() string {
	return c.url
}

func (c *repositoryContext) MarshalJSON() ([]byte, error) {
	return json.Marshal(c)
}

func (c *repositoryContext) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, c)
}

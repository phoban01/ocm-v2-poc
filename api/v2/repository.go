package v2

import (
	"context"
	"io"
)

type Repository interface {
	RepositoryStorage

	Context() RepositoryContext
	List() ([]Component, error)
	Get(name string, version string) (Component, error)
	Write(context.Context, Component) error
	Delete() error
}

type RepositoryStorage interface {
	ReadBlob(context.Context, string) (io.ReadCloser, error)
	WriteBlob(context.Context, Access) (Access, error)
}

type RepositoryContext interface {
	Location() string
	Type() string
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

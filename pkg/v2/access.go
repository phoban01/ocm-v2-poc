package v2

import "io"

type AccessType string

type Access interface {
	Type() AccessType

	Data() (io.ReadCloser, error)

	MarshalJSON() ([]byte, error)
}

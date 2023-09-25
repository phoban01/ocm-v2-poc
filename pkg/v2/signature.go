package v2

type Signature interface {
	Name() string

	Digest() (string, error)
}

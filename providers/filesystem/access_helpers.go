package filesystem

import v2 "github.com/phoban01/ocm-v2/api/v2"

type AccessOption func(*accessor)

func WithMediaType(mediaType string) func(*accessor) {
	return func(a *accessor) {
		a.mediaType = mediaType
	}
}

func FromFile(path string, opts ...AccessOption) (v2.Access, error) {
	a := &accessor{
		filepath:   path,
		repository: &repository{},
	}
	for _, f := range opts {
		f(a)
	}
	return a, nil
}

func FromBytes(data []byte, opts ...AccessOption) (v2.Access, error) {
	a := &accessor{
		data:       data,
		repository: &repository{},
	}
	for _, f := range opts {
		f(a)
	}
	return a, nil
}

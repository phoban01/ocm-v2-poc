package configr

import "errors"

var ErrNotFound = errors.New("config not found")

type Configuration interface {
	Get(string) (string, error)
}

type StaticConfig map[string]string

func (s StaticConfig) Get(k string) (string, error) {
	v, ok := s[k]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}

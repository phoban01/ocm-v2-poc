package filesystem

import (
	v2 "github.com/phoban01/ocm-v2/api/v2"
)

func (r *repository) ReadBlob(string) (v2.Access, error) {
	return nil, nil
}

func (r *repository) WriteBlob(v2.Access) error {
	return nil
}

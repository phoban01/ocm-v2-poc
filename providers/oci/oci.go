package oci

import (
	v2 "github.com/phoban01/ocm-v2/api/v2"
)

func FromImage(ref string) (v2.Access, error) {
	return &accessor{ref: ref}, nil
}

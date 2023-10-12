package query

import (
	"errors"
	"io"

	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

func GetResource(component v2.Component, meta types.ObjectMeta) (v2.Resource, error) {
	resources, err := component.Resources()
	if err != nil {
		return nil, err
	}

	for _, item := range resources {
		if item.Name() == meta.Name && item.Type() == meta.Type {
			return item, nil
		}
	}

	return nil, errors.New("resource not found")
}

func GetResourceData(component v2.Component, meta types.ObjectMeta) (io.ReadCloser, error) {
	resource, err := GetResource(component, meta)
	if err != nil {
		return nil, err
	}

	acc, err := resource.Access()
	if err != nil {
		return nil, err
	}

	return acc.Data()
}

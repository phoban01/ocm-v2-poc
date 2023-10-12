package provider

import (
	"errors"
	"fmt"

	"github.com/phoban01/ocm-v2/api/v2/types"
)

var providers = make(map[string]Provider)

func Register(p Provider) error {
	if p.Type() == "" {
		return errors.New("type not specified")
	}
	providers[p.Type()] = p
	return nil
}

func lookup(accessType string, resource types.Resource) (Provider, error) {
	p, ok := providers[accessType]
	if !ok {
		return nil, fmt.Errorf("no provider registered for id: %s", accessType)
	}
	return p, nil
}

package provider

import (
	"encoding/json"
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

func lookup(resource types.Resource) (Provider, error) {
	acc := make(map[string]any)
	if err := json.Unmarshal(resource.Access, &acc); err != nil {
		return nil, err
	}
	p, ok := providers[acc["type"].(string)]
	if !ok {
		return nil, fmt.Errorf("no provider registered for id: %s", acc["type"])
	}
	return p, nil
}

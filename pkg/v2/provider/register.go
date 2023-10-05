package provider

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

type Provider interface {
	v2.Access
	Decode(resource types.Resource) (v2.Access, error)
}

var providers = make(map[string]Provider)

func Register(p Provider) error {
	if p.Type() == "" {
		return errors.New("type not specified")
	}
	if p.MediaType() == "" {
		return errors.New("media-type not specified")
	}
	hv := UniqueIDFor(string(p.Type()), p.MediaType(), p.Labels())
	providers[hv] = p
	return nil
}

func Lookup(resource types.Resource) (v2.Access, error) {
	acc := make(map[string]any)
	if err := json.Unmarshal(resource.Access, &acc); err != nil {
		return nil, err
	}
	hv := UniqueIDFor(acc["type"].(string), acc["mediaType"].(string), nil)
	p, ok := providers[hv]
	if !ok {
		return nil, fmt.Errorf("no provider registered for id: %s", hv)
	}
	return p.Decode(resource)
}

func UniqueIDFor(artifactType, mediaType string, labels map[string]string) string {
	id := fmt.Sprintf("%s:%s", artifactType, mediaType)
	if labels != nil {
		return id + fmt.Sprintf(":%s", hashMap(labels))
	}
	return id
}

func hashMap(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		h.Write([]byte(k))
		h.Write([]byte(m[k]))
	}

	return hex.EncodeToString(h.Sum(nil))
}

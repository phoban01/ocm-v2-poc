package types

import "encoding/json"

type ResourceType string

type MediaType string

type Decodable interface {
	Resource | Source | Reference
}

type Resource struct {
	Name   string          `json:"name"`
	Type   ResourceType    `json:"type"`
	Access json.RawMessage `json:"access"`
	Digest *Digest         `json:"digest"`
}

type Source struct {
	Name   string          `json:"name"`
	Type   string          `json:"type"`
	Access json.RawMessage `json:"access"`
}

type Reference struct {
	Name   string          `json:"name"`
	Type   string          `json:"type"`
	Access json.RawMessage `json:"access"`
}

type Signature struct {
	Name          string        `json:"name"`
	Value         string        `json:"type"`
	SignatureInfo SignatureInfo `json:"signature"`
}

type Digest struct {
	HashAlgorithm          string `json:"hash_algorithm"`
	NormalisationAlgorithm string `json:"normalisation_algorithm"`
	Value                  string `json:"value"`
}

type SignatureInfo struct {
	Algorithm string `json:"algorithm"`
	MediaType string `json:"media_type"`
	Value     string `json:"value"`
	Issuer    string `json:"issuer"`
}

package types

import "encoding/json"

type MediaType string

type ObjectMeta struct {
	Name        string            `json:"name"`
	Type        ResourceType      `json:"type,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Version     string            `json:"version,omitempty"`
}

type Resource struct {
	ObjectMeta `json:",inline"`
	Access     json.RawMessage `json:"access"`
	Digest     *Digest         `json:"digest"`
}

type Source struct {
	ObjectMeta `json:",inline"`
	Access     json.RawMessage `json:"access"`
}

type Reference struct {
	ObjectMeta `json:",inline"`
	Access     json.RawMessage `json:"access"`
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

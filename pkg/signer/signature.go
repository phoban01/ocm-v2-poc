package signer

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"os"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

type sig struct {
	name       string
	resources  []v2.Resource
	privateKey []byte
}

var _ v2.Signature = (*sig)(nil)

func New(name string, privateKeyPath string, resources ...v2.Resource) (v2.Signature, error) {
	data, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	return &sig{
		name:       name,
		resources:  resources,
		privateKey: data,
	}, nil
}

func (s *sig) Name() string {
	return s.name
}

func (s *sig) Digest() (*types.Digest, error) {
	return s.sign()
}

func (s sig) Info() *types.SignatureInfo {
	return nil
}

func (s sig) sign() (*types.Digest, error) {
	block, _ := pem.Decode(s.privateKey)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	var digest string

	for _, r := range s.resources {
		d, err := r.Digest()
		if err != nil {
			return nil, err
		}
		digest += d.Value + "\n"
	}

	if digest == "" {
		return nil, errors.New("no digest provided")
	}

	hashed := sha256.Sum256([]byte(digest))

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}

	return &types.Digest{
		Value: hex.EncodeToString(signature),
	}, nil
}

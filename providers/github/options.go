package github

import (
	"github.com/phoban01/ocm-v2/api/v2/authn"
	"github.com/phoban01/ocm-v2/api/v2/configr"
)

type Option func(*repository)

func WithAuth(a authn.Authenticator) Option {
	return func(r *repository) {
		r.auth = a
	}
}

func WithKeychain(k authn.Keychain) Option {
	return func(r *repository) {
		r.keychain = k
	}
}

func WithConfig(c configr.Configuration) Option {
	return func(r *repository) {
		r.config = c
	}
}

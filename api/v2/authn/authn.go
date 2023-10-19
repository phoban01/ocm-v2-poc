package authn

import "crypto/x509"

type Authenticator interface {
	Authorization() (*AuthConfig, error)
}

// username & password statisfy basic auth
// access & refresh tokens statisfy oauth
// certificate can be used for x509 auth flows
type AuthConfig struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Auth     string `json:"auth,omitempty"`

	AccessToken string `json:"accessToken,omitempty"`

	RefreshToken string `json:"refreshToken,omitempty"`

	Certificate x509.Certificate `json:"certificate,omitempty"`
}

// Resource represents a registry or repository that can be authenticated against.
type Resource interface {
	// String returns the full string representation of the target, e.g.
	// gcr.io/my-project or just gcr.io.
	String() string

	// RegistryStr returns just the registry portion of the target, e.g. for
	// gcr.io/my-project, this should just return gcr.io. This is needed to
	// pull out an appropriate hostname.
	RegistryStr() string
}

type Keychain interface {
	// Resolve looks up the most appropriate credential for the specified target.
	Resolve(Resource) (Authenticator, error)
}

type Basic struct {
	Username string
	Password string
}

func (b *Basic) Authorization() (*AuthConfig, error) {
	return &AuthConfig{
		Username: b.Username,
		Password: b.Password,
	}, nil
}

type Bearer struct {
	Token string
}

func (b *Bearer) Authorization() (*AuthConfig, error) {
	return &AuthConfig{
		AccessToken: b.Token,
	}, nil
}

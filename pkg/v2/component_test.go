package v2_test

import (
	"testing"

	"github.com/phoban01/ocm-v2/pkg/v2/builder"
	"github.com/phoban01/ocm-v2/pkg/v2/file"
	"github.com/phoban01/ocm-v2/pkg/v2/mutate"
)

func TestComponent(t *testing.T) {
	c := builder.New("test", "v1.0.0", "acme.org")

	if c.Version() != "v1.0.0" {
		t.Errorf("Version: expected %s, but got %s", "v1.0.0", c.Version())
	}

	x := file.New("data", "./testdata/file-resource.txt")

	digest, err := x.Digest()
	if err != nil {
		t.Errorf("Digest: expected err to be nil but got: %v", err)
	}

	wantDigest := "6fc7b2fbc88b39bac5296c7382578f401abab4e554bf3781cb25a3ee0bea5e46"
	if digest != wantDigest {
		t.Errorf("Digest: expected digest to equal %s but got: %s", wantDigest, digest)
	}

	c = mutate.AddResources(c, x)

	wantLen := 1
	resources, err := c.Resources()
	if err != nil {
		t.Errorf("Add Resources: expected err to be nil but got: %v", err)
	}

	if len(resources) != wantLen {
		t.Errorf(
			"Add Resources: expected resource length to be %d but got: %d",
			wantLen,
			len(resources),
		)
	}
}

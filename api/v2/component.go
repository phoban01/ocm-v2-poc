package v2

import "github.com/phoban01/ocm-v2/api/v2/types"

type Component interface {
	// Version returns the version string of the component.
	Version() string

	// Version returns the version string of the component.
	RepositoryContext() ([]RepositoryContext, error)

	// Provider returns the Provider metadata for te Component.
	Provider() (*types.Provider, error)

	// Descriptor returns the Descriptor for the component.
	Descriptor() (*types.Descriptor, error)

	// Resources returns the list of resources that the component contains.
	Resources() ([]Resource, error)

	// Sources returns the list of sources that the component contains.
	Sources() ([]Source, error)

	// References returns the list of component references that the component contains.
	References() ([]Reference, error)

	// Signatures returns the list of component signatures that the component contains.
	Signatures() ([]Signature, error)
}

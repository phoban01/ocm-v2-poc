package v2

type Component interface {
	// Version returns the version string of the component.
	Version() string

	// Version returns the version string of the component.
	StorageContext() ([]StorageContext, error)

	// Provider returns the Provider metadata for te Component.
	Provider() (*Provider, error)

	// Descriptor returns the Descriptor for the component.
	Descriptor() (*Descriptor, error)

	// Resources returns the list of resources that the component contains.
	Resources() ([]Resource, error)

	// Sources returns the list of sources that the component contains.
	Sources() ([]Source, error)

	// References returns the list of component references that the component contains.
	References() ([]Reference, error)

	// Signatures returns the list of component signatures that the component contains.
	Signatures() ([]Signature, error)
}

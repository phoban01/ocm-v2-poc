package v2

type Repository interface {
	RepositoryStorage

	Context() RepositoryContext
	List() ([]Component, error)
	Get(name string, version string) (Component, error)
	Write(Component) error
	Delete() error
}

type RepositoryStorage interface {
	ReadBlob(string) (Access, error)
	WriteBlob(Access) error
}

type RepositoryContext interface {
	Location() string
	Type() string
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

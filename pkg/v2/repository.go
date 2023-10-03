package v2

type Repository interface {
	Context() *RepositoryContext

	List() ([]Component, error)
	Get(name string, version string) (Component, error)
	Write(Component) error
	Delete() error
}

type RepositoryContext struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

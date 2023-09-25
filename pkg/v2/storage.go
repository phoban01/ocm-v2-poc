package v2

type Storage interface {
	Write(Component) error
	Context() *StorageContext
}

type StorageContext struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

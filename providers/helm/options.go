package helm

type AccessOption func(*accessor)

func WithMediaType(mediaType string) func(*accessor) {
	return func(a *accessor) {
		a.mediaType = mediaType
	}
}

func WithAccessType(accessType string) func(*accessor) {
	return func(a *accessor) {
		a.accessType = accessType
	}
}

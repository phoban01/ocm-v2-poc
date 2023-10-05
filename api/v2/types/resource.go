package types

type ResourceType string

const (
	Blob      ResourceType = "blob"
	OCIImage  ResourceType = "ociImage"
	HelmChart ResourceType = "helmChart"
)

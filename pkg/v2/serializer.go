package v2

// func init() {
// 	c := `{
//   metadata: name: =~ ".+\\..+"
//   version: =~ "^v\\d+\\.\\d+\\.\\d+$"
// }`
//
// 	cuego.MustConstrain(&serializationFormat{}, c)
// }
//
// type serializationFormat struct {
// 	Metadata          Metadata            `json:"metadata"`
// 	Provider          Provider            `json:"provider"`
// 	Version           string              `json:"version"`
// 	RepositoryContext []RepositoryContext `json:"repository_context"`
// 	Resources         []resource          `json:"resources"`
// 	Signatures        []signature         `json:"signatures"`
// }
//
// type resource struct {
// 	Metadata     `json:",inline"`
// 	Digest       Digest `json:"digest"`
// 	MediaType    string `json:"media_type"`
// 	ResourceType string `json:"resource_type"`
// 	Access       Access `json:"access"`
// }
//
// type access struct {
// 	Type     string `json:"type"`
// 	Location string `json:"location"`
// }
//
// type signature struct {
// 	Metadata `       json:",inline"`
// 	Digest   *Digest `json:"digest"`
// }

// func serialize(d *Descriptor) (*serializationFormat, error) {
// 	sf := &serializationFormat{
// 		Metadata:          d.Metadata,
// 		RepositoryContext: d.RepositoryContext,
// 		Provider:          d.Provider,
// 		Version:           d.Version,
// 	}
//
// 	for _, r := range d.Resources {
// 		dig, err := r.Digest()
// 		if err != nil {
// 			return nil, err
// 		}
// 		rx := resource{
// 			Metadata: Metadata{
// 				Name:   r.Name(),
// 				Labels: r.Labels(),
// 			},
// 			ResourceType: string(r.ResourceType()),
// 			MediaType:    string(r.MediaType()),
// 			Digest:       dig,
// 			Access:       r.Access(),
// 		}
// 		sf.Resources = append(sf.Resources, rx)
// 	}
//
// 	for _, r := range d.Signatures {
// 		dig, err := r.Digest()
// 		if err != nil {
// 			return nil, err
// 		}
// 		sx := signature{
// 			Metadata: Metadata{
// 				Name: r.Name(),
// 			},
// 			Digest: dig,
// 		}
// 		sf.Signatures = append(sf.Signatures, sx)
// 	}
// 	return sf, nil
// }

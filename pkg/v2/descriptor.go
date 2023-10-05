package v2

import "github.com/phoban01/ocm-v2/pkg/v2/types"

// Descriptor represents the Component Descriptor in a structured way.
type Descriptor struct {
	types.ObjectMeta  `json:",inline"`
	Provider          Provider            `json:"provider"`
	RepositoryContext []RepositoryContext `json:"repository_context,omitempty"`
	Resources         []types.Resource    `json:"resources,omitempty"`
	Sources           []types.Signature   `json:"sources,omitempty"`
	References        []types.Reference   `json:"references,omitempty"`
	Signatures        []types.Signature   `json:"signatures,omitempty"`
}

type Provider struct {
	Name string `json:"name"`
}

// func (d *Descriptor) MarshalJSON() ([]byte, error) {
// so, err := serialize(d)
// if err != nil {
// 	return nil, err
// }
// if err := cuego.Validate(so); err != nil {
// 	// TODO: print all validation errors
// 	errs := errors.Errors(err)
// 	return nil, fmt.Errorf("validation error: %w", errs[len(errs)-1])
// }
// return json.Marshal(d)
// }

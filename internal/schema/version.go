package schema

import "errors"

type Version struct {
	Name string `json:"name"`
	// Constraint string `json:"constraint,omitempty"`
	Content string `json:"content,omitempty"`
}

func (v *Version) Validate() error {
	if v.Name == "" {
		return errors.New("dependency name is requried")
	}
	return nil
}

// func (v *Version) AsConstraint() string {
// 	if v.Constraint != "" {
// 		return v.Constraint
// 	}
// 	return v.Name
// }

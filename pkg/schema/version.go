package schema

import "errors"

type Version struct {
	Name string `json:"name"`
	Link string `json:"link,omitempty"`
}

func (v *Version) Validate() error {
	if v.Name == "" {
		return errors.New("dependency name is requried")
	}
	return nil
}

package schema

type Version struct {
	Name    string `json:"name"`
	Content string `json:"content,omitempty"`
}

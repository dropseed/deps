package pullrequest

import (
	"github.com/dropseed/deps/internal/config"
	"github.com/dropseed/deps/internal/schema"
)

// Pullrequest stores the basic data
type Pullrequest struct {
	Base         string
	Head         string
	Title        string
	Body         string
	Dependencies *schema.Dependencies
	Config       *config.Dependency
}

// NewPullrequest creates a Pullrequest using env variables
func NewPullrequest(base string, head string, deps *schema.Dependencies, cfg *config.Dependency) (*Pullrequest, error) {
	return &Pullrequest{
		Base:         base,
		Head:         head,
		Title:        deps.Title,
		Body:         deps.Description,
		Dependencies: deps,
		Config:       cfg,
	}, nil
}

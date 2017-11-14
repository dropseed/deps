package config

import (
	"flag"
)

// Flags store the go flags used by the user
type Flags struct {
	Branch             string
	Title              string
	Body               string
	DependenciesSchema string
	TitleFromSchema    bool
	BodyFromSchema     bool
}

// ParseFlags loads flags into Flags
func ParseFlags() *Flags {
	branch := flag.String("branch", "", "branch that pull request will be created from")
	title := flag.String("title", "", "pull request title")
	body := flag.String("body", "", "pull request body")

	dependenciesSchema := flag.String("dependencies-schema", "", "dependencies.io schema for the dependencies being acted on")
	titleFromSchema := flag.Bool("title-from-schema", false, "automatically generate the title from the dependencies-schema")
	bodyFromSchema := flag.Bool("body-from-schema", false, "automatically generate the body from the dependencies-schema")

	flag.Parse()

	return &Flags{
		Branch:             *branch,
		Title:              *title,
		Body:               *body,
		DependenciesSchema: *dependenciesSchema,
		TitleFromSchema:    *titleFromSchema,
		BodyFromSchema:     *bodyFromSchema,
	}
}

package config

import (
	"errors"
	"strings"

	"github.com/dependencies-io/pullrequest/internal/app/schema"
)

// TitleFromConfig generates a PR title based on the flags used
func (config *Config) TitleFromConfig() (string, error) {
	// user supplied title overrides all
	if config.Flags.Title != "" {
		return config.Flags.Title, nil
	}

	if config.Flags.DependenciesSchema != "" {
		schema, err := schema.NewSchemaFromString([]byte(config.Flags.DependenciesSchema))
		if err != nil {
			return "", err
		}

		if config.Flags.TitleFromSchema {
			title, err := schema.GenerateTitleFromSchema()
			if err != nil {
				return "", err
			}

			return title, nil
		}
	}

	return "", errors.New("\"title\" is required or needs the appropriate settings to be generated from the dependencies-schema")
}

// RelatedPRTitleSearchFromConfig generates the PR title search query based on the flags and env settings used
func (config *Config) RelatedPRTitleSearchFromConfig() (string, error) {
	// user supplied title overrides all
	if config.Flags.RelatedPRTitleSearch != "" {
		return config.Flags.RelatedPRTitleSearch, nil
	}

	if config.Flags.DependenciesSchema != "" {
		// generate the title
		schema, err := schema.NewSchemaFromString([]byte(config.Flags.DependenciesSchema))
		if err != nil {
			return "", err
		}

		title, err := schema.GenerateRelatedPRTitleSearchFromSchema()
		if err != nil {
			return "", err
		}

		return title, nil
	}

	return "", errors.New("Cannot use related PR behaviors without a title to search (use --dependencies-schema or --related-pr-title-search)")
}

// BodyFromConfig generates a PR body based on the flags used
func (config *Config) BodyFromConfig() (string, error) {
	maxBodyLength := 65535
	body := ""

	if config.Flags.Body != "" {

		// user supplied body overrides all
		body = config.Flags.Body

	} else if config.Flags.DependenciesSchema != "" && config.Flags.BodyFromSchema {

		schema, err := schema.NewSchemaFromString([]byte(config.Flags.DependenciesSchema))
		if err != nil {
			return "", err
		}

		body, err = schema.GenerateBodyFromSchema()
		if err != nil {
			return "", err
		}
	}

	if body != "" {
		// look for additional user content to add to the body
		if config.EnvSettings.PullrequestNotes != "" {
			body = strings.TrimSpace(config.EnvSettings.PullrequestNotes) + "\n\n---\n\n" + strings.TrimSpace(body)
		}

		// trim the pr body string to a max of this size,
		// should rarely happen but this way API call should still be success
		if len(body) > maxBodyLength {
			body = body[:maxBodyLength]
		}

		return body, nil
	}

	return "", errors.New("the \"body\" flag is required or needs the appropriate settings to be generated from the dependencies-schema")
}

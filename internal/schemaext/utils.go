package schemaext

import (
	"fmt"
	"strings"
)

func dependencyNameForDisplay(name string) string {

	prefixes := []string{
		"https://github.com/",
		"https://gitlab.com/",
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) {
			name = name[len(prefix):]
		}
	}

	suffixes := []string{
		".git/",
		".git",
	}

	for _, suffix := range suffixes {
		if strings.HasSuffix(name, suffix) {
			name = name[:len(name)-len(suffix)]
		}
	}

	return name
}

func optionalMarkdownLink(s, url string) string {
	if url == "" {
		return s
	}
	return fmt.Sprintf("[%s](%s)", s, url)
}

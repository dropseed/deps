package version

import "fmt"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var WithMeta = fmt.Sprintf("%v\ncommit %v\nbuilt at %v", version, commit, date)

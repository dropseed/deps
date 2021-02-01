package lag

import "os"

func Run() error {
	// infer files we know about
	// see if we have last install info

	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return nil
	}

	// More complicated to run a command than I thought
	// https://direnv.net/docs/hook.html
	// https://github.com/direnv/direnv/blob/master/shell_zsh.go

	print("git!")

	return nil
}

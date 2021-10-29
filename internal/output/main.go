package output

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
)

var Verbosity = 0

func isGitHubActions() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}

func shouldColorize() bool {
	isTerm := terminal.IsTerminal(int(os.Stdout.Fd()))
	if isGitHubActions() {
		// https://github.com/fatih/color#github-actions
		color.NoColor = false
	}
	return isTerm || isGitHubActions()
}

func IsDebug() bool {
	return Verbosity > 0
}

func Event(f string, args ...interface{}) {
	if shouldColorize() && IsDebug() {
		color.Set(color.FgMagenta)
		print("> ")
		color.Unset()
		color.Set(color.Bold)
	}
	fmt.Printf(f+"\n", args...)
	if shouldColorize() && IsDebug() {
		color.Unset()
	}
}

func Debug(f string, args ...interface{}) {
	if !IsDebug() {
		return
	}
	if shouldColorize() {
		color.Set(color.FgCyan)
		print("> ")
		color.Unset()
	}
	fmt.Printf(f+"\n", args...)
}

func Warning(f string, args ...interface{}) {
	color.Set(color.FgYellow)
	fmt.Printf(f+"\n", args...)
	color.Unset()
}

func Error(f string, args ...interface{}) {
	color.Set(color.FgRed)
	fmt.Printf(f+"\n", args...)
	color.Unset()
}

func Success(f string, args ...interface{}) {
	color.Set(color.FgGreen)
	fmt.Printf(f+"\n", args...)
	color.Unset()
}

func Subtle(f string, args ...interface{}) {
	color.Set(color.Faint)
	fmt.Printf(f+"\n", args...)
	color.Unset()
}

func Unstyled(f string, args ...interface{}) {
	color.Unset()
	fmt.Printf(f+"\n", args...)
}

func StartSection(f string, args ...interface{}) {
	if isGitHubActions() {
		// https://docs.github.com/en/actions/learn-github-actions/workflow-commands-for-github-actions#grouping-log-lines
		formatted := fmt.Sprintf(f, args...)
		fmt.Printf("::group::%s\n", formatted)
	} else {
		Event(f, args...)
	}
}

func EndSection() {
	if isGitHubActions() {
		fmt.Println("::endgroup::")
	}
}

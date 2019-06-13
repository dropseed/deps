package output

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
)

var Verbosity = 0

func IsDebug() bool {
	return Verbosity > 0
}

func Event(f string, args ...interface{}) {
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		color.Set(color.FgMagenta)
		print("> ")
		color.Unset()
	}
	fmt.Printf(f+"\n", args...)
}

func Debug(f string, args ...interface{}) {
	if !IsDebug() {
		return
	}
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
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

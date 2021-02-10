package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

const hookSlug = "deps"
const hookCmd = "lag"

type hookContext struct {
	// SelfPath is the unescaped absolute path to direnv
	SelfPath string
	HookCmd  string
	HookSlug string
}

const bashHook = `
_{{.HookSlug}}_hook() {
  local previous_exit_status=$?;
  trap -- '' SIGINT;
  "{{.SelfPath}}" {{.HookCmd}};
  trap - SIGINT;
  return $previous_exit_status;
};
if ! [[ "${PROMPT_COMMAND:-}" =~ _{{.HookSlug}}_hook ]]; then
  PROMPT_COMMAND="_{{.HookSlug}}_hook${PROMPT_COMMAND:+;$PROMPT_COMMAND}"
fi
`

const zshHook = `
_{{.HookSlug}}_hook() {
  trap -- '' SIGINT;
  "{{.SelfPath}}" {{.HookCmd}};
  trap - SIGINT;
}
typeset -ag precmd_functions;
if [[ -z ${precmd_functions[(r)_{{.HookSlug}}_hook]} ]]; then
  precmd_functions=( _{{.HookSlug}}_hook ${precmd_functions[@]} )
fi
`

const usage = `## BASH

Add the following line at the end of the ~/.bashrc file:

eval "$(deps shellhook bash)"

Make sure it appears even after rvm, git-prompt and other shell extensions that manipulate the prompt.

## ZSH

Add the following line at the end of the ~/.zshrc file:

eval "$(deps shellhook zsh)"
`

var shellhookCmd = &cobra.Command{
	Use:    "shellhook",
	Args:   cobra.MinimumNArgs(1),
	Hidden: true,
	Long:   usage,
	Run: func(cmd *cobra.Command, args []string) {
		shellType := args[0]

		selfPath, err := os.Executable()
		if err != nil {
			printErrAndExitFailure(err)
		}

		// Convert Windows path if needed
		selfPath = strings.Replace(selfPath, "\\", "/", -1)
		ctx := hookContext{
			SelfPath: selfPath,
			HookCmd:  hookCmd,
			HookSlug: hookSlug,
		}

		hookStr := ""

		if shellType == "bash" {
			hookStr = bashHook
		} else if shellType == "zsh" {
			hookStr = zshHook
		} else {
			printErrAndExitFailure(fmt.Errorf("unknown target shell '%s'", shellType))
		}

		hookTemplate, err := template.New("hook").Parse(hookStr)
		if err != nil {
			printErrAndExitFailure(err)
		}

		err = hookTemplate.Execute(os.Stdout, ctx)
		if err != nil {
			printErrAndExitFailure(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(shellhookCmd)
}

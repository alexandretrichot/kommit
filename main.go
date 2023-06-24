package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// RootCmd is the root command for kommit
// take a conventional commit message
// it will parse the message and create a commit message with the correct format
// allow the usage of tab completion to select the type of commit and the scope based on completion at runtime
var RootCmd = &cobra.Command{
	Use:   "kommit [type] [scope] [message]",
	Short: "A conventional commit message generator",
	Long: `A conventional commit message generator
		Example: kommit feat auth add auth feature`,
	Args: cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		c := Commit{
			Type:    args[0],
			Scope:   args[1],
			Message: strings.Join(args[2:], " "),
		}

		if err := createCommit(c); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
	ValidArgsFunction: TypeGet,
}

func TypeGet(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	switch len(args) {
	case 0:
		return getTypes(), cobra.ShellCompDirectiveNoFileComp
	case 1:
		return getScopes(), cobra.ShellCompDirectiveNoFileComp
	case 2:
		return getMessages(), cobra.ShellCompDirectiveNoFileComp
	default:
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
}

var CompletionCmd = &cobra.Command{
	Use:                   "completion [bash|zsh|fish|powershell]",
	Short:                 "Generate completion script",
	Long:                  "To load completions",
	CompletionOptions:     cobra.CompletionOptions{DisableDefaultCmd: true},
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	RootCmd.AddCommand(CompletionCmd)
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"
	"os/exec"
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
		Example: kommit feat auth "add auth feature"`,
	Args: cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		cType := args[0]
		cScope := args[1]
		// rest of the args are the message
		cMessage := strings.Join(args[2:], " ")

		fmt.Printf("%s(%s): %s\n", cType, cScope, cMessage)
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

type Commit struct {
	Type    string
	Scope   string
	Message string
}

func getMessages() []string {
	commits, err := getCommits(10)
	if err != nil {
		return nil
	}

	// get all the messages from the commits
	messages := make([]string, len(commits))
	for i, commit := range commits {
		messages[i] = commit.Message
	}

	return messages
}

func getTypes() []string {
	return []string{"feat", "refactor", "chore", "fix", "style", "perf", "test", "docs"}
}

func getScopes() []string {
	commits, err := getCommits(500)
	if err != nil {
		return nil
	}

	// get all the scopes from the commits and make unique
	scopes := make([]string, len(commits)+1)
	scopes[0] = ""
	for i, commit := range commits {
		scopes[i+1] = commit.Scope
	}

	return unique(scopes)
}

func unique(items []string) []string {
	// make a map of the items
	// the key is the item and the value is true
	// this will make all the items unique
	uniqueItems := make(map[string]bool)
	for _, item := range items {
		uniqueItems[item] = true
	}

	// make a slice of the keys
	// this will be the unique items
	uniqueSlice := make([]string, len(uniqueItems))
	i := 0
	for item := range uniqueItems {
		uniqueSlice[i] = item
		i++
	}

	return uniqueSlice
}

func getCommits(limit int) ([]Commit, error) {
	// get the git log
	cmd := exec.Command("git", "log", "--pretty=format:%s")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// split the output by new line
	lines := strings.Split(string(out), "\n")

	// make a slice of commits
	commits := []Commit{}

	// loop through the lines and parse the commit
	for _, line := range lines {
		if commit, err := parseCommit(line); err == nil {
			commits = append(commits, commit)
		}
	}

	return commits[:limit], nil
}

// return a commit from a line or nil if parsing failed
// the line should be in the format of `type(scope): message` or `type: message`
// the scope is optional
func parseCommit(line string) (Commit, error) {
	// split the line by `:`
	// the first item should be the type and scope
	// the second item should be the message
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return Commit{}, fmt.Errorf("Failed to parse commit message `%v`", line)
	}

	// split the type and scope by `(` and `)`
	// the first item should be the type
	// the second item should be the scope
	typeScope := strings.Split(parts[0], "(")
	cType := typeScope[0]
	cScope := ""
	if len(typeScope) > 1 {
		cScope = strings.TrimSuffix(typeScope[1], ")")
	}

	// trim the message
	cMessage := strings.TrimSpace(parts[1])

	return Commit{cType, cScope, cMessage}, nil
}

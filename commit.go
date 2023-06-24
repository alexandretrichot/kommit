package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Commit struct {
	Type    string
	Scope   string
	Message string
}

func (c Commit) String() string {
	if c.Scope == "" || c.Scope == "-" {
		return fmt.Sprintf("%s: %s", c.Type, c.Message)
	}

	return fmt.Sprintf("%s(%s): %s", c.Type, c.Scope, c.Message)
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

func createCommit(c Commit) error {
	// build the commit message
	cmd := exec.Command("git", "commit", "-m", c.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

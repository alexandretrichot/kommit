package main

func getTypes() []string {
	return []string{"feat", "refactor", "chore", "fix", "style", "perf", "test", "docs"}
}

func getScopes() []string {
	commits, err := getCommits(200)
	if err != nil {
		return nil
	}

	// get all the scopes from the commits and make unique
	scopes := make([]string, len(commits)+1)
	scopes[0] = "-"
	for i, commit := range commits {
		scopes[i+1] = commit.Scope
	}

	return unique(scopes)
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

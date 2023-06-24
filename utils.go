package main

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

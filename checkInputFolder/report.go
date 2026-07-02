package FolderCheckerInput

import (
	"fmt"
	"sort"
)

func SortIssues(issues []Issue) {

	sort.Slice(
		issues,
		func(i, j int) bool {

			if issues[i].Type != issues[j].Type {
				return issues[i].Type < issues[j].Type
			}

			return issues[i].Path < issues[j].Path

		},
	)
}

func Print(issues []Issue, stats interface{}) {

	SortIssues(issues)

	fmt.Println(
		"Filesystem portability check",
	)

	fmt.Println()

	for _, issue := range issues {

		fmt.Printf(
			"ERROR %-28s %s\n",
			issue.Type,
			issue.Path,
		)

		if issue.Info != "" {

			fmt.Println(
				" x  ",
				issue.Info,
			)
		}
	}

	fmt.Println()

	fmt.Printf(
		"Errors: %d\n",
		len(issues),
	)
}

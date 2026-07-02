package checkInputFolder

import (
	"fmt"
	"os"
	"sort"
)

func (chk *checker) report() {
	issueLen := len(chk.issuesDirectorySiblings)

	if issueLen > 0 {
		chk.printLocalizedListHeader(DirectorySiblingsFound, issueLen)

		for _, issue := range chk.issuesDirectorySiblings {
			fmt.Fprintf(os.Stderr, "    - "+issue.Path+"\n")
		}
	}

	//---------------------------------------------------------------------

	issueLen = len(chk.issuesInvalidWindowsCharacters)

	if issueLen > 0 {
		chk.printLocalizedListHeader(InvalidWindowsChar, issueLen)

		for _, issue := range chk.issuesInvalidWindowsCharacters {
			fmt.Fprintf(os.Stderr, "    - "+issue.Path+"\n")
		}
	}

	//---------------------------------------------------------------------

	issueLen = len(chk.issuesTrailingCharacters)

	if issueLen > 0 {
		chk.printLocalizedListHeader(TrailingDotSpace, issueLen)

		for _, issue := range chk.issuesTrailingCharacters {
			fmt.Fprintf(os.Stderr, "    - "+issue.Path+"\n")
		}
	}

	//---------------------------------------------------------------------

	issueLen = len(chk.issuesReservedWindowsName)

	if issueLen > 0 {
		chk.printLocalizedListHeader(ReservedWindowsName, issueLen)

		for _, issue := range chk.issuesReservedWindowsName {
			fmt.Fprintf(os.Stderr, "    - "+issue.Path+"\n")
		}
	}

	//---------------------------------------------------------------------

	issueLen = len(chk.issuesNameLength)

	if issueLen > 0 {
		chk.printLocalizedListHeader(NameTooLong, issueLen)

		for _, issue := range chk.issuesNameLength {
			fmt.Fprintf(os.Stderr, "    - "+issue.Path+"\n")
		}
	}

	//---------------------------------------------------------------------

	issueLen = len(chk.issuesPathLength)

	if issueLen > 0 {
		chk.printLocalizedListHeader(PathTooLong, issueLen)

		for _, issue := range chk.issuesPathLength {
			fmt.Fprintf(os.Stderr, "    - "+issue.Path+"\n")
		}
	}

	//---------------------------------------------------------------------

	issueLen = len(chk.issuesUnicodeNormalization)

	if issueLen > 0 {
		chk.printLocalizedListHeader(UnicodeCollision, issueLen)

		for _, issue := range chk.issuesUnicodeNormalization {
			fmt.Fprintf(os.Stderr, "    - "+issue.Path+"\n")
		}
	}

	//---------------------------------------------------------------------

	issueLen = len(chk.issuesSymLink)

	if issueLen > 0 {
		chk.printLocalizedListHeader(SymbolicLinkDetected, issueLen)

		for _, issue := range chk.issuesSymLink {
			fmt.Fprintf(os.Stderr, "    - "+issue.Path+"\n")
		}
	}
}

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

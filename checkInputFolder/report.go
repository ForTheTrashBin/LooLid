package checkInputFolder

import (
	"fmt"
	"os"
)

func (chk *checker) report() {
	issueLen := len(chk.issuesDirectorySiblings)

	if issueLen > 0 {
		chk.printLocalizedListHeader(CheckError_DirectorySiblingsFound, issueLen)

		for _, issue := range chk.issuesDirectorySiblings {
			fmt.Fprintf(os.Stderr, "    - "+issue.Path+"\n")
		}
	}

	//---------------------------------------------------------------------

	issueLen = len(chk.issuesInvalidWindowsChar)

	if issueLen > 0 {
		chk.printLocalizedListHeader(CheckError_InvalidWindowsChar, issueLen)

		for _, issue := range chk.issuesInvalidWindowsChar {
			fmt.Fprintf(os.Stderr, "    - "+issue.Path+"\n")
		}
	}

	//---------------------------------------------------------------------

	issueLen = len(chk.issuesTrailingDotSpace)

	if issueLen > 0 {
		chk.printLocalizedListHeader(CheckError_TrailingDotSpace, issueLen)

		for _, issue := range chk.issuesTrailingDotSpace {
			fmt.Fprintf(os.Stderr, "    - "+issue.Path+"\n")
		}
	}

	//---------------------------------------------------------------------

	issueLen = len(chk.issuesReservedWindowsName)

	if issueLen > 0 {
		chk.printLocalizedListHeader(CheckError_ReservedWindowsName, issueLen)

		for _, issue := range chk.issuesReservedWindowsName {
			fmt.Fprintf(os.Stderr, "    - "+issue.Path+"\n")
		}
	}

	//---------------------------------------------------------------------

	issueLen = len(chk.issuesFileNameLength)

	if issueLen > 0 {
		chk.printLocalizedListHeader(CheckError_FileNameTooLong, issueLen)

		for _, issue := range chk.issuesFileNameLength {
			fmt.Fprintf(os.Stderr, "    - "+issue.Path+"\n")
		}
	}

	//---------------------------------------------------------------------

	issueLen = len(chk.issuesPathNameLength)

	if issueLen > 0 {
		chk.printLocalizedListHeader(CheckError_PathNameTooLong, issueLen)

		for _, issue := range chk.issuesPathNameLength {
			fmt.Fprintf(os.Stderr, "    - "+issue.Path+"\n")
		}
	}

	//---------------------------------------------------------------------

	issueLen = len(chk.issuesUnicodeNormalization)

	if issueLen > 0 {
		chk.printLocalizedListHeader(CheckError_UnicodeCollision, issueLen)

		for _, issue := range chk.issuesUnicodeNormalization {
			fmt.Fprintf(os.Stderr, "    - "+issue.Path+"\n")
		}
	}

	//---------------------------------------------------------------------

	issueLen = len(chk.issuesSymLink)

	if issueLen > 0 {
		chk.printLocalizedListHeader(CheckError_SymbolicLinkDetected, issueLen)

		for _, issue := range chk.issuesSymLink {
			fmt.Fprintf(os.Stderr, "    - "+issue.Path+"\n")
		}
	}
}

/*
func SortIssues(issues []Issue) {

	sort.Slice(
		issues,
		func(i, j int) bool {

			if issues[i].Typex != issues[j].Typex {
				return issues[i].Typex < issues[j].Typex
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
			issue.Typex,
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
*/

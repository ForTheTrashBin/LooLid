package checkInputFolder

import (
	"fmt"
	"os"
	"sort"

	"github.com/ForTheTrashBin/LooLid/helper"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

//-----------------------------------------------------------------------------
// Print formatted list of 'issues' with header
//-----------------------------------------------------------------------------

func (chk *checker) printIssues(errorType ErrorType, issues []Issue) {

	issueLen := len(issues)

	if issueLen > 0 {

		helper.PrintLocalizedListHeader(chk.localizer, string(errorType), issueLen)

		//---------------------------------------------------------------------
		// Sorting for a better customer-experience
		//---------------------------------------------------------------------

		sort.Slice(issues, func(i, j int) bool {
			return issues[i].path < issues[j].path
		})

		//---------------------------------------------------------------------

		strFolder := chk.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: string(CheckError_DuplicateEntries_Folder)})
		strFile := chk.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: string(CheckError_DuplicateEntries_File)})

		lenFolder := len(strFolder)
		lenFile := len(strFile)

		var needed int

		for _, issue := range issues {
			if issue.isDir {
				needed = max(needed, lenFolder)
			} else {
				needed = max(needed, lenFile)
			}
		}

		lineFormat := fmt.Sprintf("  - %%-%ds", needed+2)

		for _, issue := range issues {
			if issue.isDir {
				fmt.Fprintf(os.Stderr, lineFormat, strFolder+":")
			} else {
				fmt.Fprintf(os.Stderr, lineFormat, strFile+":")
			}

			fmt.Fprintf(os.Stderr, "%s\n", issue.path)
		}
	}
}

//-----------------------------------------------------------------------------
// Print formatted list of 'duplicates' with header
//-----------------------------------------------------------------------------

func (chk *checker) printDuplicates() {

	numDuplicateGroups := len(chk.duplicateGroups)

	if numDuplicateGroups > 0 {

		helper.PrintLocalizedListHeader(chk.localizer, string(CheckError_DuplicateEntries), numDuplicateGroups)

		//---------------------------------------------------------------------
		// Sorting for a better customer-experience
		//---------------------------------------------------------------------

		sort.Slice(chk.duplicateGroups, func(i, j int) bool {
			return chk.duplicateGroups[i].groupName < chk.duplicateGroups[j].groupName
		})

		for looperGroups := 0; looperGroups < len(chk.duplicateGroups); looperGroups++ {
			duplicateGroup := &chk.duplicateGroups[looperGroups]

			sort.Slice(duplicateGroup.items, func(i, j int) bool {
				return duplicateGroup.items[i].itemName < duplicateGroup.items[j].itemName
			})
		}

		//---------------------------------------------------------------------

		strFolder := chk.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: string(CheckError_DuplicateEntries_Folder)})
		strFile := chk.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: string(CheckError_DuplicateEntries_File)})

		lenFolder := len(strFolder)
		lenFile := len(strFile)

		//---------------------------------------------------------------------

		for _, duplicateGroup := range chk.duplicateGroups {

			fmt.Fprintf(os.Stderr, "  - %s: %s\n", strFolder, duplicateGroup.groupName)

			var needed int

			for _, items := range duplicateGroup.items {
				if items.isDir {
					needed = max(needed, lenFolder)
				} else {
					needed = max(needed, lenFile)
				}
			}

			lineheaderFormat := fmt.Sprintf("    - %%-%ds", needed+2)

			for _, items := range duplicateGroup.items {

				if items.isDir {
					fmt.Fprintf(os.Stderr, lineheaderFormat, strFolder+":")
				} else {
					fmt.Fprintf(os.Stderr, lineheaderFormat, strFile+":")
				}

				fmt.Fprintf(os.Stderr, "%s\n", items.itemName)
			}
		}
	}
}

//-----------------------------------------------------------------------------
// Print all issues
//-----------------------------------------------------------------------------

func (chk *checker) reportInputFolderErrors() {

	chk.printIssues(CheckError_DirectorySiblingsFound, chk.issuesInputDirectorySiblings)
}

func (chk *checker) reportBulkDataErrors() {

	chk.printIssues(CheckError_InvalidWindowsChar, chk.issuesInvalidWindowsChar)
	chk.printIssues(CheckError_TrailingDotSpace, chk.issuesTrailingDotSpace)
	chk.printIssues(CheckError_ReservedWindowsName, chk.issuesReservedWindowsName)
	chk.printIssues(CheckError_FileNameTooLong, chk.issuesFileNameLength)
	chk.printIssues(CheckError_PathNameTooLong, chk.issuesPathNameLength)
	chk.printIssues(CheckError_UnicodeCollision, chk.issuesUnicodeNormalization)
	chk.printIssues(CheckError_SymbolicLinkDetected, chk.issuesSymLink)
	chk.printIssues(CheckError_ConfigFile, chk.issuesConfigFile)

	chk.printDuplicates()
}

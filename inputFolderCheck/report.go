package inputFolderCheck

import (
	"fmt"
	"os"
	"slices"
	"sort"

	"github.com/ForTheTrashBin/LooLid/helper/constants"
	"github.com/ForTheTrashBin/LooLid/helper/nutsandbolts"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

//-----------------------------------------------------------------------------
// Print formatted list of 'issues' with header
//-----------------------------------------------------------------------------

func (chk *checker) printIssues(MessageID string, issues []Issue) {

	issueLen := len(issues)

	if issueLen > 0 {

		nutsandbolts.PrintLocalizedListHeader(chk.localizer, MessageID, issueLen)

		//---------------------------------------------------------------------
		// Sorting for a better customer-experience
		//---------------------------------------------------------------------

		sort.Slice(issues, func(i, j int) bool {
			return issues[i].filePath < issues[j].filePath
		})

		//---------------------------------------------------------------------

		strFolder := chk.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: constants.Label_Folder})
		strFile := chk.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: constants.Label_File})

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

			fmt.Fprintf(os.Stderr, "%s\n", issue.filePath)
		}
	}
}

//-----------------------------------------------------------------------------
// Print formatted list of 'duplicates' with header
//-----------------------------------------------------------------------------

func (chk *checker) printDuplicates() {

	numDuplicateGroups := len(chk.duplicateGroups)

	if numDuplicateGroups > 0 {

		nutsandbolts.PrintLocalizedListHeader(chk.localizer, constants.CheckError_DuplicateEntries, numDuplicateGroups)

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

		strFolder := chk.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: constants.Label_Folder})
		strFile := chk.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: constants.Label_File})

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

	chk.printIssues(constants.CheckError_DirectorySiblingsFound, chk.issuesInputDirectorySiblings)
}

func (chk *checker) reportTemplatesFolderErrors() {

	strFile := chk.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: constants.Label_File})
	strFolder := chk.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: constants.Label_Folder})
	strExtension := chk.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: constants.Label_Extension})

	//-------------------------------------------------------------------------

	numEntries := len(chk.violations.Directories)

	if numEntries > 0 {

		nutsandbolts.PrintLocalizedListHeader(chk.localizer, constants.CheckError_SubfoldersFound, numEntries)

		//---------------------------------------------------------------------
		// Sorting for a better customer-experience
		//---------------------------------------------------------------------

		slices.Sort(chk.violations.Directories)

		//---------------------------------------------------------------------

		for _, dirName := range chk.violations.Directories {

			fmt.Fprintf(os.Stderr, "  - %s: %s\n", strFolder, dirName)
		}
	}

	//-------------------------------------------------------------------------

	numEntries = len(chk.violations.NoExtension)

	if numEntries > 0 {

		nutsandbolts.PrintLocalizedListHeader(chk.localizer, constants.CheckError_NoExtensions, numEntries) // TODO

		//---------------------------------------------------------------------
		// Sorting for a better customer-experience
		//---------------------------------------------------------------------

		slices.Sort(chk.violations.NoExtension)

		//---------------------------------------------------------------------

		for _, fileName := range chk.violations.NoExtension {

			fmt.Fprintf(os.Stderr, "  - %s: %s\n", strFile, fileName)
		}
	}

	//-------------------------------------------------------------------------

	numEntries = len(chk.violations.DupExtensions)

	if numEntries > 0 {

		nutsandbolts.PrintLocalizedListHeader(chk.localizer, constants.CheckError_DuplicateExtension, numEntries) // TODO

		//---------------------------------------------------------------------
		// Sorting for a better customer-experience
		//---------------------------------------------------------------------

		keySlice := make([]string, 0, numEntries)

		for key := range chk.violations.DupExtensions {

			keySlice = append(keySlice, key)
		}

		slices.Sort(keySlice)

		//---------------------------------------------------------------------

		for _, key := range keySlice {

			fmt.Fprintf(os.Stderr, "  - %s: %s\n", strExtension, key)

			fileNames := chk.violations.DupExtensions[key]

			slices.Sort(fileNames)

			//-----------------------------------------------------------------

			for _, fileName := range chk.violations.DupExtensions[key] {

				fmt.Fprintf(os.Stderr, "    - %s: %s\n", strFile, fileName)
			}

		}
	}
}

func (chk *checker) reportBulkDataErrors() {

	chk.printIssues(constants.CheckError_InvalidWindowsChar, chk.issuesInvalidWindowsChar)
	chk.printIssues(constants.CheckError_TrailingDotSpace, chk.issuesTrailingDotSpace)
	chk.printIssues(constants.CheckError_ReservedWindowsName, chk.issuesReservedWindowsName)
	chk.printIssues(constants.CheckError_FileNameTooLong, chk.issuesFileNameLength)
	chk.printIssues(constants.CheckError_PathNameTooLong, chk.issuesPathNameLength)
	chk.printIssues(constants.CheckError_UnicodeCollision, chk.issuesUnicodeNormalization)
	chk.printIssues(constants.CheckError_SymbolicLinkDetected, chk.issuesSymLink)
	chk.printIssues(constants.CheckError_ConfigFile, chk.issuesConfigFile)

	chk.printDuplicates()
}

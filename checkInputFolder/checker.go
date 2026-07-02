package checkInputFolder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type checker struct {
	localizer                      *i18n.Localizer
	pureAppName                    string
	issuesTrailingCharacters       []Issue
	issuesInvalidWindowsCharacters []Issue
	issuesReservedWindowsName      []Issue
	issuesNameLength               []Issue
	issuesPathLength               []Issue
	issuesUnicodeNormalization     []Issue
	issuesSymLink                  []Issue
	issuesDirectorySiblings        []Issue
}

func newChecker(localizer *i18n.Localizer, pureAppName string) *checker {
	chk := &checker{
		localizer:   localizer,
		pureAppName: pureAppName,
	}

	return chk
}

func (chk *checker) checkIssues() bool {
	var result = true

	if len(chk.issuesDirectorySiblings) > 0 {
		result = false
	}

	if len(chk.issuesTrailingCharacters) > 0 {
		result = false
	}

	if len(chk.issuesInvalidWindowsCharacters) > 0 {
		result = false
	}

	if len(chk.issuesReservedWindowsName) > 0 {
		result = false
	}

	if len(chk.issuesNameLength) > 0 {
		result = false
	}

	if len(chk.issuesPathLength) > 0 {
		result = false
	}

	if len(chk.issuesUnicodeNormalization) > 0 {
		result = false
	}

	if len(chk.issuesSymLink) > 0 {
		result = false
	}

	return result
}

func (chk *checker) printLocalizedListHeader(errorType ErrorType, listLength int) {
	localizedMessage := chk.localizer.MustLocalize(
		&i18n.LocalizeConfig{
			MessageID: string(errorType),
			TemplateData: map[string]int{
				"Count": listLength,
			},
			PluralCount: listLength,
		})

	fmt.Fprintf(os.Stderr, localizedMessage)
}

func (chk *checker) printLocalizedMessage(errorType ErrorType, value1 string, value2 string) bool {
	localizedMessage := chk.localizer.MustLocalize(
		&i18n.LocalizeConfig{
			MessageID: string(errorType),
			TemplateData: map[string]string{
				"value1": value1,
				"value2": value2,
			},
		})

	fmt.Fprintf(os.Stderr, localizedMessage)

	return true
}

func (chk *checker) check(inputFolder string) bool {

	//-------------------------------------------------------------------------
	// get fileinfo of given input-path
	//-------------------------------------------------------------------------

	fileInfo, err := os.Stat(inputFolder)

	if err != nil {
		return chk.printLocalizedMessage(NoFileInfo, inputFolder, err.Error())
	}

	//-------------------------------------------------------------------------
	// Check if input-path points to a directory
	//-------------------------------------------------------------------------

	if !fileInfo.IsDir() {
		return chk.printLocalizedMessage(NoDirectory, inputFolder, "")
	}

	//-------------------------------------------------------------------------
	// Check, if directory-name starts with a dot (.)
	//-------------------------------------------------------------------------

	baseName := filepath.Base(inputFolder)

	// Do not treat the "." and ".." directories as hidden directories.

	if (baseName != ".") && (baseName != "..") {
		if strings.HasPrefix(baseName, ".") {
			return chk.printLocalizedMessage(DotHiddenDirectory, inputFolder, "")
		}
	}

	//-------------------------------------------------------------------------
	// On windows directory must not be hidden or system
	//-------------------------------------------------------------------------

	if isHiddenOrSystemOnWindows(inputFolder) {
		return chk.printLocalizedMessage(DirectoryHiddenOrSystem, inputFolder, "")
	}

	//-------------------------------------------------------------------------
	// Check, if directory is readable (Just try and check the result)
	//-------------------------------------------------------------------------

	_, err = os.ReadDir(inputFolder)

	if err != nil {
		if os.IsPermission(err) {
			chk.printLocalizedMessage(DirectoryNoPermission, inputFolder, "")
		} else {
			chk.printLocalizedMessage(DirectoryNotReadable, inputFolder, err.Error())
		}

		return true
	}

	//-------------------------------------------------------------------------
	// Make input-path absolute and evaluate sym-links
	//-------------------------------------------------------------------------

	inputPathAbs, err := filepath.Abs(inputFolder)

	if err != nil {
		return chk.printLocalizedMessage(NoPathAbs, inputFolder, err.Error())
	}

	inputPathAbs, err = filepath.EvalSymlinks(inputPathAbs)

	if err != nil {
		return chk.printLocalizedMessage(NoSymLinks, inputPathAbs, err.Error())
	}

	//-------------------------------------------------------------------------
	// Get working-directory (which is an absolute path) and evaluate sym-links
	//-------------------------------------------------------------------------

	workingDirAbs, err := os.Getwd()

	if err != nil {
		return chk.printLocalizedMessage(NoWorkingDir, "", err.Error())
	}

	workingDirAbs, err = filepath.EvalSymlinks(workingDirAbs)

	if err != nil {
		return chk.printLocalizedMessage(NoSymLinks, workingDirAbs, err.Error())
	}

	//-------------------------------------------------------------------------
	// 'Compare' and evaluate if workingDirAbs is inside inputPathAbs
	//-------------------------------------------------------------------------

	relativePath, err := filepath.Rel(inputPathAbs, workingDirAbs)

	if err != nil {
		return chk.printLocalizedMessage(NoPathRel, inputPathAbs, err.Error())
	}

	if (relativePath != "..") && ((len(relativePath) < 3) || (relativePath[:3] != ".."+string(filepath.Separator))) {
		return chk.printLocalizedMessage(WorkingDirInInput, inputPathAbs, err.Error())
	}

	//-------------------------------------------------------------------------
	// Try to read siblings to detect directories and files with ambiguous names
	//-------------------------------------------------------------------------

	parentPath := filepath.Dir(inputPathAbs)

	dirEntries, err := os.ReadDir(parentPath)

	if err != nil {
		return chk.printLocalizedMessage(DirectoryNotReadable, parentPath, err.Error())
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			if dirEntry.Name() != baseName { // it's NOT me!
				if strings.ToLower(dirEntry.Name()) == strings.ToLower(baseName) {
					chk.issuesDirectorySiblings = append(
						chk.issuesDirectorySiblings,
						Issue{
							Type: DirectorySiblingsFound,
							Path: dirEntry.Name(),
							Info: inputFolder,
						},
					)
				}
			}
		} else {
			if strings.Contains(strings.ToLower(dirEntry.Name()), strings.ToLower(baseName)) {
				chk.issuesDirectorySiblings = append(
					chk.issuesDirectorySiblings,
					Issue{
						Type: DirectorySiblingsFound,
						Path: dirEntry.Name(),
						Info: inputFolder,
					},
				)
			}
		}
	}

	entries, err := scanDirectory(inputFolder)

	if err != nil {
		return false
	}

	chk.validateEntries(entries)

	var issues []Issue

	issues = append(
		issues,
		analyzeCollisions(entries)...,
	)

	return chk.checkIssues()
}

func Check(localizer *i18n.Localizer, pureAppName string, inputFolder string) bool {
	checker := newChecker(
		localizer,
		pureAppName,
	)

	return checker.check(inputFolder)
}

func CheckAndReport(localizer *i18n.Localizer, pureAppName string, inputfolder string) bool {
	checker := newChecker(
		localizer,
		pureAppName,
	)

	result := checker.check(inputfolder)

	if !result {
		checker.report()
	}

	return result
}

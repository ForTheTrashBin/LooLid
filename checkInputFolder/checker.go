package checkInputFolder

import (
	"LooLid/helper"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type checker struct {
	localizer   *i18n.Localizer
	inputFolder string
	pureAppName string

	issuesInputDirectorySiblings []Issue // checker.go
	issuesInvalidWindowsChar     []Issue // rules.go
	issuesTrailingDotSpace       []Issue // rules.go
	issuesReservedWindowsName    []Issue // rules.go
	issuesFileNameLength         []Issue // rules.go
	issuesPathNameLength         []Issue // rules.go
	issuesUnicodeNormalization   []Issue // rules.go
	issuesSymLink                []Issue // rules.go
	issuesConfigFile             []Issue // scanner.go

	duplicateGroups []DuplicateGroup
}

func newChecker(localizer *i18n.Localizer, inputFolder string, pureAppName string) *checker {
	chk := &checker{
		localizer:   localizer,
		inputFolder: inputFolder,
		pureAppName: pureAppName,
	}

	return chk
}

func (chk *checker) isInputDirectoryCorrect() bool {
	result := true

	if len(chk.issuesInputDirectorySiblings) > 0 {
		result = false
	}

	return result
}

func (chk *checker) isAllTheBulkDataCorrect() bool {

	result := true

	if len(chk.issuesInvalidWindowsChar) > 0 {
		result = false
	}

	if len(chk.issuesTrailingDotSpace) > 0 {
		result = false
	}

	if len(chk.issuesReservedWindowsName) > 0 {
		result = false
	}

	if len(chk.issuesFileNameLength) > 0 {
		result = false
	}

	if len(chk.issuesPathNameLength) > 0 {
		result = false
	}

	if len(chk.issuesUnicodeNormalization) > 0 {
		result = false
	}

	if len(chk.issuesSymLink) > 0 {
		result = false
	}

	if len(chk.issuesConfigFile) > 0 {
		return false
	}

	if len(chk.duplicateGroups) > 0 {
		result = false
	}

	return result
}

func (chk *checker) printLocalizedMessage(errorType ErrorType, value1 string, value2 string) bool {
	helper.PrintLocalizedMessage(chk.localizer, string(errorType), value1, value2)

	return false
}

func (chk *checker) checkInputFolder() bool {

	//-------------------------------------------------------------------------
	// get fileinfo of given input-path
	//-------------------------------------------------------------------------

	fileInfo, err := os.Stat(chk.inputFolder)

	if err != nil {

		return chk.printLocalizedMessage(CheckError_NoFileInfo, chk.inputFolder, err.Error())
	}

	//-------------------------------------------------------------------------
	// Check if input-path points to a directory
	//-------------------------------------------------------------------------

	if !fileInfo.IsDir() {

		return chk.printLocalizedMessage(CheckError_NoDirectory, chk.inputFolder, "")
	}

	//-------------------------------------------------------------------------
	// Check, if directory-name starts with a dot (.)
	//-------------------------------------------------------------------------

	baseName := filepath.Base(chk.inputFolder)

	// Do not treat the "." and ".." directories as hidden directories.

	if (baseName != ".") && (baseName != "..") {

		if strings.HasPrefix(baseName, ".") {

			return chk.printLocalizedMessage(CheckError_DotHiddenDirectory, chk.inputFolder, "")
		}
	}

	//-------------------------------------------------------------------------
	// On windows directory must not be hidden or system
	//-------------------------------------------------------------------------

	if isHiddenOrSystemOnWindows(chk.inputFolder) {

		return chk.printLocalizedMessage(CheckError_DirectoryHiddenOrSystem, chk.inputFolder, "")
	}

	//-------------------------------------------------------------------------
	// Check whether the directory name could be confused with the configuration file
	//-------------------------------------------------------------------------

	if strings.ToLower(baseName) == strings.ToLower(helper.GetConfigFileName()) {

		return chk.printLocalizedMessage(CheckError_Confusion, chk.inputFolder, "")
	}

	//-------------------------------------------------------------------------
	// Check, if directory is readable (Just try and check the result)
	//-------------------------------------------------------------------------

	_, err = os.ReadDir(chk.inputFolder)

	if err != nil {
		if os.IsPermission(err) {

			return chk.printLocalizedMessage(CheckError_DirectoryNoPermission, chk.inputFolder, "")
		} else {

			return chk.printLocalizedMessage(CheckError_DirectoryNotReadable, chk.inputFolder, err.Error())
		}
	}

	//-------------------------------------------------------------------------
	// Make input-path absolute and evaluate sym-links
	//-------------------------------------------------------------------------

	inputPathAbs, err := filepath.Abs(chk.inputFolder)

	if err != nil {

		return chk.printLocalizedMessage(CheckError_NoPathAbs, chk.inputFolder, err.Error())
	}

	inputPathAbs, err = filepath.EvalSymlinks(inputPathAbs)

	if err != nil {

		return chk.printLocalizedMessage(CheckError_NoSymLinks, inputPathAbs, err.Error())
	}

	//-------------------------------------------------------------------------
	// Get working-directory (which is an absolute path) and evaluate sym-links
	//-------------------------------------------------------------------------

	workingDirAbs, err := os.Getwd()

	if err != nil {

		return chk.printLocalizedMessage(CheckError_NoWorkingDir, "", err.Error())
	}

	workingDirAbs, err = filepath.EvalSymlinks(workingDirAbs)

	if err != nil {

		return chk.printLocalizedMessage(CheckError_NoSymLinks, workingDirAbs, err.Error())
	}

	//-------------------------------------------------------------------------
	// 'Compare' and evaluate if workingDirAbs is inside inputPathAbs
	//-------------------------------------------------------------------------

	if true {
		relativePath, err := filepath.Rel(inputPathAbs, workingDirAbs)

		if err != nil {

			return chk.printLocalizedMessage(CheckError_NoPathRel, inputPathAbs, err.Error())
		}

		if (relativePath != "..") && ((len(relativePath) < 3) || (relativePath[:3] != ".."+string(filepath.Separator))) {

			return chk.printLocalizedMessage(CheckError_WorkingDirInInput, inputPathAbs, "")
		}
	}

	//-------------------------------------------------------------------------
	// Try to read siblings to detect directories and files with ambiguous names
	//-------------------------------------------------------------------------

	parentPath := filepath.Dir(inputPathAbs)

	dirEntries, err := os.ReadDir(parentPath)

	if err != nil {

		return chk.printLocalizedMessage(CheckError_DirectoryNotReadable, parentPath, err.Error())
	}

	for _, dirEntry := range dirEntries {

		if dirEntry.IsDir() {

			if dirEntry.Name() != baseName { // it's NOT me!

				processedEntryName, _ := strings.CutPrefix(strings.ToLower(dirEntry.Name()), ".")
				processedBaseName, _ := strings.CutPrefix(strings.ToLower(baseName), ".")

				if processedEntryName == processedBaseName {

					chk.issuesInputDirectorySiblings = append(
						chk.issuesInputDirectorySiblings,
						Issue{
							isDir: dirEntry.IsDir(),
							path:  dirEntry.Name(),
							info:  chk.inputFolder,
						},
					)
				}
			}
		}
	}

	return chk.isInputDirectoryCorrect()
}

func (chk *checker) checkBulkData() bool {
	err := chk.scanDirectories()

	if err != nil {
		return false
	}

	return chk.isAllTheBulkDataCorrect()
}

func CheckInputFolder(localizer *i18n.Localizer, inputFolder string, pureAppName string) bool {

	result := false

	checker := newChecker(localizer, inputFolder, pureAppName)

	if result = checker.checkInputFolder(); result {

		result = checker.checkBulkData()
	}

	return result
}

func CheckInputFolderAndReport(localizer *i18n.Localizer, inputfolder string, pureAppName string) bool {

	result := false

	checker := newChecker(localizer, inputfolder, pureAppName)

	if result = checker.checkInputFolder(); result {

		if result = checker.checkBulkData(); !result {

			checker.reportBulkDataErrors()
		}
	} else {

		checker.reportInputFolderErrors()
	}

	return result
}

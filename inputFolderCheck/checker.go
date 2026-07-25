package inputFolderCheck

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ForTheTrashBin/LooLid/blackwhite"
	"github.com/ForTheTrashBin/LooLid/helper/constants"
	"github.com/ForTheTrashBin/LooLid/helper/nutsandbolts"
	"github.com/ForTheTrashBin/LooLid/helper/osspecific"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/afero"
	"github.com/theckman/yacspin"
)

type checker struct {
	localizer   *i18n.Localizer
	inputFolder string

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

func newChecker(localizer *i18n.Localizer, inputFolder string) *checker {

	chk := &checker{

		localizer:   localizer,
		inputFolder: inputFolder,
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

func (chk *checker) isBulkDataCorrect() bool {

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

func (chk *checker) getLocalizedMessage(MessageId string, value1 string, value2 string) string {

	return chk.localizer.MustLocalize(
		&i18n.LocalizeConfig{
			MessageID: MessageId,
			TemplateData: map[string]string{
				"value1": value1,
				"value2": value2,
			},
		})
}

func (chk *checker) checkInputFolder() error {

	//-------------------------------------------------------------------------
	// get fileinfo of given input-path
	//-------------------------------------------------------------------------

	fileInfo, err := os.Stat(chk.inputFolder)

	if err != nil {

		return errors.New(chk.getLocalizedMessage(constants.CheckError_NoFileInfo, chk.inputFolder, err.Error()))
	}

	//-------------------------------------------------------------------------
	// Check if input-path points to a directory
	//-------------------------------------------------------------------------

	if !fileInfo.IsDir() {

		return errors.New(chk.getLocalizedMessage(constants.CheckError_NoDirectory, chk.inputFolder, ""))
	}

	//-------------------------------------------------------------------------
	// Check, if directory-name starts with a dot (.)
	//-------------------------------------------------------------------------

	baseName := filepath.Base(chk.inputFolder)

	// Do not treat the "." and ".." directories as hidden directories.

	if (baseName != ".") && (baseName != "..") {

		if strings.HasPrefix(baseName, ".") {

			return errors.New(chk.getLocalizedMessage(constants.CheckError_DotHiddenDirectory, chk.inputFolder, ""))
		}
	}

	//-------------------------------------------------------------------------
	// On windows the directory must not be hidden or system
	//-------------------------------------------------------------------------

	if osspecific.IsHiddenOrSystem(chk.inputFolder) {

		return errors.New(chk.getLocalizedMessage(constants.CheckError_DirectoryHiddenOrSystem, chk.inputFolder, ""))
	}

	//-------------------------------------------------------------------------
	// Check whether the directory name could be confused with the configuration file
	//-------------------------------------------------------------------------

	if strings.EqualFold(baseName, nutsandbolts.GetConfigFileName()) {

		return errors.New(chk.getLocalizedMessage(constants.CheckError_Confusion, chk.inputFolder, ""))
	}

	//-------------------------------------------------------------------------
	// Check, if directory is readable (Just try and check the result)
	//-------------------------------------------------------------------------

	_, err = os.ReadDir(chk.inputFolder)

	if err != nil {

		if os.IsPermission(err) {

			return errors.New(chk.getLocalizedMessage(constants.CheckError_DirectoryNoPermission, chk.inputFolder, ""))
		} else {

			return errors.New(chk.getLocalizedMessage(constants.CheckError_DirectoryNotReadable, chk.inputFolder, err.Error()))
		}
	}

	//-------------------------------------------------------------------------
	// Make input-path absolute and evaluate sym-links
	//-------------------------------------------------------------------------

	inputPathAbs, err := filepath.Abs(chk.inputFolder)

	if err != nil {

		return errors.New(chk.getLocalizedMessage(constants.CheckError_NoPathAbs, chk.inputFolder, err.Error()))
	}

	inputPathAbs, err = filepath.EvalSymlinks(inputPathAbs)

	if err != nil {

		return errors.New(chk.getLocalizedMessage(constants.CheckError_NoSymLinks, inputPathAbs, err.Error()))
	}

	//-------------------------------------------------------------------------
	// Get working-directory (which is an absolute path) and evaluate sym-links
	//-------------------------------------------------------------------------

	workingDirAbs, err := os.Getwd()

	if err != nil {

		return errors.New(chk.getLocalizedMessage(constants.CheckError_NoWorkingDir, "", err.Error()))
	}

	workingDirAbs, err = filepath.EvalSymlinks(workingDirAbs)

	if err != nil {

		return errors.New(chk.getLocalizedMessage(constants.CheckError_NoSymLinks, workingDirAbs, err.Error()))
	}

	//-------------------------------------------------------------------------
	// 'Compare' and evaluate if workingDirAbs is inside inputPathAbs
	//-------------------------------------------------------------------------

	if true {

		relativePath, err := filepath.Rel(inputPathAbs, workingDirAbs)

		if err != nil {

			return errors.New(chk.getLocalizedMessage(constants.CheckError_NoPathRel, inputPathAbs, err.Error()))
		}

		if (relativePath != "..") && ((len(relativePath) < 3) || (relativePath[:3] != ".."+string(filepath.Separator))) {

			return errors.New(chk.getLocalizedMessage(constants.CheckError_WorkingDirInInput, inputPathAbs, ""))
		}
	}

	//-------------------------------------------------------------------------
	// Try to read siblings to detect directories and files with ambiguous names
	//-------------------------------------------------------------------------

	parentPath := filepath.Dir(inputPathAbs)

	dirEntries, err := os.ReadDir(parentPath)

	if err != nil {

		return errors.New(chk.getLocalizedMessage(constants.CheckError_DirectoryNotReadable, parentPath, err.Error()))
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
							isDir:    dirEntry.IsDir(),
							filePath: dirEntry.Name(),
							info:     chk.inputFolder,
						},
					)

					err = constants.ErrInputFolderNotCorrect
				}
			}
		}
	}

	return err
}

func (chk *checker) checkBulkData() error {

	diskFs := afero.NewOsFs()

	var bwConfig blackwhite.Config

	if err := chk.scanDir(diskFs, chk.inputFolder, bwConfig); err != nil {

		return err
	}

	if !chk.isBulkDataCorrect() {

		return constants.ErrBulkDataNotCorrect
	}

	return nil
}

//-----------------------------------------------------------------------------
// Check the input-directory and it's content AND print error messages
//-----------------------------------------------------------------------------

func InputFolderCheck(localizer *i18n.Localizer, inputfolder string) error {

	var err error

	checker := newChecker(localizer, inputfolder)

	if err = checker.checkInputFolder(); err == nil {

		if err = checker.checkBulkData(); err == nil {

			checker.reportBulkDataErrors()
		}
	} else {

		checker.reportInputFolderErrors()
	}

	return err
}

//-----------------------------------------------------------------------------
// Check the input-directory and it's content AND print error messages asynchronous
//-----------------------------------------------------------------------------

func InputFolderCheckAsync(localizer *i18n.Localizer, inputfolder string, sigCh chan os.Signal) error {

	checker := newChecker(localizer, inputfolder)

	spinnerSuffix := checker.getLocalizedMessage(constants.SpinnerSuffixInputfolderCheck, "", "")
	spinnerStopMessage := checker.getLocalizedMessage(constants.SpinnerStopMessageDone, "", "")

	spinnerConfig := yacspin.Config{
		Frequency:         constants.Spinner_FrequencyMS * time.Millisecond,
		CharSet:           yacspin.CharSets[constants.Spinner_CharSet],
		Suffix:            " " + spinnerSuffix,
		SuffixAutoColon:   true,
		StopCharacter:     constants.Spinner_StopCharacter,
		StopColors:        []string{constants.Spinner_StopColor},
		StopFailCharacter: constants.Spinner_StopFailCharacter,
		StopFailColors:    []string{constants.Spinner_StopFailColor},
		StopMessage:       spinnerStopMessage,
	}

	spinner, err := yacspin.New(spinnerConfig)

	if err != nil {
		panic(fmt.Errorf("spinner init failed: %w", err))
	}

	spinner.Reverse()

	if err := spinner.Start(); err != nil {
		panic(fmt.Errorf("spinner start failed: %w", err))
	}

	defer spinner.Stop()

	//-------------------------------------------------------------------------

	doneChannel := make(chan error, 1)
	panicChannel := make(chan error, 1)

	go func() {

		defer func() {

			if rec := recover(); rec != nil {

				switch value := rec.(type) {

				case error:

					panicChannel <- value

				case string:

					panicChannel <- errors.New(value)

				default:

					panicChannel <- errors.New("Recovered panic without type")
				}
			}

			panicChannel <- errors.New("Recovered panic without type")
		}()

		if err := checker.checkInputFolder(); err == nil {

			doneChannel <- checker.checkBulkData()
		} else {

			doneChannel <- err
		}
	}()

	//-------------------------------------------------------------------------

	select {

	case <-sigCh:

		stopFailMessage := checker.getLocalizedMessage(constants.CheckError_AbortedByUser, "", "")

		spinner.StopFailMessage(stopFailMessage)

		spinner.StopFail()

		return constants.ErrInterrupted

	case err := <-doneChannel:

		if err != nil {

			if err == constants.ErrInputFolderNotCorrect {

				stopFailMessage := checker.getLocalizedMessage(constants.CheckError_InputfolderIncorrect, "", "")

				spinner.StopFailMessage(stopFailMessage)

				spinner.StopFail()

				checker.reportInputFolderErrors()

				return err
			}

			if err == constants.ErrBulkDataNotCorrect {

				stopFailMessage := checker.getLocalizedMessage(constants.CheckError_BulkdataIncorrect, "", "")

				spinner.StopFailMessage(stopFailMessage)

				spinner.StopFail()

				checker.reportBulkDataErrors()

				return err
			}

			spinner.StopFailMessage(err.Error())

			spinner.StopFail()

			return err
		}

		return err

	case err := <-panicChannel:

		spinner.StopFailMessage(checker.getLocalizedMessage(constants.SpinnerStopMessageError, "", ""))

		spinner.StopFail()

		panic(err) // panics in main
	}
}

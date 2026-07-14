package inputFolderWrite

import (
	"fmt"
	"os"
	"time"

	"github.com/ForTheTrashBin/LooLid/helper/constants"
	"github.com/ForTheTrashBin/LooLid/helper/osspecific"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/afero"
	"github.com/theckman/yacspin"
)

type writer struct {
	localizer    *i18n.Localizer
	outputFolder string
	memFs        *afero.Fs
}

func newWriter(localizer *i18n.Localizer, outputFolder string, memFs *afero.Fs, x int) *writer {

	wrt := &writer{

		localizer:    localizer,
		outputFolder: outputFolder,
		memFs:        memFs,
	}

	return wrt
}

func (wrt *writer) getLocalizedMessage(MessageID string, value1 string, value2 string) string {

	return wrt.localizer.MustLocalize(
		&i18n.LocalizeConfig{
			MessageID: MessageID,
			TemplateData: map[string]string{
				"value1": value1,
				"value2": value2,
			},
		})
}

func xxxstat(fs afero.Fs, path string) (os.FileInfo, error) {

	if fs, ok := fs.(afero.Lstater); ok {

		fi, _, err := fs.LstatIfPossible(path)

		return fi, err
	}

	return fs.Stat(path)
}

func XXXPreserveTimes(sourceInfo os.FileInfo, destFs afero.Fs, dest string) error {

	timeSpec := osspecific.GetTimeSpec(sourceInfo)

	return destFs.Chtimes(dest, timeSpec.TimeAccess, timeSpec.TimeModify)
}

func xxxchmod(fs afero.Fs, dir string, mode os.FileMode, reported *error) {

	if err := fs.Chmod(dir, mode); *reported == nil {

		*reported = err
	}
}

func xxxcloseFile(f afero.File, reported *error) {

	if err := f.Close(); *reported == nil {

		*reported = err
	}
}

func (wrt *writer) writeInputFolder() error {

	var doDebug bool = false

	if doDebug {

		fmt.Println("*********************************************************************************")
		fmt.Println("*** writeInputFolder")
		fmt.Println("*********************************************************************************")
	}

	destFs := afero.NewOsFs()

	err := cleanDir(destFs, wrt.outputFolder, *wrt.memFs)

	if err != nil {

		return err
	}

	if doDebug {

		fmt.Println("*** cleanDir Ok!")
		fmt.Println("---------------------------------------------------------------------------------")
	}

	//-------------------------------------------------------------------------

	err = syncDir(*wrt.memFs, "", destFs, wrt.outputFolder, time.Now(), 0)

	if err != nil {

		panic(err)
	}

	if doDebug {

		fmt.Println("*** syncDir OK!")
		fmt.Println("*********************************************************************************")
	}

	return nil
}

func InputFolderWrite(localizer *i18n.Localizer, outputfolder string, memFs *afero.Fs) error {

	return nil
}

func InputFolderWriteAsync(localizer *i18n.Localizer, outputfolder string, memFs *afero.Fs, sigCh chan os.Signal) error {

	writer := newWriter(localizer, outputfolder, memFs, 1)

	spinnerSuffix := writer.getLocalizedMessage(constants.SpinnerSuffixInputFolderWrite, "", "")
	spinnerStopMessage := writer.getLocalizedMessage(constants.SpinnerStopMessage, "", "")

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

	go func() {

		doneChannel <- writer.writeInputFolder()
	}()

	//-------------------------------------------------------------------------

	select {

	case <-sigCh:

		stopFailMessage := writer.getLocalizedMessage(constants.CheckError_AbortedByUser, "", "")

		spinner.StopFailMessage(stopFailMessage)

		spinner.StopFail()

		return constants.ErrInterrupted

	case err := <-doneChannel:

		if err != nil {

			if err == constants.ErrInputFolderNotCorrect {

				stopFailMessage := writer.getLocalizedMessage(constants.CheckError_InputfolderIncorrect, "", "")

				spinner.StopFailMessage(stopFailMessage)

				spinner.StopFail()

				return err
			}

			if err == constants.ErrBulkDataNotCorrect {

				stopFailMessage := writer.getLocalizedMessage(constants.CheckError_BulkdataIncorrect, "", "")

				spinner.StopFailMessage(stopFailMessage)

				spinner.StopFail()

				return err
			}

			spinner.StopFailMessage(err.Error())

			spinner.StopFail()

			return err
		}

		return err
	}
}

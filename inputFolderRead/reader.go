package inputFolderRead

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ForTheTrashBin/LooLid/blackwhite"
	"github.com/ForTheTrashBin/LooLid/helper/constants"
	"github.com/ForTheTrashBin/LooLid/helper/osspecific"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/afero"
	"github.com/theckman/yacspin"
)

type reader struct {
	localizer   *i18n.Localizer
	inputFolder string
	memFs       *afero.Fs
}

func newReader(localizer *i18n.Localizer, inputFolder string, memFs *afero.Fs) *reader {

	rdr := &reader{

		localizer:   localizer,
		inputFolder: inputFolder,
		memFs:       memFs,
	}

	return rdr
}

func (rdr *reader) getLocalizedMessage(MessageID string, value1 string, value2 string) string {

	return rdr.localizer.MustLocalize(
		&i18n.LocalizeConfig{
			MessageID: MessageID,
			TemplateData: map[string]string{
				"value1": value1,
				"value2": value2,
			},
		})
}

func stat(fs afero.Fs, path string) (os.FileInfo, error) {

	if fs, ok := fs.(afero.Lstater); ok {

		fi, _, err := fs.LstatIfPossible(path)

		return fi, err
	}

	return fs.Stat(path)
}

func PreserveTimes(sourceInfo os.FileInfo, destFs afero.Fs, dest string) error {

	timeSpec := osspecific.GetTimeSpec(sourceInfo)

	return destFs.Chtimes(dest, timeSpec.TimeAccess, timeSpec.TimeModify)
}

func chmod(fs afero.Fs, dir string, mode os.FileMode, reported *error) {

	if err := fs.Chmod(dir, mode); *reported == nil {

		*reported = err
	}
}

func closeFile(f afero.File, reported *error) {

	if err := f.Close(); *reported == nil {

		*reported = err
	}
}

func (rdr *reader) readInputFolder() error {

	diskFs := afero.NewOsFs()

	var bwConfig blackwhite.Config

	return copyDir(rdr.localizer, diskFs, rdr.inputFolder, *rdr.memFs, "", bwConfig)
}

func InputFolderRead(localizer *i18n.Localizer, inputfolder string, memFs *afero.MemMapFs) error {

	return nil // TODO:
}

func InputFolderReadAsync(localizer *i18n.Localizer, inputfolder string, memFs *afero.Fs, sigCh chan os.Signal) error {

	reader := newReader(localizer, inputfolder, memFs)

	spinnerSuffix := reader.getLocalizedMessage(constants.SpinnerSuffixInputFolderRead, "", "")
	spinnerStopMessage := reader.getLocalizedMessage(constants.SpinnerStopMessageDone, "", "")

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

		doneChannel <- reader.readInputFolder()
	}()

	//-------------------------------------------------------------------------

	select {

	case <-sigCh:

		stopFailMessage := reader.getLocalizedMessage(constants.CheckError_AbortedByUser, "", "")

		spinner.StopFailMessage(stopFailMessage)

		spinner.StopFail()

		return constants.ErrInterrupted

	case err := <-doneChannel:

		if err != nil {

			if err == constants.ErrInputFolderNotCorrect {

				stopFailMessage := reader.getLocalizedMessage(constants.CheckError_InputfolderIncorrect, "", "")

				spinner.StopFailMessage(stopFailMessage)

				spinner.StopFail()

				return err
			}

			if err == constants.ErrBulkDataNotCorrect {

				stopFailMessage := reader.getLocalizedMessage(constants.CheckError_BulkdataIncorrect, "", "")

				spinner.StopFailMessage(stopFailMessage)

				spinner.StopFail()

				return err
			}

			spinner.StopFailMessage(err.Error())

			spinner.StopFail()

			return err
		}

		return err

	case err := <-panicChannel:

		spinner.StopFailMessage(reader.getLocalizedMessage(constants.SpinnerStopMessageError, "", ""))

		spinner.StopFail()

		panic(err) // panics in main
	}
}

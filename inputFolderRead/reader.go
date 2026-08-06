package inputFolderRead

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ForTheTrashBin/LooLid/configParser"
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

func (rdr *reader) readInputFolder(ctx context.Context) error {

	select {

	case <-ctx.Done():

		return ctx.Err()

	default:

		diskFs := afero.NewOsFs()

		rulesConfig := configParser.NewRulesConfig() // Start with default config

		return copyDir(ctx, rdr.localizer, diskFs, rdr.inputFolder, *rdr.memFs, "", rulesConfig, 0)
	}
}

func InputFolderRead(ctx context.Context, localizer *i18n.Localizer, inputfolder string, memFs *afero.Fs) error {

	return newReader(localizer, inputfolder, memFs).readInputFolder(ctx)
}

func InputFolderReadAsync(ctx context.Context, localizer *i18n.Localizer, inputfolder string, memFs *afero.Fs) error {

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

	if err = spinner.Start(); err != nil {

		panic(fmt.Errorf("spinner start failed: %w", err))
	}

	defer spinner.Stop() // Don't forget to stop the spinner when done

	//-------------------------------------------------------------------------

	if err = reader.readInputFolder(ctx); err != nil {

		if errors.Is(err, context.Canceled) {

			spinner.StopFailMessage(reader.getLocalizedMessage(constants.CheckError_AbortedByUser, "", ""))
		} else {

			spinner.StopFailMessage(err.Error())
		}

		spinner.StopFail()
	}

	return err
}

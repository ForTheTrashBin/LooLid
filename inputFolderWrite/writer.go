package inputFolderWrite

import (
	"context"
	"errors"
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

func newWriter(localizer *i18n.Localizer, outputFolder string, memFs *afero.Fs) *writer {

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

func (wrt *writer) writeInputFolder(ctx context.Context) error {

	destFs := afero.NewOsFs()

	err := cleanDir(ctx, destFs, wrt.outputFolder, *wrt.memFs)

	if err == nil {

		err = syncDir(ctx, *wrt.memFs, "", destFs, wrt.outputFolder, time.Now(), 0)
	}

	return err
}

func InputFolderWrite(ctx context.Context, localizer *i18n.Localizer, outputfolder string, memFs *afero.Fs) error {

	return newWriter(localizer, outputfolder, memFs).writeInputFolder(ctx)
}

func InputFolderWriteAsync(ctx context.Context, localizer *i18n.Localizer, outputfolder string, memFs *afero.Fs) error {

	writer := newWriter(localizer, outputfolder, memFs)

	//-------------------------------------------------------------------------
	// Disable echo for the duration of the processing
	//-------------------------------------------------------------------------

	if err := osspecific.SetEcho(false); err == nil {

		// Remember: defered functions are executed in LIFO order, so FlushStdin will be called before SetEcho(true)

		defer osspecific.SetEcho(true) // Restore echo on exit
		defer osspecific.FlushStdin()  // Flush stdin on exit
	}

	//-------------------------------------------------------------------------

	spinnerSuffix := writer.getLocalizedMessage(constants.SpinnerSuffixInputFolderWrite, "", "")
	spinnerStopMessage := writer.getLocalizedMessage(constants.SpinnerStopMessageDone, "", "")

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

	if err = writer.writeInputFolder(ctx); err != nil {

		if errors.Is(err, context.Canceled) {

			spinner.StopFailMessage(writer.getLocalizedMessage(constants.CheckError_AbortedByUser, "", ""))
		} else {

			spinner.StopFailMessage(err.Error())
		}

		spinner.StopFail()
	}

	return err
}

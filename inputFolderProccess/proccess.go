package inputFolderProccess

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ForTheTrashBin/LooLid/helper/constants"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/afero"
	"github.com/theckman/yacspin"
)

type proccess struct {
	localizer *i18n.Localizer
	memFs     *afero.Fs
}

func newProccess(localizer *i18n.Localizer, memFs *afero.Fs) *proccess {

	prc := &proccess{

		localizer: localizer,
		memFs:     memFs,
	}

	return prc
}

func (prc *proccess) getLocalizedMessage(MessageID string, value1 string, value2 string) string {

	return prc.localizer.MustLocalize(
		&i18n.LocalizeConfig{
			MessageID: MessageID,
			TemplateData: map[string]string{
				"value1": value1,
				"value2": value2,
			},
		})
}

func (prc *proccess) proccessInputFolder() error {

	time.Sleep(250 * time.Millisecond)

	return nil
}

func InputFolderProccess(localizer *i18n.Localizer, memFs *afero.MemMapFs) error {

	return nil // TODO:
}

func InputFolderProccessAsync(localizer *i18n.Localizer, memFs *afero.Fs, sigCh chan os.Signal) error {

	proccess := newProccess(localizer, memFs)

	spinnerSuffix := proccess.getLocalizedMessage(constants.SpinnerSuffixInputFolderProccess, "", "")
	spinnerStopMessage := proccess.getLocalizedMessage(constants.SpinnerStopMessageDone, "", "")

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

					panicChannel <- constants.ErrRecoveredPanicWithoutType
				}
			} else {

				panicChannel <- constants.ErrRecoveredPanicWithoutType
			}
		}()

		doneChannel <- proccess.proccessInputFolder()
	}()

	//-------------------------------------------------------------------------

	select {

	case <-sigCh:

		stopFailMessage := proccess.getLocalizedMessage(constants.CheckError_AbortedByUser, "", "")

		spinner.StopFailMessage(stopFailMessage)

		spinner.StopFail()

		return constants.ErrInterrupted

	case err := <-doneChannel:

		if err != nil {

			spinner.StopFailMessage(err.Error())

			spinner.StopFail()
		}

		return err

	case err := <-panicChannel:

		spinner.StopFailMessage(proccess.getLocalizedMessage(constants.SpinnerStopMessageError, "", ""))

		spinner.StopFail()

		panic(err) // panics in main
	}
}

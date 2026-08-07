package inputFolderProccess

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ForTheTrashBin/LooLid/helper/constants"
	"github.com/ForTheTrashBin/LooLid/helper/osspecific"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/afero"
	"github.com/theckman/yacspin"
)

type proccessor struct {
	localizer *i18n.Localizer
	memFs     *afero.Fs
}

func newProccessor(localizer *i18n.Localizer, memFs *afero.Fs) *proccessor {

	prc := &proccessor{

		localizer: localizer,
		memFs:     memFs,
	}

	return prc
}

func (prc *proccessor) getLocalizedMessage(MessageID string, value1 string, value2 string) string {

	return prc.localizer.MustLocalize(
		&i18n.LocalizeConfig{
			MessageID: MessageID,
			TemplateData: map[string]string{
				"value1": value1,
				"value2": value2,
			},
		})
}

func (prc *proccessor) proccessInputFolder(ctx context.Context) error {

	for range 40 {

		select {

		case <-ctx.Done():

			return ctx.Err()

		default:
		}

		time.Sleep(125 * time.Millisecond)
	}

	return nil
}

func InputFolderProccess(ctx context.Context, localizer *i18n.Localizer, memFs *afero.Fs) error {

	return newProccessor(localizer, memFs).proccessInputFolder(ctx)
}

func InputFolderProccessAsync(ctx context.Context, localizer *i18n.Localizer, memFs *afero.Fs) error {

	proccessor := newProccessor(localizer, memFs)

	//-------------------------------------------------------------------------
	// Disable echo for the duration of the processing
	//-------------------------------------------------------------------------

	if err := osspecific.SetEcho(false); err == nil {

		// Remember: defered functions are executed in LIFO order, so FlushStdin will be called before SetEcho(true)

		defer osspecific.SetEcho(true) // Restore echo on exit
		defer osspecific.FlushStdin()  // Flush stdin on exit
	}

	//-------------------------------------------------------------------------

	spinnerSuffix := proccessor.getLocalizedMessage(constants.SpinnerSuffixInputFolderProccess, "", "")
	spinnerStopMessage := proccessor.getLocalizedMessage(constants.SpinnerStopMessageDone, "", "")

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

	if err = proccessor.proccessInputFolder(ctx); err != nil {

		if errors.Is(err, context.Canceled) {

			spinner.StopFailMessage(proccessor.getLocalizedMessage(constants.CheckError_AbortedByUser, "", ""))
		} else {

			spinner.StopFailMessage(err.Error())
		}

		spinner.StopFail()
	}

	return err
}

package commands

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ForTheTrashBin/LooLid/checkInputFolder"
	"github.com/ForTheTrashBin/LooLid/helper"
	"github.com/ForTheTrashBin/LooLid/helper/constants"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/theckman/yacspin"
)

//-----------------------------------------------------------------------------
// 'locals' for printing a usage-message
//-----------------------------------------------------------------------------

var buildUsageMessage string

func printBuildUsageMessage() {
	fmt.Fprintf(os.Stderr, buildUsageMessage)
}

//-----------------------------------------------------------------------------

type BuildCommand struct {
	localizer   *i18n.Localizer
	pureAppName string
}

func (cmd *BuildCommand) GetName() string {
	return "build"
}

func (cmd *BuildCommand) Parse(localizer *i18n.Localizer, pureAppName string, args []string) error {

	cmd.localizer = localizer
	cmd.pureAppName = pureAppName

	buildUsageMessage = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    "flag." + cmd.GetName() + "UsageMessage",
		TemplateData: map[string]string{"pureAppName": pureAppName},
	})

	flagSet := flag.NewFlagSet(cmd.GetName(), flag.ContinueOnError)

	flagSet.Usage = printBuildUsageMessage

	flagSet.SetOutput(io.Discard) // prevent output from flag-library (no i18n)

	if err := flagSet.Parse(args); err != nil {
		return err
	}

	if flagSet.NArg() > 0 {
		printBuildUsageMessage()

		return helper.ErrInvalidNumberOfCommands
	}

	return nil
}

func runTask(name string, duration time.Duration, sigCh chan os.Signal) error {

	spinnerConfig := yacspin.Config{
		Frequency:         constants.Spinner_FrequencyMS * time.Millisecond,
		CharSet:           yacspin.CharSets[constants.Spinner_CharSet],
		Suffix:            " " + name,
		SuffixAutoColon:   true,
		StopCharacter:     constants.Spinner_StopCharacter,
		StopColors:        []string{constants.Spinner_StopColor},
		StopFailCharacter: constants.Spinner_StopFailCharacter,
		StopFailColors:    []string{constants.Spinner_StopFailColor},
		StopMessage:       "done",
	}

	spinner, err := yacspin.New(spinnerConfig)

	if err != nil {
		return fmt.Errorf("spinner init failed: %w", err)
	}

	spinner.Reverse()

	if err := spinner.Start(); err != nil {
		return fmt.Errorf("spinner start failed: %w", err)
	}

	defer spinner.Stop()

	//-------------------------------------------------------------------------

	doneCh := make(chan error, 1)

	go func() {

		timer := time.NewTimer(duration)

		defer timer.Stop()

		<-timer.C

		// Beispiel: künstlicher Fehler in Step 2
		if name == "Dummy Step 4: proccessing acknowledge" {
			doneCh <- errors.New("simulated processing error")
			return
		}

		doneCh <- nil
	}()

	select {

	case <-sigCh:

		spinner.StopFailMessage("aborted by user")

		spinner.StopFail()

		return constants.ErrInterrupted

	case err := <-doneCh:

		if err != nil {

			spinner.StopFailMessage(err.Error())

			spinner.StopFail()

			return err
		}

		return nil
	}
}

func (cmd *BuildCommand) Execute() error {

	// inputFolder := "/home/u32800"
	inputFolder := constants.AppConfig_DefaultInputDirectory
	// inputFolder := "/home/u32800/Dokumente/Development/Websites/Website_de/hugo/content"

	//-------------------------------------------------------------------------

	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	defer signal.Stop(sigCh)

	//-------------------------------------------------------------------------

	if err := checkInputFolder.CheckInputFolderAsync(cmd.localizer, inputFolder, cmd.pureAppName, sigCh); err != nil {
		return nil
	}

	if err := runTask("Dummy Step 2: processing data", 1*time.Second, sigCh); err != nil {
		return nil
	}

	if err := runTask("Dummy Step 3: sending data", 1*time.Second, sigCh); err != nil {
		return nil
	}

	if err := runTask("Dummy Step 4: proccessing acknowledge", 1*time.Second, sigCh); err != nil {
		return nil
	}

	if err := runTask("Dummy Step 5: saving results", 1*time.Second, sigCh); err != nil {
		return nil
	}

	return nil
}

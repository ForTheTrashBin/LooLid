package commands

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/ForTheTrashBin/LooLid/helper/constants"
	"github.com/ForTheTrashBin/LooLid/helper/nutsandbolts"
	"github.com/ForTheTrashBin/LooLid/inputFolderCheck"
	"github.com/ForTheTrashBin/LooLid/inputFolderProccess"
	"github.com/ForTheTrashBin/LooLid/inputFolderRead"
	"github.com/ForTheTrashBin/LooLid/inputFolderWrite"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/afero"
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
	localizer *i18n.Localizer
}

func (cmd *BuildCommand) GetName() string {
	return "build"
}

func (cmd *BuildCommand) Parse(localizer *i18n.Localizer, pureAppName string, args []string) error {

	cmd.localizer = localizer

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

		return nutsandbolts.ErrInvalidNumberOfCommands
	}

	return nil
}

func (cmd *BuildCommand) Execute() error {

	inputFolder := constants.AppConfig_DefaultInputFolder
	outputFolder := constants.AppConfig_DefaultOutputFolder

	//-------------------------------------------------------------------------

	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	defer signal.Stop(sigCh)

	//-------------------------------------------------------------------------

	if err := inputFolderCheck.InputFolderCheckAsync(cmd.localizer, inputFolder, sigCh); err != nil {

		return nil
	}

	memFs := afero.NewMemMapFs()

	if err := inputFolderRead.InputFolderReadAsync(cmd.localizer, inputFolder, &memFs, sigCh); err != nil {

		return nil
	}

	if err := inputFolderProccess.InputFolderProccessAsync(cmd.localizer, &memFs, sigCh); err != nil {

		return nil
	}

	if err := inputFolderWrite.InputFolderWriteAsync(cmd.localizer, outputFolder, &memFs, sigCh); err != nil {

		return nil
	}

	return nil
}

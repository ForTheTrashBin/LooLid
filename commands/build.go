package commands

import (
	"context"
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

	fmt.Fprint(os.Stderr, buildUsageMessage)
}

//-----------------------------------------------------------------------------

type BuildCommand struct {
	localizer *i18n.Localizer
}

func (cmd *BuildCommand) GetName() string {
	return "build"
}

func (cmd *BuildCommand) Parse(localizer *i18n.Localizer, args []string) error {

	cmd.localizer = localizer

	buildUsageMessage = localizer.MustLocalize(&i18n.LocalizeConfig{

		MessageID: "flag." + cmd.GetName() + "UsageMessage",

		TemplateData: map[string]string{

			"pureAppName":    nutsandbolts.GetPureAppName(),
			"configFileName": nutsandbolts.GetConfigFileName(),
		},
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

	//-------------------------------------------------------------------------
	// Create a context that is canceled (ctx.Done()) when an interrupt signal is received
	//-------------------------------------------------------------------------

	ctx, stopSignaling := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	defer stopSignaling() // Don't forget to stop the signal notification when done

	//-------------------------------------------------------------------------

	inputFolder := constants.AppConfig_DefaultInputFolder
	outputFolder := constants.AppConfig_DefaultOutputFolder
	templatesFolder := constants.AppConfig_DefaultTemplatesFolder

	//-------------------------------------------------------------------------

	if err := inputFolderCheck.InputFolderCheckAsync(ctx, cmd.localizer, inputFolder, templatesFolder); err != nil {

		return nil
	}

	if ctx.Err() == nil {

		memFs := afero.NewMemMapFs()

		if err := inputFolderRead.InputFolderReadAsync(ctx, cmd.localizer, inputFolder, &memFs); err != nil {

			return nil
		}

		if ctx.Err() == nil {

			if err := inputFolderProccess.InputFolderProccessAsync(ctx, cmd.localizer, &memFs); err != nil {

				return nil
			}
		}

		if ctx.Err() == nil {

			if err := inputFolderWrite.InputFolderWriteAsync(ctx, cmd.localizer, outputFolder, &memFs); err != nil {

				return nil
			}
		}
	}

	return nil
}

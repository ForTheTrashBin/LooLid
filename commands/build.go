package commands

import (
	"LooLid/checkInputFolder"
	"LooLid/helper"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/nicksnyder/go-i18n/v2/i18n"
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

func (cmd *BuildCommand) Execute() error {

	if checkInputFolder.CheckAndReport(cmd.localizer, cmd.pureAppName, "checkInputFolder/test/testInputFolder") {
		fmt.Println("********** build is OK ***************")
	} else {
		fmt.Println("********** build is abborted ***************")
	}

	return nil
}

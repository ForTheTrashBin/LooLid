package Commands

import (
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

	buildUsageMessage, _ = localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "flag." + cmd.GetName() + "UsageMessage",
		TemplateData: map[string]string{"pureAppName": pureAppName},
	})

	flagSet := flag.NewFlagSet(cmd.GetName(), flag.ContinueOnError)

	flagSet.Usage = printBuildUsageMessage

	flagSet.SetOutput(io.Discard) // prevent output from flag-library (no i18n)

	if err := flagSet.Parse(args); err != nil {
		return err
	}

	return nil
}

func (cmd *BuildCommand) Execute() error {
	fmt.Println("************** Execute build **************")
	/*
		issues, stats, err := FolderCheckerInput.Check("~")

		if err != nil {
			fmt.Println("Fatal:", err)
			os.Exit(2)
		}

		FolderCheckerInput.Print(issues, stats)

		if len(issues) > 0 {
			os.Exit(1)
		}
	*/
	return nil
}

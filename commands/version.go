package commands

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

//-----------------------------------------------------------------------------
// 'locals' for printing a usage-message
//-----------------------------------------------------------------------------

var versionUsageMessage string

func printVersionUsageMessage() {
	fmt.Fprintf(os.Stderr, versionUsageMessage)
}

type VersionCommand struct {
	localizer   *i18n.Localizer
	pureAppName string
}

func (cmd *VersionCommand) GetName() string {
	return "version"
}

func (cmd *VersionCommand) Parse(localizer *i18n.Localizer, pureAppName string, args []string) error {

	cmd.localizer = localizer
	cmd.pureAppName = pureAppName

	versionUsageMessage, _ = localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "flag." + cmd.GetName() + "UsageMessage",
		TemplateData: map[string]string{"pureAppName": pureAppName},
	})

	flagSet := flag.NewFlagSet(cmd.GetName(), flag.ContinueOnError)

	flagSet.Usage = printVersionUsageMessage

	flagSet.SetOutput(io.Discard) // prevent output from flag-library (no i18n)

	if err := flagSet.Parse(args); err != nil {
		return err
	}

	if flagSet.NArg() > 0 {
		printVersionUsageMessage()

		return errors.New("invalid number of commands")
	}

	return nil
}

func (cmd *VersionCommand) Execute() error {
	fmt.Println(cmd.pureAppName, "v0.1.0")

	return nil
}

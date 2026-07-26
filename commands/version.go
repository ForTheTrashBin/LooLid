package commands

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ForTheTrashBin/LooLid/helper/nutsandbolts"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

//-----------------------------------------------------------------------------
// 'locals' for printing a usage-message
//-----------------------------------------------------------------------------

var versionUsageMessage string

func printVersionUsageMessage() {

	fmt.Fprint(os.Stderr, versionUsageMessage)
}

type VersionCommand struct {
	localizer *i18n.Localizer
}

func (cmd *VersionCommand) GetName() string {
	return "version"
}

func (cmd *VersionCommand) Parse(localizer *i18n.Localizer, args []string) error {

	cmd.localizer = localizer

	versionUsageMessage = localizer.MustLocalize(&i18n.LocalizeConfig{

		MessageID: "flag." + cmd.GetName() + "UsageMessage",

		TemplateData: map[string]string{

			"pureAppName":    nutsandbolts.GetPureAppName(),
			"configFileName": nutsandbolts.GetConfigFileName(),
		},
	})

	flagSet := flag.NewFlagSet(cmd.GetName(), flag.ContinueOnError)

	flagSet.Usage = printVersionUsageMessage

	flagSet.SetOutput(io.Discard) // prevent output from flag-library (no i18n)

	if err := flagSet.Parse(args); err != nil {
		return err
	}

	if flagSet.NArg() > 0 {
		printVersionUsageMessage()

		return nutsandbolts.ErrInvalidNumberOfCommands
	}

	return nil
}

func (cmd *VersionCommand) Execute() error {
	fmt.Println(nutsandbolts.GetPureAppName(), "v0.1.0")

	return nil
}

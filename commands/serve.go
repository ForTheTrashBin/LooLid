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

var serveUsageMessage string

func printServeUsageMessage() {
	fmt.Fprintf(os.Stderr, serveUsageMessage)
}

//-----------------------------------------------------------------------------

type ServeCommand struct {
	localizer   *i18n.Localizer
	pureAppName string
	port        int
	doUpdate    bool
}

func (cmd *ServeCommand) GetName() string {
	return "serve"
}

func (cmd *ServeCommand) Parse(localizer *i18n.Localizer, pureAppName string, args []string) error {

	cmd.localizer = localizer
	cmd.pureAppName = pureAppName

	serveUsageMessage = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    "flag." + cmd.GetName() + "UsageMessage",
		TemplateData: map[string]string{"pureAppName": pureAppName},
	})

	flagSet := flag.NewFlagSet(cmd.GetName(), flag.ContinueOnError)

	flagSet.Usage = printServeUsageMessage

	flagSet.SetOutput(io.Discard) // prevent output from flag-library (no i18n)

	flagSet.IntVar(&cmd.port, "p", 3000, "Port the server is listening")
	flagSet.IntVar(&cmd.port, "port", 3000, "Port the server is listening")

	flagSet.BoolVar(&cmd.doUpdate, "u", false, "Automatic update of browser")
	flagSet.BoolVar(&cmd.doUpdate, "update", false, "Automatic update of browser")

	if err := flagSet.Parse(args); err != nil {
		return err
	}

	if flagSet.NArg() > 0 {
		printServeUsageMessage()

		return nutsandbolts.ErrInvalidNumberOfCommands
	}

	return nil
}

func (cmd *ServeCommand) Execute() error {
	if cmd.doUpdate {
		fmt.Println("************** Execute serve on port:", cmd.port, " with Update ON **************")
	} else {
		fmt.Println("************** Execute serve on port:", cmd.port, " with Update OFF **************")
	}

	return nil
}

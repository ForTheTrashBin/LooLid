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

var serveUsageMessage string

func printServeUsageMessage() {
	fmt.Fprintf(os.Stderr, serveUsageMessage)
}

//-----------------------------------------------------------------------------

type ServeCommand struct {
	localizer   *i18n.Localizer
	pureAppName string
	port        int
}

func (cmd *ServeCommand) GetName() string {
	return "serve"
}

func (cmd *ServeCommand) Parse(localizer *i18n.Localizer, pureAppName string, args []string) error {

	cmd.localizer = localizer
	cmd.pureAppName = pureAppName

	serveUsageMessage, _ = localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "flag." + cmd.GetName() + "UsageMessage",
		TemplateData: map[string]string{"pureAppName": pureAppName},
	})

	flagSet := flag.NewFlagSet(cmd.GetName(), flag.ContinueOnError)

	flagSet.Usage = printServeUsageMessage

	flagSet.SetOutput(io.Discard) // prevent output from flag-library (no i18n)

	flagSet.IntVar(&cmd.port, "p", 3000, "Port the server is listening")
	flagSet.IntVar(&cmd.port, "port", 3000, "Port the server is listening")

	if err := flagSet.Parse(args); err != nil {
		return err
	}

	if flagSet.NArg() > 0 {
		printServeUsageMessage()

		return errors.New("invalid number of commands")
	}

	return nil
}

func (cmd *ServeCommand) Execute() error {
	fmt.Println("**************Execute serve on port:", cmd.port, " **************")

	return nil
}

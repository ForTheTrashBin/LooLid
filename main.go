package main

import (
	"LooLid/commands"
	"LooLid/helper"
	"embed"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

//-----------------------------------------------------------------------------
// Enbed message-files direct into the app, so NO externel files are needed
//-----------------------------------------------------------------------------

//go:embed locales/active.*.toml
var LocalFS embed.FS

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

var mainUsageMessage string

func printMainUsageMessage() {
	fmt.Fprintf(os.Stderr, mainUsageMessage)
}

//-----------------------------------------------------------------------------

type commandItem interface {
	GetName() string

	Parse(
		localizer *i18n.Localizer,
		pureAppName string,
		arguments []string,
	) error

	Execute() error
}

func testableMain(args []string) int {

	pureAppName := helper.GetPureAppName()

	//-------------------------------------------------------------------------
	// Initialize i18n-bundle to use for the lifetime of the application
	//-------------------------------------------------------------------------

	bundle := i18n.NewBundle(language.English)

	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal) // Bundle now can read yaml-formatted message-files
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal) // Bundle now can read toml-formatted message-files

	//-------------------------------------------------------------------------
	// Enbed message-files direct into the app, so NO externel files are needed
	//-------------------------------------------------------------------------

	if _, err := bundle.LoadMessageFileFS(LocalFS, "locales/active.en.toml"); err != nil {
		panic(err)
	}

	if _, err := bundle.LoadMessageFileFS(LocalFS, "locales/active.de.toml"); err != nil {
		panic(err)
	}

	localizer := i18n.NewLocalizer(bundle, language.German.String())

	templateData := map[string]string{"pureAppName": pureAppName}

	mainUsageMessage, _ = localizer.Localize(&i18n.LocalizeConfig{TemplateData: templateData, MessageID: "flag.mainUsageMessage"})

	flagSet := flag.NewFlagSet(pureAppName, flag.ContinueOnError)

	flagSet.Usage = printMainUsageMessage

	if err := flagSet.Parse(args[1:]); err != nil {
		if err == flag.ErrHelp {
			return 2
		} else {
			return 1
		}
	}

	if flagSet.NArg() == 0 {
		printMainUsageMessage()
		return 0 // 2
	}

	commandItems := []commandItem{
		&commands.BuildCommand{},
		&commands.ServeCommand{},
		&commands.VersionCommand{},
	}

	commandName := flagSet.Arg(0)

	for _, commandItem := range commandItems {
		if commandItem.GetName() == commandName {
			if err := commandItem.Parse(localizer, pureAppName, flagSet.Args()[1:]); err != nil {
				if err != flag.ErrHelp {
					messageId, varItems := helper.GetFlagMessage(err)

					templateData := make(map[string]string)

					for index, varItem := range varItems {
						templateData["value"+strconv.Itoa(index+1)] = varItem
					}

					myMessage, _ := localizer.Localize(
						&i18n.LocalizeConfig{
							MessageID:    messageId,
							TemplateData: templateData,
						})

					fmt.Fprintln(os.Stderr, myMessage)

					return 0 // 1
				} else {
					return 0
				}
			}

			if err := commandItem.Execute(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return 0 // 1
			}

			return 0
		}
	}

	printMainUsageMessage()

	unknownSubcommand, _ := localizer.Localize(
		&i18n.LocalizeConfig{
			MessageID: "flag.unknownSubCommand",
			TemplateData: map[string]string{
				"pureAppName": pureAppName,
				"commandName": commandName,
			},
		})

	fmt.Fprintf(os.Stderr, unknownSubcommand)

	return 0 // 1
}

func main() {
	os.Exit(testableMain(os.Args))
}

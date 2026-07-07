package main

import (
	"LooLid/commands"
	"LooLid/helper"
	"LooLid/helper/locales"
	"embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
	//
	// Englich is the bundle default language.
	// Every unsupported language automatically falls back to English.
	//-------------------------------------------------------------------------

	bundle := i18n.NewBundle(language.English)

	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal) // Bundle now can read yaml-formatted message-files
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal) // Bundle now can read toml-formatted message-files

	//-------------------------------------------------------------------------
	// Enbed message-files direct into the app, so NO externel files are needed
	//-------------------------------------------------------------------------

	helper.Must(bundle.LoadMessageFileFS(LocalFS, "locales/active.en.toml"))
	helper.Must(bundle.LoadMessageFileFS(LocalFS, "locales/active.de.toml"))

	//-------------------------------------------------------------------------
	// read and prepare 'mainUsageMessage'
	//-------------------------------------------------------------------------

	systemLanguage, err := locales.GetLanguage()

	if err != nil {
		systemLanguage = language.English.String()
	}

	localizer := i18n.NewLocalizer(bundle, systemLanguage)

	mainUsageMessage = localizer.MustLocalize(
		&i18n.LocalizeConfig{
			MessageID:    "flag.mainUsageMessage",
			TemplateData: map[string]string{"pureAppName": pureAppName},
		})

	//-------------------------------------------------------------------------

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
		return 2
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

					myMessage := localizer.MustLocalize(
						&i18n.LocalizeConfig{
							MessageID:    messageId,
							TemplateData: templateData,
						})

					fmt.Fprintln(os.Stderr, myMessage)

					return 1
				} else {
					return 0
				}
			}

			if err := commandItem.Execute(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return 1
			}

			return 0
		}
	}

	printMainUsageMessage()

	unknownSubcommand := localizer.MustLocalize(
		&i18n.LocalizeConfig{
			MessageID: "flag.unknownSubCommand",
			TemplateData: map[string]string{
				"pureAppName": pureAppName,
				"commandName": commandName,
			},
		})

	fmt.Fprintf(os.Stderr, unknownSubcommand)

	return 1
}

//-----------------------------------------------------------------------------
// Supress private information, when printing the stack-trace at panic()
// Use "go build -trimpath ." to avoid having privat information in binary
//-----------------------------------------------------------------------------

func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "\x1b[31mInternal error:\x1b[0m %v\n", r)

			pcs := make([]uintptr, 32)

			n := runtime.Callers(3, pcs)

			frames := runtime.CallersFrames(pcs[:n])

			for {
				frame, more := frames.Next()

				fmt.Printf("%s:%d\n",
					filepath.Base(frame.File),
					frame.Line)

				if !more {
					break
				}
			}
		}
	}()

	testableMain(os.Args)
}

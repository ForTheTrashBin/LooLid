package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/ForTheTrashBin/LooLid/commands"
	"github.com/ForTheTrashBin/LooLid/helper/locales"
	"github.com/ForTheTrashBin/LooLid/helper/nutsandbolts"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.yaml.in/yaml/v3"
	"golang.org/x/text/language"
)

//-----------------------------------------------------------------------------
// Embed message files directly into the app, so NO external files are needed
//-----------------------------------------------------------------------------

//go:embed locales/active.*.toml
var LocalFS embed.FS

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

var mainUsageMessage string

func printMainUsageMessage() {

	fmt.Fprint(os.Stderr, mainUsageMessage)
}

//-----------------------------------------------------------------------------

func printLocalizedFlagMessage(localizer *i18n.Localizer, err error) {

	messageId, varItems := nutsandbolts.GetFlagMessage(err)

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

}

type commandItem interface {
	GetName() string

	Parse(
		localizer *i18n.Localizer,
		arguments []string,
	) error

	Execute() error
}

func testableMain(args []string) int {

	//-------------------------------------------------------------------------
	// Initialize i18n-bundle to use for the lifetime of the application
	//
	// English is the bundle default language.
	// Every unsupported language automatically falls back to English.
	//-------------------------------------------------------------------------

	bundle := i18n.NewBundle(language.English)

	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal) // Bundle now can read yaml-formatted message-files
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal) // Bundle now can read toml-formatted message-files

	//-------------------------------------------------------------------------
	// Embed message files directly into the app, so NO external files are needed
	//-------------------------------------------------------------------------

	nutsandbolts.Must(bundle.LoadMessageFileFS(LocalFS, "locales/active.en.toml"))
	nutsandbolts.Must(bundle.LoadMessageFileFS(LocalFS, "locales/active.de.toml"))

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

			MessageID: "flag.mainUsageMessage",
			TemplateData: map[string]string{
				"pureAppName":    nutsandbolts.GetPureAppName(),
				"configFileName": nutsandbolts.GetConfigFileName()},
		})

	//-------------------------------------------------------------------------

	flagSet := flag.NewFlagSet(nutsandbolts.GetPureAppName(), flag.ContinueOnError)

	flagSet.Usage = printMainUsageMessage

	flagSet.SetOutput(io.Discard) // prevent output from flag-library (no i18n)

	if err := flagSet.Parse(args[1:]); err != nil {

		if err != flag.ErrHelp {

			printLocalizedFlagMessage(localizer, err)

			return 1
		} else {

			return 0
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

			if err := commandItem.Parse(localizer, flagSet.Args()[1:]); err != nil {

				if err != flag.ErrHelp {

					printLocalizedFlagMessage(localizer, err)

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

				"value1": commandName,
			},
		})

	fmt.Fprintf(os.Stderr, "%s\n", unknownSubcommand)

	return 1
}

//-----------------------------------------------------------------------------
// Suppress private information when printing the stack trace at panic()
// Use "go build -trimpath ." to avoid having private information in the binary
//-----------------------------------------------------------------------------

func main() {

	defer func() {

		if recover := recover(); recover != nil {

			fmt.Fprintf(os.Stderr, "\x1b[31mInternal error:\x1b[0m %v\n", recover)

			pcs := make([]uintptr, 32)

			n := runtime.Callers(3, pcs)

			frames := runtime.CallersFrames(pcs[:n])

			for {

				frame, more := frames.Next()

				fmt.Printf("%s:%d\n", filepath.Base(frame.File), frame.Line)

				if !more {

					break
				}
			}
		}
	}()

	testableMain(os.Args)
}

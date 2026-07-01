package main

import (
	"LooLid/Commands"
	"embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

// ----------------
// Test ist ein Test
// ...und das ist noch ein Test
// Und das ist ein weiterer Test
// ------------------
func getPureAppName() string {

	//-------------------------------------------------------------------------
	// "os.Executable()" seems to be more reliable, so only use "os.Args" when nessesarry
	//-------------------------------------------------------------------------

	executablePath, err := os.Executable()

	if err != nil {
		executablePath = os.Args[0]
	}

	//-------------------------------------------------------------------------

	executablePath = executablePath + ".com.exe" // TODO

	executableBase := filepath.Base(executablePath) // Could be "app.exe", "app.v2.exe", "app.com", "app.v2", "app", ...

	//-------------------------------------------------------------------------
	// Remove repeatedly if there are multiple "renamed executable" extensions appended to the end.
	//-------------------------------------------------------------------------

	trimmed := true

	executableEndings := []string{".exe", ".com", ".cmd", ".bat"}

	for trimmed {
		trimmed = false

		ext := filepath.Ext(executableBase)

		for _, executableEnding := range executableEndings {
			if strings.EqualFold(ext, executableEnding) {
				executableBase = strings.TrimSuffix(executableBase, ext)
				trimmed = true
				break
			}
		}
	}

	return executableBase
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

var mainUsageMessage string

func mainUsage() {
	fmt.Fprintf(os.Stderr, mainUsageMessage)
}

//-----------------------------------------------------------------------------

type command interface {
	GetName() string

	Parse(
		localizer *i18n.Localizer,
		pureAppName string,
		arguments []string,
	) error

	Execute() error
}

func testableMain(args []string) int {

	pureAppName := getPureAppName()

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

	localizer := i18n.NewLocalizer(bundle, language.English.String())

	templateData := map[string]string{"pureAppName": pureAppName}

	mainUsageMessage, _ = localizer.Localize(&i18n.LocalizeConfig{TemplateData: templateData, MessageID: "flag.mainUsageMessage"})

	flagSet := flag.NewFlagSet(pureAppName, flag.ContinueOnError)

	flagSet.Usage = mainUsage

	if err := flagSet.Parse(args[1:]); err != nil {
		if err == flag.ErrHelp {
			return 2
		} else {
			return 1
		}
	}

	if flagSet.NArg() == 0 {
		mainUsage()
		return 2
	}

	commands := []command{
		&Commands.BuildCommand{},
		&Commands.ServeCommand{},
		&Commands.VersionCommand{},
	}

	commandName := flagSet.Arg(0)

	for _, command := range commands {
		if command.GetName() == commandName {
			if err := command.Parse(localizer, pureAppName, flagSet.Args()[1:]); err != nil {
				if err != flag.ErrHelp {
					/*
										"bad flag syntax: %s", s
										"flag provided but not defined: -%s", name
										"invalid boolean value %q for -%s: %v", value, name, err
										"invalid boolean flag %s: %v", name, err
										"flag needs an argument: -%s", name
										"invalid value %q for flag -%s: %v", value, name, err
						                "Zuviele Kommandos"
					*/
					fmt.Fprintln(os.Stderr, err)

					return 1
				} else {
					return 0
				}
			}
			if err := command.Execute(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return 1
			}
			return 0
		}
	}

	unknownSubcommand, _ := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "flag.unknownSubCommand",
		TemplateData: map[string]string{
			"pureAppName": pureAppName,
			"commandName": commandName,
		},
	})

	fmt.Fprintf(os.Stderr, unknownSubcommand)

	return 1
}

func main() {
	os.Exit(testableMain(os.Args))
}

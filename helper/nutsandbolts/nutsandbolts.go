package nutsandbolts

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ForTheTrashBin/LooLid/helper/constants"
)

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

var ErrInvalidNumberOfCommands = errors.New("invalid number of commands")

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

func GetPureAppName() string {

	//-------------------------------------------------------------------------
	// "os.Executable()" seems to be more reliable, so only use "os.Args" when nessesarry
	//-------------------------------------------------------------------------

	executablePath, err := os.Executable()

	if err != nil {
		executablePath = os.Args[0]
	}

	//-------------------------------------------------------------------------

	executableBase := filepath.Base(executablePath) // Could be "app.exe", "app.v2.exe", "app.com", "app.v2", "app", ...

	//-------------------------------------------------------------------------
	// Remove repeatedly if there are multiple "(renamed) executable" extensions appended to the end.
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

func GetConfigFileName() string {

	executableBase := GetPureAppName()

	executableBase, _ = strings.CutPrefix(executableBase, ".")
	executableBase, _ = strings.CutSuffix(executableBase, ".")

	configFileExtension := constants.AppConfig_ConfigFileExtension

	configFileExtension, _ = strings.CutPrefix(configFileExtension, ".")
	configFileExtension, _ = strings.CutSuffix(configFileExtension, ".")

	return executableBase + "." + configFileExtension
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

/*
*** "bad flag syntax: %s", s
*** "flag provided but not defined: -%s", name
*** "invalid boolean value %q for -%s: %v", value, name, err
*** "invalid boolean flag %s: %v", name, err
*** "flag needs an argument: -%s", name
*** "invalid value %q for flag -%s: %v", value, name, err
*** "invalid number of commands"
 */

func GetFlagMessage(err error) (string, []string) {

	if after, found := strings.CutPrefix(err.Error(), "bad flag syntax: "); found {

		return "flag.bad_flag_syntax", []string{after}
	} else {

		if after, found := strings.CutPrefix(err.Error(), "flag provided but not defined: "); found {

			return "flag.flag_provided_but_not_defined", []string{after}
		} else {

			if after, found := strings.CutPrefix(err.Error(), "invalid boolean value "); found {

				var values []string = make([]string, 2)

				if before, after, found := strings.Cut(after, " for "); found {

					values[0] = before

					if before, _, found := strings.Cut(after, ":"); found {
						values[1] = before
					}
				}

				return "flag.invalid_boolean_value", values
			} else {

				if after, found := strings.CutPrefix(err.Error(), "invalid boolean flag "); found {

					var values []string = make([]string, 2)

					if before, after, found := strings.Cut(after, ": "); found {

						values[0] = before
						values[1] = after
					}

					return "flag.invalid_boolean_flag", []string{after}
				} else {

					if after, found := strings.CutPrefix(err.Error(), "flag needs an argument: "); found {

						return "flag.flag_needs_an_argument", []string{after}
					} else {

						if after, found := strings.CutPrefix(err.Error(), "invalid value "); found {

							var values []string = make([]string, 2)

							if before, after, found := strings.Cut(after, " for flag "); found {

								values[0] = before

								if before, _, found := strings.Cut(after, ":"); found {

									values[1] = before
								}
							}

							return "flag.invalid_value", values
						} else {

							if after, found = strings.CutPrefix(err.Error(), ErrInvalidNumberOfCommands.Error()); found {

								return "flag.invalid_number_of_commands", nil
							} else {

								return "", nil
							}
						}
					}
				}
			}
		}
	}
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

func SortFileInfos(fileInfos []os.FileInfo, directoriesFirst bool) {

	sort.Slice(fileInfos, func(i, j int) bool {

		iDir := fileInfos[i].IsDir()
		jDir := fileInfos[j].IsDir()

		if iDir != jDir {

			return directoriesFirst == iDir
		}

		return fileInfos[i].Name() < fileInfos[j].Name()
	})
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

func Must[T any](x T, err error) T {
	if err != nil {
		panic(err)
	}
	return x
}

package inputFolderRead

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ForTheTrashBin/LooLid/blackwhite"
	"github.com/ForTheTrashBin/LooLid/helper/nutsandbolts"
	"github.com/ForTheTrashBin/LooLid/helper/osspecific"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/afero"
)

func copyDir(localizer *i18n.Localizer, sourceFs afero.Fs, sourceFolder string, destFs afero.Fs, destFolder string, bwConfig blackwhite.Config) error {

	var doDebug bool = false

	if doDebug {

		fmt.Println("*********************************************************************************")
		fmt.Print("*** copyDir sourceFolder <", sourceFolder, "> destFolder <", destFolder, ">\n")
		fmt.Println("*********************************************************************************")
	}

	sourceFileInfos, err := afero.ReadDir(sourceFs, sourceFolder)

	if err == nil {

		nutsandbolts.SortFileInfos(sourceFileInfos, false)

		configFileName := nutsandbolts.GetConfigFileName()

		for idx := range sourceFileInfos {

			if doDebug {

				fmt.Println("*** Search for", configFileName, "in", idx, sourceFileInfos[idx].Name(), sourceFileInfos[idx].IsDir())
			}

			if !sourceFileInfos[idx].IsDir() && sourceFileInfos[idx].Name() == configFileName {

				localConfig, ok, err := blackwhite.LoadConfigFile(localizer, filepath.Join(sourceFolder, configFileName))

				if err == nil {

					if ok {

						// Create a local independent copy of bwConfig so modifications
						// (Merge/assignment) for this folder do not affect the caller
						// after recursion returns. We must copy the slices to avoid
						// sharing the underlying array.
						localBW := bwConfig
						localBW.Blacklist.Rules = append([]blackwhite.Rule(nil), bwConfig.Blacklist.Rules...)
						localBW.Whitelist.Rules = append([]blackwhite.Rule(nil), bwConfig.Whitelist.Rules...)

						if localConfig.DoBlacklistInherit() {

							localBW.Blacklist.MergeRuleSet(&localConfig.Blacklist)
						} else {

							localBW.Blacklist = localConfig.Blacklist
						}

						if localConfig.DoWhitelistInherit() {

							localBW.Whitelist.MergeRuleSet(&localConfig.Whitelist)
						} else {

							localBW.Whitelist = localConfig.Whitelist
						}

						// Use the local copy for subsequent recursion and processing
						bwConfig = localBW
					}
				} else {

					return err
				}
				//Found

				break
			}
		}

		for idx, sourceFileInfo := range sourceFileInfos {

			if doDebug {

				fmt.Println("*** Index                :", idx)
				fmt.Println("*** sourceFolder         :", sourceFolder)
				fmt.Println("*** destFolder           :", destFolder)
				fmt.Println("*** sourceFileInfo.IsDir :", sourceFileInfo.IsDir())
				fmt.Println("*** sourceFileInfo.Name  :", sourceFileInfo.Name())
				fmt.Println("*** sourceFileInfo.Size  :", sourceFileInfo.Size())
				fmt.Println("*** sourceFileInfo.Mode  :", sourceFileInfo.Mode())
			}

			isListed := bwConfig.IsListed(sourceFileInfo.Name(), sourceFileInfo.IsDir())

			if !isListed {

				if sourceFileInfo.IsDir() {

					if !osspecific.IsHiddenOrSystem(sourceFolder) { // TODO: Path?

						newSourcePath := filepath.Join(sourceFolder, sourceFileInfo.Name())
						newDestPath := filepath.Join(destFolder, sourceFileInfo.Name())

						if err = destFs.Mkdir(newDestPath, os.ModePerm); err == nil {

							if err = copyDir(localizer, sourceFs, newSourcePath, destFs, newDestPath, bwConfig); err == nil {

								if err = destFs.Chmod(newDestPath, sourceFileInfo.Mode()); err != nil {

									return err
								}
							} else {

								return err
							}
						} else {

							return err
						}
					}
				} else {

					if err = copyFile(
						sourceFs,
						filepath.Join(sourceFolder, sourceFileInfo.Name()),
						destFs,
						filepath.Join(destFolder, sourceFileInfo.Name())); err != nil {

						return err
					}
				}
			}
		}
	}

	return err
}

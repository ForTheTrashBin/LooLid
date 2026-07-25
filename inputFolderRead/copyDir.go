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

func copyDir(localizer *i18n.Localizer, sourceFs afero.Fs, sourceFolder string, destFs afero.Fs, destFolder string, bwConfig *blackwhite.BWConfig, depth int) error {

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

		deepCopyBWConfig := bwConfig

		for idx := range sourceFileInfos {

			if doDebug {

				fmt.Println("*** Search for", configFileName, "in", idx, sourceFileInfos[idx].Name(), sourceFileInfos[idx].IsDir())
			}

			if !sourceFileInfos[idx].IsDir() && sourceFileInfos[idx].Name() == configFileName {

				localBWConfig, err := blackwhite.NewBWConfigFromFile(localizer, filepath.Join(sourceFolder, configFileName), depth)

				if err != nil {

					return err
				}

				deepCopyBWConfig = blackwhite.NewBWConfigFromBWConfig(bwConfig) // Deep copy

				if localBWConfig.DoBlacklistInherit() {

					deepCopyBWConfig.Blacklist.MergeRuleSet(&localBWConfig.Blacklist)
				} else {

					deepCopyBWConfig.Blacklist = blackwhite.NewRuleSetFromRuleSet(&localBWConfig.Blacklist)
				}

				if localBWConfig.DoWhitelistInherit() {

					deepCopyBWConfig.Whitelist.MergeRuleSet(&localBWConfig.Whitelist)
				} else {

					deepCopyBWConfig.Whitelist = blackwhite.NewRuleSetFromRuleSet(&localBWConfig.Whitelist)
				}

				break // Should be only ONE config-file per folder, so we can stop searching
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

			if !deepCopyBWConfig.IsListed(sourceFileInfo.Name(), sourceFileInfo.IsDir()) {

				if sourceFileInfo.IsDir() {

					if !osspecific.IsHiddenOrSystem(sourceFolder) { // TODO: Path?

						newSourcePath := filepath.Join(sourceFolder, sourceFileInfo.Name())
						newDestPath := filepath.Join(destFolder, sourceFileInfo.Name())

						if err = destFs.Mkdir(newDestPath, os.ModePerm); err == nil {

							if err = copyDir(localizer, sourceFs, newSourcePath, destFs, newDestPath, deepCopyBWConfig, depth+1); err == nil {

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

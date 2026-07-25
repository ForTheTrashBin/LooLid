package inputFolderRead

import (
	"os"
	"path/filepath"

	"github.com/ForTheTrashBin/LooLid/blackwhite"
	"github.com/ForTheTrashBin/LooLid/helper/nutsandbolts"
	"github.com/ForTheTrashBin/LooLid/helper/osspecific"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/afero"
)

func copyDir(localizer *i18n.Localizer, sourceFs afero.Fs, sourceFolder string, destFs afero.Fs, destFolder string, bwConfig *blackwhite.BWConfig, depth int) error {

	//-------------------------------------------------------------------------
	// Get all files and directories in the source folder
	//-------------------------------------------------------------------------

	sourceFileInfos, err := afero.ReadDir(sourceFs, sourceFolder)

	//-------------------------------------------------------------------------

	if err == nil {

		deepCopyBWConfig := bwConfig // Copy of pointers, NOT a deep-copy

		configFileName := nutsandbolts.GetConfigFileName()

		//---------------------------------------------------------------------

		nutsandbolts.SortFileInfos(sourceFileInfos, false)

		//---------------------------------------------------------------------
		// Iterate over all files and directories and try to find a config file
		//---------------------------------------------------------------------

		for _, sourceFileInfo := range sourceFileInfos {

			if !sourceFileInfo.IsDir() && (sourceFileInfo.Name() == configFileName) {

				localBWConfig, err := blackwhite.NewBWConfigFromFile(localizer, filepath.Join(sourceFolder, configFileName), depth)

				if err != nil {

					return err
				}

				//-------------------------------------------------------------

				deepCopyBWConfig = blackwhite.NewBWConfigFromBWConfig(bwConfig) // Deep copy now

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

				break // There should be only ONE config-file per folder, so we can stop searching
			}
		}

		for _, sourceFileInfo := range sourceFileInfos {

			if !deepCopyBWConfig.IsListed(sourceFileInfo.Name(), sourceFileInfo.IsDir()) {

				if !osspecific.IsHiddenOrSystem(sourceFolder) {

					newSource := filepath.Join(sourceFolder, sourceFileInfo.Name())
					newDest := filepath.Join(destFolder, sourceFileInfo.Name())

					if sourceFileInfo.IsDir() {

						if err = destFs.Mkdir(newDest, os.ModePerm); err != nil {

							return err
						}

						if err = copyDir(localizer, sourceFs, newSource, destFs, newDest, deepCopyBWConfig, depth+1); err != nil {

							return err
						}

						if err = destFs.Chmod(newDest, sourceFileInfo.Mode()); err != nil {

							return err
						}
					} else {

						if err = copyFile(sourceFs, newSource, destFs, newDest); err != nil {

							return err
						}
					}
				}
			}
		}
	}

	return err
}

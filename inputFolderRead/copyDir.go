package inputFolderRead

import (
	"os"
	"path/filepath"

	"github.com/ForTheTrashBin/LooLid/configParser"
	"github.com/ForTheTrashBin/LooLid/helper"
	"github.com/ForTheTrashBin/LooLid/helper/nutsandbolts"
	"github.com/ForTheTrashBin/LooLid/helper/osspecific"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/afero"
)

func copyDir(localizer *i18n.Localizer, sourceFs afero.Fs, sourceFolder string, destFs afero.Fs, destFolder string, rulesConfig *configParser.RulesConfig, depth int) error {

	//-------------------------------------------------------------------------
	// Get all files and directories in the source folder
	//-------------------------------------------------------------------------

	sourceFileInfos, err := afero.ReadDir(sourceFs, sourceFolder)

	//-------------------------------------------------------------------------

	if err == nil {

		deepCopyRulesConfig := rulesConfig // Copy of pointers, NOT a deep-copy

		configFileName := nutsandbolts.GetConfigFileName()

		//---------------------------------------------------------------------

		nutsandbolts.SortFileInfos(sourceFileInfos, false)

		//---------------------------------------------------------------------
		// Iterate over all files and directories and try to find a config file
		//---------------------------------------------------------------------

		for _, sourceFileInfo := range sourceFileInfos {

			if !sourceFileInfo.IsDir() && (sourceFileInfo.Name() == configFileName) {

				localRulesConfig, err := configParser.NewRulesConfigFromFile(localizer, filepath.Join(sourceFolder, configFileName), depth)

				if err != nil {

					return err
				}

				//-------------------------------------------------------------

				deepCopyRulesConfig = configParser.NewRulesConfigFromRulesConfig(rulesConfig) // Deep copy now

				if localRulesConfig.DoBlacklistInherit() {

					deepCopyRulesConfig.Blacklist.MergeRuleSet(&localRulesConfig.Blacklist)
				} else {

					deepCopyRulesConfig.Blacklist = configParser.CloneRuleSet(&localRulesConfig.Blacklist)
				}

				if localRulesConfig.DoWhitelistInherit() {

					deepCopyRulesConfig.Whitelist.MergeRuleSet(&localRulesConfig.Whitelist)
				} else {

					deepCopyRulesConfig.Whitelist = configParser.CloneRuleSet(&localRulesConfig.Whitelist)
				}

				deepCopyRulesConfig.FrontMatter = configParser.CloneFrontMatter(localRulesConfig.FrontMatter)

				if len(deepCopyRulesConfig.FrontMatter) > 0 {
					frontMatterFile := filepath.Join(destFolder, "loolid-frontmatter")
					frontMatterData, err := helper.SerializeFrontMatter(deepCopyRulesConfig.FrontMatter)
					if err != nil {
						return err
					}
					if err := afero.WriteFile(destFs, frontMatterFile, frontMatterData, 0o600); err != nil {
						return err
					}
				}

				break // There should be only ONE config-file per folder, so we can stop searching
			}
		}

		for _, sourceFileInfo := range sourceFileInfos {

			if !deepCopyRulesConfig.IsListed(sourceFileInfo.Name(), sourceFileInfo.IsDir()) {

				if !osspecific.IsHiddenOrSystem(sourceFolder) {

					newSource := filepath.Join(sourceFolder, sourceFileInfo.Name())
					newDest := filepath.Join(destFolder, sourceFileInfo.Name())

					if sourceFileInfo.IsDir() {

						if err = destFs.Mkdir(newDest, os.ModePerm); err != nil {

							return err
						}

						if err = copyDir(localizer, sourceFs, newSource, destFs, newDest, deepCopyRulesConfig, depth+1); err != nil {

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

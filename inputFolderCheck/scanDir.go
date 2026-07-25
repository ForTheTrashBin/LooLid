package inputFolderCheck

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ForTheTrashBin/LooLid/blackwhite"
	"github.com/ForTheTrashBin/LooLid/helper/nutsandbolts"
	"github.com/ForTheTrashBin/LooLid/helper/osspecific"
	"github.com/spf13/afero"
	"golang.org/x/text/unicode/norm"
)

func (chk *checker) scanDir(sourceFs afero.Fs, sourceFolder string, bwConfig *blackwhite.BWConfig, depth int) error {

	//-------------------------------------------------------------------------
	// Get all files and directories in the source folder
	//-------------------------------------------------------------------------

	sourceFileInfos, err := afero.ReadDir(sourceFs, sourceFolder)

	//-------------------------------------------------------------------------

	if err == nil {

		configFileName := nutsandbolts.GetConfigFileName()

		configFileNameLower := strings.ToLower(configFileName)

		//---------------------------------------------------------------------

		nutsandbolts.SortFileInfos(sourceFileInfos, false) // Sort fileInfos, directories last

		//---------------------------------------------------------------------
		// Iterate over all files and directories in the source folder and validate them on per-item basis
		//---------------------------------------------------------------------

		var hasDirs bool = false // Keep track if there are any directories in the source folder

		for _, fileInfo := range sourceFileInfos {

			filePath := filepath.Join(sourceFolder, fileInfo.Name())

			if !osspecific.IsHiddenOrSystem(filePath) {

				//-------------------------------------------------------------
				// Validate a single file or directory and add any issues to the checker's issue lists
				//-------------------------------------------------------------

				chk.validateEntry(Entry{

					filePath:  filePath,
					entryName: fileInfo.Name(),
					isDir:     fileInfo.IsDir(),
					isSymlink: (fileInfo.Mode() & os.ModeSymlink) != 0,
				})

				//-------------------------------------------------------------
				// Check for ambiguous filenames of configuration-files
				//-------------------------------------------------------------

				testEntryName, _ := strings.CutPrefix(strings.ToLower(fileInfo.Name()), ".")

				if fileInfo.IsDir() {

					hasDirs = true

					if testEntryName == configFileNameLower {

						chk.issuesConfigFile = append(chk.issuesConfigFile, Issue{isDir: true, filePath: filePath})
					}
				} else {

					if fileInfo.Name() != configFileName {

						if testEntryName == configFileNameLower {

							chk.issuesConfigFile = append(chk.issuesConfigFile, Issue{isDir: false, filePath: filePath})
						}
					}
				}
			}
		}

		//---------------------------------------------------------------------
		//---------------------------------------------------------------------

		sourceFileInfoCount := len(sourceFileInfos)

		if sourceFileInfoCount >= 2 { // There have to be at least 2 items in a folder to have duplicates

			duplicateGroup := DuplicateGroup{groupName: sourceFolder}

			for _, fileInfo := range sourceFileInfos {

				duplicateGroup.items = append(duplicateGroup.items, DuplicateItem{

					itemName: fileInfo.Name(),
					isDir:    fileInfo.IsDir(),
				})
			}

			for memberIdx := len(duplicateGroup.items) - 1; memberIdx >= 0; memberIdx-- {

				groupMember, _ := strings.CutPrefix(strings.ToLower(duplicateGroup.items[memberIdx].itemName), ".")

				var siblingsFound bool = false

				for testerIdx := 0; !siblingsFound && (testerIdx < len(duplicateGroup.items)); testerIdx++ {

					if memberIdx != testerIdx {

						groupMemberTest, _ := strings.CutPrefix(strings.ToLower(duplicateGroup.items[testerIdx].itemName), ".")

						if groupMember == groupMemberTest {

							siblingsFound = true
						} else {

							if norm.NFC.String(groupMember) == norm.NFC.String(groupMemberTest) {

								siblingsFound = true
							}
						}
					}
				}

				if !siblingsFound {

					duplicateGroup.items = slices.Delete(duplicateGroup.items, memberIdx, memberIdx+1)
				}
			}

			if len(duplicateGroup.items) >= 2 {

				chk.duplicateGroups = append(chk.duplicateGroups, duplicateGroup)
			}
		}

		if hasDirs {

			deepCopyBWConfig := blackwhite.NewBWConfigFromBWConfig(bwConfig)

			for _, sourceFileInfo := range sourceFileInfos {

				if !sourceFileInfo.IsDir() && sourceFileInfo.Name() == configFileName {

					localBWConfig, err := blackwhite.NewBWConfigFromFile(chk.localizer, filepath.Join(sourceFolder, configFileName), depth)

					if err != nil {

						return err
					}

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

			for _, sourceFileInfo := range sourceFileInfos {

				if !deepCopyBWConfig.IsListed(sourceFileInfo.Name(), sourceFileInfo.IsDir()) {

					if sourceFileInfo.IsDir() {

						if !osspecific.IsHiddenOrSystem(sourceFolder) { // TODO: Path?

							newSourcePath := filepath.Join(sourceFolder, sourceFileInfo.Name())

							if err = chk.scanDir(sourceFs, newSourcePath, deepCopyBWConfig, depth+1); err != nil {

								return err
							}
						}
					}
				}
			}
		}
	}

	return err
}

package checkInputFolder

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ForTheTrashBin/LooLid/helper/nutsandbolts"
	"golang.org/x/text/unicode/norm"
)

func (chk *checker) scanDirectories() error {

	var entries []Entry

	if err := filepath.WalkDir(chk.inputFolder, func(path string, dirEntry os.DirEntry, err error) error {

		//---------------------------------------------------------------------
		// This error will abort the 'WalkDir-function
		//---------------------------------------------------------------------

		if err != nil {
			return err
		}

		//---------------------------------------------------------------------
		// Skip root, because it was checked before
		//---------------------------------------------------------------------

		if path == chk.inputFolder {
			return nil
		}

		//---------------------------------------------------------------------
		// Get the relative path (to root) of this entry
		//---------------------------------------------------------------------

		xrelativePath, err := filepath.Rel(chk.inputFolder, path)

		if err != nil {

			return err
		}

		xrelativePath = filepath.ToSlash(xrelativePath)

		//---------------------------------------------------------------------
		// Get fileInfo of this entry
		//---------------------------------------------------------------------

		fileInfo, err := dirEntry.Info()

		if err != nil {
			return err
		}

		//---------------------------------------------------------------------
		// Create an entry and save it to the list
		//---------------------------------------------------------------------

		entry := Entry{
			filePath:  path,
			entryName: dirEntry.Name(),
			isDir:     fileInfo.IsDir(),
			isSymlink: fileInfo.Mode()&os.ModeSymlink != 0,
		}

		entries = append(entries, entry)

		//---------------------------------------------------------------------
		// Check,if this ia a directory to be skiped
		//---------------------------------------------------------------------

		if entry.isDir {

			_, err = os.ReadDir(chk.inputFolder + string(filepath.Separator) + entry.filePath)

			if err != nil {

				if os.IsPermission(err) {

					return filepath.SkipDir
				}
			}

			if isHiddenOrSystemOnWindows(chk.inputFolder + string(filepath.Separator) + entry.filePath) {

				return filepath.SkipDir
			}
		}

		return nil
	}); err == nil {

		//---------------------------------------------------------------------
		// Run checks on individual files
		//---------------------------------------------------------------------

		chk.validateEntries(entries)

		//---------------------------------------------------------------------
		// Enter all folders and files into a 2D-map representing a folder tree
		//---------------------------------------------------------------------

		folders := make(map[string]map[string]bool)

		for _, entry := range entries {

			folderName := filepath.Dir(entry.filePath)

			folder, foundFolder := folders[folderName]

			if !foundFolder {

				folder = make(map[string]bool)

				folders[folderName] = folder
			}

			folder[entry.entryName] = entry.isDir
		}

		//---------------------------------------------------------------------

		if false {

			for keyFolder, folder := range folders {

				fmt.Println("Folder:", keyFolder)

				for keyFolderEntry, isDir := range folder {

					if isDir {

						fmt.Println("   Folder Entry DIR:", keyFolderEntry)
					} else {

						fmt.Println("   Folder Entry FIL:", keyFolderEntry)
					}
				}
			}
		}

		//---------------------------------------------------------------------
		// Check for ambiguous filename of configuration-files
		//---------------------------------------------------------------------

		configFileName := nutsandbolts.GetConfigFileName()

		configFileNameLower := strings.ToLower(configFileName)

		for keyFolder, folder := range folders {

			for keyFolderEntry, isDir := range folder { // entry.entryName, entry.isDir

				testEntryName, _ := strings.CutPrefix(strings.ToLower(keyFolderEntry), ".")

				if isDir {

					if testEntryName == configFileNameLower {

						chk.issuesConfigFile = append(
							chk.issuesConfigFile, Issue{
								isDir:    isDir,
								filePath: keyFolder + string(filepath.Separator) + keyFolderEntry,
							})
					}
				} else {

					if keyFolderEntry != configFileName {

						if testEntryName == configFileNameLower {

							chk.issuesConfigFile = append(
								chk.issuesConfigFile, Issue{
									isDir:    isDir,
									filePath: keyFolder + string(filepath.Separator) + keyFolderEntry,
								})
						}
					}
				}
			}
		}

		//---------------------------------------------------------------------
		// Compare 'DuplicateItems' in each 'DuplicateGroup' and delete non duplicates
		//---------------------------------------------------------------------

		for keyFolder, folder := range folders {

			if len(folder) >= 2 {

				duplicateGroup := DuplicateGroup{
					groupName: keyFolder,
				}

				for keyFolderItem, isDir := range folder { // entry.entryName, entry.isDir

					duplicateGroup.items = append(
						duplicateGroup.items,
						DuplicateItem{
							itemName: keyFolderItem,
							isDir:    isDir,
						})
				}

				chk.duplicateGroups = append(chk.duplicateGroups, duplicateGroup)
			}
		}

		//---------------------------------------------------------------------

		for looperGroups := len(chk.duplicateGroups) - 1; looperGroups >= 0; looperGroups-- {

			duplicateGroup := &chk.duplicateGroups[looperGroups]

			for looperMember := len(duplicateGroup.items) - 1; looperMember >= 0; looperMember-- {

				groupMember, _ := strings.CutPrefix(strings.ToLower(duplicateGroup.items[looperMember].itemName), ".")

				var siblingsFound bool = false

				for looperTest := 0; !siblingsFound && (looperTest < len(duplicateGroup.items)); looperTest++ {

					if looperMember != looperTest {

						groupMemberTest, _ := strings.CutPrefix(strings.ToLower(duplicateGroup.items[looperTest].itemName), ".")

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

					duplicateGroup.items = slices.Delete(duplicateGroup.items, looperMember, looperMember+1)
				}
			}
		}

		//---------------------------------------------------------------------
		// Cleanup of empty 'duplicateGroups'
		//---------------------------------------------------------------------

		for looperGroups := len(chk.duplicateGroups) - 1; looperGroups >= 0; looperGroups-- {
			if len(chk.duplicateGroups[looperGroups].items) <= 1 {
				chk.duplicateGroups = slices.Delete(chk.duplicateGroups, looperGroups, looperGroups+1)
			}
		}

		return err

	} else {

		return err
	}
}

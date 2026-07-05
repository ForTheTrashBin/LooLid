package checkInputFolder

import (
	"LooLid/helper"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"golang.org/x/text/unicode/norm"
)

func (chk *checker) scanDirectories() error {

	var entries []Entry

	if err := helper.WalkDir(chk.inputFolder, func(depth int, path string, dirEntry os.DirEntry, err error) error {

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

		relativePath, err := filepath.Rel(chk.inputFolder, path)

		if err != nil {
			return err
		}

		relativePath = filepath.ToSlash(relativePath)

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
			depth:        depth,
			relativePath: relativePath,
			entryName:    dirEntry.Name(),
			isDir:        fileInfo.IsDir(),
			isSymlink:    fileInfo.Mode()&os.ModeSymlink != 0,
		}

		entries = append(entries, entry)

		//---------------------------------------------------------------------
		// Check,if this ia a directory to be skiped
		//---------------------------------------------------------------------

		if entry.isDir {
			_, err = os.ReadDir(chk.inputFolder + string(filepath.Separator) + entry.relativePath)

			if err != nil {
				if os.IsPermission(err) {
					return filepath.SkipDir
				}
			}

			if isHiddenOrSystemOnWindows(chk.inputFolder + string(filepath.Separator) + entry.relativePath) {
				return filepath.SkipDir
			}
		}

		return nil
	},
	); err == nil {

		//---------------------------------------------------------------------
		// Sort eintries by depth of the directory structure and file/directory name
		//---------------------------------------------------------------------

		sort.Slice(entries, func(i, j int) bool {
			if entries[i].depth < entries[j].depth {
				return true
			} else {
				if entries[i].depth == entries[j].depth {
					return entries[i].relativePath < entries[j].relativePath
				}
			}

			return false
		})

		//---------------------------------------------------------------------
		// Run checks on individual files
		//---------------------------------------------------------------------

		chk.validateEntries(entries)

		//---------------------------------------------------------------------
		// Group entries by level an directory
		//---------------------------------------------------------------------

		var groupIndex int = -1
		var groupDepth int
		var groupDirPath string

		for index, entry := range entries {

			if groupIndex < 0 {

				groupIndex = index
				groupDepth = entry.depth
				groupDirPath = filepath.Dir(entry.relativePath)
			} else {

				if (groupDepth != entry.depth) || (groupDirPath != filepath.Dir(entry.relativePath)) {

					if (index - groupIndex) > 0 {

						duplicateGroup := DuplicateGroup{
							groupName: groupDirPath,
						}

						for looper := groupIndex; looper < index; looper++ {
							duplicateGroup.items = append(
								duplicateGroup.items,
								DuplicateItem{
									itemName: entries[looper].entryName,
									isDir:    entries[looper].isDir,
								})
						}

						chk.duplicateGroups = append(chk.duplicateGroups, duplicateGroup)
					}

					groupIndex = index
					groupDepth = entry.depth
					groupDirPath = filepath.Dir(entry.relativePath)
				}
			}
		}

		if (len(entries) - groupIndex) > 0 {
			duplicateGroup := DuplicateGroup{
				groupName: groupDirPath,
			}

			for looper := groupIndex; looper < len(entries); looper++ {
				duplicateGroup.items = append(
					duplicateGroup.items,
					DuplicateItem{
						itemName: entries[looper].entryName,
						isDir:    entries[looper].isDir,
					})
			}

			chk.duplicateGroups = append(chk.duplicateGroups, duplicateGroup)
		}

		//---------------------------------------------------------------------
		// Check for ambiguous filename of configuration-files
		//---------------------------------------------------------------------

		configFileName := helper.GetConfigFileName()

		configFileNameLower := strings.ToLower(configFileName)

		for _, duplicateGroup := range chk.duplicateGroups {

			for _, duplicateItem := range duplicateGroup.items {

				itemName, _ := strings.CutPrefix(strings.ToLower(duplicateItem.itemName), ".")

				if duplicateItem.isDir {

					if itemName == configFileNameLower {

						chk.issuesConfigFile = append(
							chk.issuesConfigFile, Issue{
								isDir: duplicateItem.isDir,
								path:  duplicateGroup.groupName + string(filepath.Separator) + duplicateItem.itemName,
							})
					}
				} else {

					if duplicateItem.itemName != configFileName {

						if itemName == configFileNameLower {

							chk.issuesConfigFile = append(
								chk.issuesConfigFile, Issue{
									isDir: duplicateItem.isDir,
									path:  duplicateGroup.groupName + string(filepath.Separator) + duplicateItem.itemName,
								})
						}
					}
				}
			}
		}

		//---------------------------------------------------------------------
		// Compare 'DuplicateItems' in each 'DuplicateGroup' and delete non duplicates
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

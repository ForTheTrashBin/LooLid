package FolderCheckerInput

import (
	"fmt"
	"os"
	"path/filepath"
)

func scanDirectory(root string) ([]Entry, []Issue, Statistics, error) {

	var entries []Entry
	var issues []Issue

	var stats Statistics

	err := filepath.WalkDir(
		root,
		func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if path == root {
				return nil
			}

			//-----------------------------------------------------------------

			info, err := entry.Info()

			if err != nil {
				return err
			}

			relative, err := filepath.Rel(root, path)

			if err != nil {
				return err
			}

			relative = filepath.ToSlash(relative)

			item := Entry{
				Path: relative,
				Name: entry.Name(),

				IsDir: info.IsDir(),

				IsSymlink: info.Mode()&os.ModeSymlink != 0,
			}

			fmt.Println("Entry:", item)

			entries = append(
				entries,
				item,
			)

			if item.IsDir {
				stats.Directories++
			} else {
				stats.Files++
			}

			if item.IsSymlink {

				issues = append(
					issues,
					Issue{
						Type: SymbolicLinkDetected,
						Path: relative,
						Info: "symbolic link detected",
					},
				)
			}

			if item.IsDir {
				_, err = os.ReadDir(root + "/" + item.Path)

				if err != nil {
					if os.IsPermission(err) {
						return filepath.SkipDir
					}
				}
			}

			return nil
		},
	)

	return entries, issues, stats, err

}

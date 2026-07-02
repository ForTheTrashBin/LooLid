package checkInputFolder

import (
	"os"
	"path/filepath"
)

func scanDirectory(root string) ([]Entry, error) {

	var entries []Entry

	err := filepath.WalkDir(
		root,
		func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			//-----------------------------------------------------------------
			// Skip root, because it was checked before
			//-----------------------------------------------------------------

			if path == root {
				return nil
			}

			//-----------------------------------------------------------------
			// Get the relative path (to root) of this entry
			//-----------------------------------------------------------------

			relativePath, err := filepath.Rel(root, path)

			if err != nil {
				return err
			}

			relativePath = filepath.ToSlash(relativePath)

			//-----------------------------------------------------------------
			// Get fileInfo of this entry
			//-----------------------------------------------------------------

			fileInfo, err := entry.Info()

			if err != nil {
				return err
			}

			//-----------------------------------------------------------------
			/*
				fmt.Println("Path         :", path)
				fmt.Println("RelativePath :", relativePath)

				if fileInfo.IsDir() {
					fmt.Println("IsDirectory  : true")
				} else {
					fmt.Println("IsDirectory  : false")
				}

				if fileInfo.Mode()&os.ModeSymlink != 0 {
					fmt.Println("IsSymLink    : true")
				} else {
					fmt.Println("IsSymLink    : false")
				}

				fmt.Println("********************************************")
			*/
			//-----------------------------------------------------------------
			// Check for "hidden" or "system" directories on Windows
			//-----------------------------------------------------------------
			/*
				fmt.Println("Path:", path)
				if entry.IsDir() {
					fmt.Println("It's a directory")
				}
			*/
			item := Entry{
				Path: relativePath,
				Name: entry.Name(),

				IsDir: fileInfo.IsDir(),

				IsSymlink: fileInfo.Mode()&os.ModeSymlink != 0,
			}

			entries = append(
				entries,
				item,
			)

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

	return entries, err

}

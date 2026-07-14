package inputFolderWrite

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ForTheTrashBin/LooLid/helper/nutsandbolts"
	"github.com/spf13/afero"
)

//-----------------------------------------------------------------------------
// Delete files/directories that are not (or no longer) present in the source
//-----------------------------------------------------------------------------

func cleanDir(destFs afero.Fs, destFolder string, sourceFs afero.Fs) error {

	var doDebug bool = false

	if doDebug {

		fmt.Println("*********************************************************************************")
		fmt.Print("*** cleanDestination destFolder: <", destFolder, ">\n")
		fmt.Println("*********************************************************************************")
	}

	dirExists, err := afero.DirExists(destFs, destFolder)

	if err != nil {

		return err
	}

	if dirExists {

		configFileName := nutsandbolts.GetConfigFileName()

		//---------------------------------------------------------------------
		// Collect files/directories to be deleted
		//---------------------------------------------------------------------

		var toBeDeleted []string

		err = afero.Walk(destFs, destFolder, func(path string, info fs.FileInfo, err error) error {

			if err != nil {

				return err
			}

			//-----------------------------------------------------------------
			// Delete all configuration files for security reasons
			//-----------------------------------------------------------------

			if !info.IsDir() && (filepath.Base(path) == configFileName) {

				toBeDeleted = append(toBeDeleted, path)

				return nil
			}

			//-----------------------------------------------------------------

			relativePath, err := filepath.Rel(destFolder, path)

			if err != nil {

				return err
			}

			_, err = sourceFs.Stat(relativePath)

			if err != nil {

				if os.IsNotExist(err) {

					toBeDeleted = append(toBeDeleted, path)
				} else {

					return err
				}
			}

			return nil
		})

		//---------------------------------------------------------------------

		for idx := len(toBeDeleted) - 1; idx >= 0; idx-- {

			if doDebug {

				fmt.Println("*** Delete from disk:", toBeDeleted[idx])
			}

			if err = destFs.Remove(toBeDeleted[idx]); err != nil {

				return err
			}
		}
	}

	return nil
}

//-----------------------------------------------------------------------------

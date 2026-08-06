package inputFolderWrite

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/ForTheTrashBin/LooLid/helper/nutsandbolts"
	"github.com/spf13/afero"
)

//-----------------------------------------------------------------------------
// Delete files/directories that are not (or no longer) present in the source
//-----------------------------------------------------------------------------

func cleanDir(ctx context.Context, destFs afero.Fs, destFolder string, sourceFs afero.Fs) error {

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
			// Check cancellation
			//-----------------------------------------------------------------

			select {

			case <-ctx.Done():

				return ctx.Err()

			default:
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

		for _, path := range slices.Backward(toBeDeleted) {

			//-----------------------------------------------------------------
			// Check cancellation
			//-----------------------------------------------------------------

			select {

			case <-ctx.Done():

				return ctx.Err()

			default:
			}

			//-----------------------------------------------------------------

			if err = destFs.Remove(path); err != nil {

				return err
			}
		}
	}

	return nil
}

//-----------------------------------------------------------------------------

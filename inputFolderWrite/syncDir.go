package inputFolderWrite

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ForTheTrashBin/LooLid/helper/nutsandbolts"
	"github.com/spf13/afero"
)

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

func syncDir(ctx context.Context, sourceFs afero.Fs, sourceFolder string, destFs afero.Fs, destFolder string, timeStamp time.Time, x int) error {

	//-------------------------------------------------------------------------
	// Check cancellation
	//-------------------------------------------------------------------------

	select {

	case <-ctx.Done():

		return ctx.Err()

	default:
	}

	//-------------------------------------------------------------------------
	// A 'local' function to calculate the MD5-hash of a file
	//-------------------------------------------------------------------------

	calculateMD5 := func(fileSys afero.Fs, path string) (string, error) {

		//---------------------------------------------------------------------
		// Check cancellation
		//---------------------------------------------------------------------

		select {

		case <-ctx.Done():

			return "", ctx.Err()

		default:
		}

		//---------------------------------------------------------------------

		file, err := fileSys.Open(path)

		if err != nil {

			return "", err
		}

		defer file.Close()

		//---------------------------------------------------------------------

		hashMD5 := md5.New() // MD5 works like a file you can write into

		if _, err := io.Copy(hashMD5, &nutsandbolts.ContextReader{Context: ctx, Reader: file}); err != nil {

			return "", err
		}

		return fmt.Sprintf("%x", hashMD5.Sum(nil)), nil // return hash as a string
	}

	//-------------------------------------------------------------------------

	sourceFileInfos, err := afero.ReadDir(sourceFs, sourceFolder)

	if err == nil {

		configFileName := nutsandbolts.GetConfigFileName()

		nutsandbolts.SortFileInfos(sourceFileInfos, true)

		for _, sourceFileInfo := range sourceFileInfos {

			//-----------------------------------------------------------------
			// Check cancellation
			//-----------------------------------------------------------------

			select {

			case <-ctx.Done():

				return ctx.Err()

			default:
			}

			//-----------------------------------------------------------------

			newSource := filepath.Join(sourceFolder, sourceFileInfo.Name())
			newDest := filepath.Join(destFolder, sourceFileInfo.Name())

			if sourceFileInfo.IsDir() {

				if err = destFs.MkdirAll(newDest, os.ModePerm); err != nil {

					return err
				}

				if err = syncDir(ctx, sourceFs, newSource, destFs, newDest, timeStamp, x+1); err != nil {

					return err
				}

				if err = destFs.Chmod(newDest, sourceFileInfo.Mode()); err != nil {

					return err
				}
			} else {

				if sourceFileInfo.Name() != configFileName { // Never copy configuration files

					destFileInfo, err := destFs.Stat(newDest)

					if err == nil {

						if destFileInfo.Size() == sourceFileInfo.Size() { // Are the files of same size

							ramHash, _ := calculateMD5(sourceFs, newSource)
							diskHash, _ := calculateMD5(destFs, newDest)

							if ramHash == diskHash { // Are the files of same content

								continue // Leave file AND timestamp untouched
							}
						}

						if err = copyFile(ctx, sourceFs, newSource, destFs, newDest); err != nil {

							return err
						}
					} else {

						if err = copyFile(ctx, sourceFs, newSource, destFs, newDest); err != nil {

							return err
						}

						/* TODO
						if err= destFs.Chtimes(destPath, timeStamp, timeStamp); err != nil{

							return err
						}

						if err = osspecific.PreserveOwner(memFs, path, diskFs, destinationPath, info); err != nil {

							return err
						}

						if err = PreserveTimes(info, diskFs, destinationPath); err != nil {

							return err
						}
						*/
					}
				}
			}
		}
	}

	return err
}

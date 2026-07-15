package inputFolderWrite

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/ForTheTrashBin/LooLid/helper/nutsandbolts"
	"github.com/spf13/afero"
)

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

func syncDir(sourceFs afero.Fs, sourceFolder string, destFs afero.Fs, destFolder string, timeStamp time.Time, x int) error {

	var doDebug bool = false

	if doDebug {

		fmt.Println("*********************************************************************************")
		fmt.Print("*** syncDir sourceFolder <", sourceFolder, "> destFolder <", destFolder, ">\n")
		fmt.Println("*********************************************************************************")
	}

	//-------------------------------------------------------------------------
	// A 'local' function to calculate the MD5-hash of a file
	//-------------------------------------------------------------------------

	calculateMD5 := func(fileSys afero.Fs, path string) (string, error) {

		file, err := fileSys.Open(path)

		if err != nil {

			return "", err
		}

		defer file.Close()

		//---------------------------------------------------------------------

		hashMD5 := md5.New() // MD5 works like a file you can write into

		if _, err := io.Copy(hashMD5, file); err != nil {

			return "", err
		}

		return fmt.Sprintf("%x", hashMD5.Sum(nil)), nil // return hash as a string
	}

	//-------------------------------------------------------------------------

	sourceFileInfos, err := afero.ReadDir(sourceFs, sourceFolder)

	if err == nil {

		sort.Slice(sourceFileInfos, func(i, j int) bool {

			if sourceFileInfos[i].IsDir() {

				if sourceFileInfos[j].IsDir() {

					return sourceFileInfos[i].Name() < sourceFileInfos[j].Name()
				} else {

					return true
				}
			} else {

				if sourceFileInfos[j].IsDir() {

					return false
				} else {

					return sourceFileInfos[i].Name() < sourceFileInfos[j].Name()
				}
			}
		})

		configFileName := nutsandbolts.GetConfigFileName()

		for idx, sourceFileInfo := range sourceFileInfos {

			if doDebug {

				fmt.Println("*** Index                :", idx)
				fmt.Println("*** sourceFolder         :", sourceFolder)
				fmt.Println("*** destFolder           :", destFolder)
				fmt.Println("*** sourceFileInfo.IsDir :", sourceFileInfo.IsDir())
				fmt.Println("*** sourceFileInfo.Name  :", sourceFileInfo.Name())
				fmt.Println("*** sourceFileInfo.Size  :", sourceFileInfo.Size())
				fmt.Println("*** sourceFileInfo.Mode  :", sourceFileInfo.Mode())
			}

			if sourceFileInfo.IsDir() {

				newSourcePath := filepath.Join(sourceFolder, sourceFileInfo.Name())
				newDestPath := filepath.Join(destFolder, sourceFileInfo.Name())

				if err = destFs.MkdirAll(newDestPath, os.ModePerm); err == nil { // ALL rights on directory

					if err = syncDir(sourceFs, newSourcePath, destFs, newDestPath, timeStamp, x+1); err == nil {

						if err = destFs.Chmod(newDestPath, sourceFileInfo.Mode()); err != nil {

							return err
						}
					} else {

						return err
					}
				} else {

					return err
				}
			} else {

				if sourceFileInfo.Name() != configFileName { // Never copy configuration files

					sourcePath := filepath.Join(sourceFolder, sourceFileInfo.Name())
					destPath := filepath.Join(destFolder, sourceFileInfo.Name())

					destFileInfo, err := destFs.Stat(destPath)

					if err == nil {

						if doDebug {

							fmt.Println("*** The file does exist")
						}

						if destFileInfo.Size() == sourceFileInfo.Size() { // Are the files of same size

							ramHash, _ := calculateMD5(sourceFs, sourcePath)
							diskHash, _ := calculateMD5(destFs, destPath)

							if ramHash == diskHash { // Are the files of same content

								if doDebug {

									fmt.Println("*** The file is ident")
								}

								continue // Leave file AND timestamp untouched
							}
						}

						if doDebug {

							fmt.Println("*** The file is NOT ident --> need to copy")
						}

						if err = copyFile(sourceFs, sourcePath, destFs, destPath); err != nil {

							return err
						}
					} else {

						if doDebug {

							fmt.Println("*** The file does NOT exist --> need to copy")
						}

						if err = copyFile(sourceFs, sourcePath, destFs, destPath); err != nil {

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

//-----------------------------------------------------------------------------

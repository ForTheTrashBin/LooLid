package inputFolderRead

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/afero"
)

func copyDir(sourceFs afero.Fs, sourcePath string, destFs afero.Fs, destPath string) error {

	var doDebug bool = false

	if doDebug {

		fmt.Println("*********************************************************************************")
		fmt.Print("*** copyDir source <", sourcePath, "> dest <", destPath, ">\n")
		fmt.Println("*********************************************************************************")
	}

	sourceFileInfos, err := afero.ReadDir(sourceFs, sourcePath)

	if err == nil {

		sort.Slice(sourceFileInfos, func(i, j int) bool {

			if sourceFileInfos[i].IsDir() {
				if sourceFileInfos[j].IsDir() {
					return sourceFileInfos[i].Name() < sourceFileInfos[j].Name()
				} else {
					return false
				}
			} else {
				if sourceFileInfos[j].IsDir() {
					return true
				} else {
					return sourceFileInfos[i].Name() < sourceFileInfos[j].Name()
				}
			}
		})

		/*
			for idx := range len(sourceFileInfos) - 1 {

				if doDebug {

					fmt.Println("*** Search for", configFile, "in", idx, sourceFileInfos[idx].Name(), sourceFileInfos[idx].IsDir())
				}

				if !sourceFileInfos[idx].IsDir() && sourceFileInfos[idx].Name() == "LooLid.config" {
					//Found

					break
				}
			}
		*/
		for idx, sourceFileInfo := range sourceFileInfos {

			if doDebug {

				fmt.Println("*** Index                :", idx)
				fmt.Println("*** sourcePath           :", sourcePath)
				fmt.Println("*** destPath             :", destPath)
				fmt.Println("*** sourceFileInfo.IsDir :", sourceFileInfo.IsDir())
				fmt.Println("*** sourceFileInfo.Name  :", sourceFileInfo.Name())
				fmt.Println("*** sourceFileInfo.Size  :", sourceFileInfo.Size())
				fmt.Println("*** sourceFileInfo.Mode  :", sourceFileInfo.Mode())
			}

			var doCopy bool = true

			if isInBlacklist(sourceFileInfo.Name(), sourceFileInfo.IsDir()) {

				doCopy = isInWhitelist(sourceFileInfo.Name(), sourceFileInfo.IsDir())
			}

			if doCopy {

				if sourceFileInfo.IsDir() {

					newSourcePath := filepath.Join(sourcePath, sourceFileInfo.Name())
					newDestPath := filepath.Join(destPath, sourceFileInfo.Name())

					if err = destFs.Mkdir(newDestPath, os.ModePerm); err == nil {

						if err = copyDir(sourceFs, newSourcePath, destFs, newDestPath); err == nil {

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

					if err = copyFile(
						sourceFs,
						filepath.Join(sourcePath, sourceFileInfo.Name()),
						destFs,
						filepath.Join(destPath, sourceFileInfo.Name())); err != nil {

						return err
					}
				}
			}
		}
	}

	return err
}

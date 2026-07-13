package inputFolderRead

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ForTheTrashBin/LooLid/helper/constants"
	"github.com/ForTheTrashBin/LooLid/helper/osspecific"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/afero"
	"github.com/theckman/yacspin"
)

type reader struct {
	localizer   *i18n.Localizer
	inputFolder string
	pureAppName string
	memFs       *afero.Fs
}

func newReader(localizer *i18n.Localizer, inputFolder string, pureAppName string, memFs *afero.Fs) *reader {

	rdr := &reader{

		localizer:   localizer,
		inputFolder: inputFolder,
		pureAppName: pureAppName,
		memFs:       memFs,
	}

	return rdr
}

func (rdr *reader) getLocalizedMessage(MessageID string, value1 string, value2 string) string {

	return rdr.localizer.MustLocalize(
		&i18n.LocalizeConfig{
			MessageID: MessageID,
			TemplateData: map[string]string{
				"value1": value1,
				"value2": value2,
			},
		})
}

func stat(fs afero.Fs, path string) (os.FileInfo, error) {

	if fs, ok := fs.(afero.Lstater); ok {

		fi, _, err := fs.LstatIfPossible(path)

		return fi, err
	}

	return fs.Stat(path)
}

func PreserveTimes(sourceInfo os.FileInfo, destFs afero.Fs, dest string) error {

	timeSpec := osspecific.GetTimeSpec(sourceInfo)

	return destFs.Chtimes(dest, timeSpec.TimeAccess, timeSpec.TimeModify)
}

func chmod(fs afero.Fs, dir string, mode os.FileMode, reported *error) {

	if err := fs.Chmod(dir, mode); *reported == nil {

		*reported = err
	}
}

func closeFile(f afero.File, reported *error) {

	if err := f.Close(); *reported == nil {

		*reported = err
	}
}

func isInBlacklist(fileName string, isDir bool) bool {

	// ignoreFiles = [ "(?i)\\.psd$", "(?i)\\.odp$", "(?i)\\.ppt$", "(?i)luftbild\\.jpg$" ]

	var bFound = false

	if !isDir {

		return strings.HasSuffix(fileName, ".psd") ||
			strings.HasSuffix(fileName, ".odp") ||
			strings.HasSuffix(fileName, ".ppt") || strings.ToLower(fileName) == "luftbild.jpg"

	}

	return bFound
}

func isInWhitelist(fileName string, isDir bool) bool {

	return false
}

func cleanDestination(diskFs afero.Fs, memFs afero.Fs, destinationPath string) error {

	var doDebug bool = false

	if doDebug {

		fmt.Println("*********************************************************************************")
		fmt.Println("**************************  cleanDestination  ***********************************")
		fmt.Println("*********************************************************************************")
	}

	exists, err := afero.DirExists(diskFs, destinationPath)

	if err != nil {

		return err
	}

	if exists {

		var toBeDeleted []string

		err = afero.Walk(diskFs, destinationPath, func(path string, info fs.FileInfo, err error) error {

			if err != nil {

				return err
			}

			relativePath, err := filepath.Rel(destinationPath, path)

			if err != nil {

				return err
			}

			_, err = memFs.Stat(relativePath)

			if err != nil {

				if os.IsNotExist(err) {

					toBeDeleted = append(toBeDeleted, path)
				} else {

					return err
				}
			}

			return nil
		})

		for idx := len(toBeDeleted) - 1; idx >= 0; idx-- {

			if doDebug {

				fmt.Println("Lösche von Datenträger", toBeDeleted[idx])
			}

			if err = diskFs.Remove(toBeDeleted[idx]); err != nil {

				return err
			}
		}
	}

	return nil
}

func calculateMD5(fsys afero.Fs, path string) (string, error) {

	f, err := fsys.Open(path)

	if err != nil {

		return "", err
	}

	defer f.Close()

	h := md5.New()

	if _, err := io.Copy(h, f); err != nil {

		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func syncDestination(diskFs afero.Fs, destinationFolder string, memFs afero.Fs, timeStamp time.Time) error {

	var doDebug bool = false

	if doDebug {

		fmt.Println("*********************************************************************************")
		fmt.Println("**************************  syncDestination  ************************************")
		fmt.Println("*********************************************************************************")
	}

	err := cleanDestination(diskFs, memFs, destinationFolder)

	if err != nil {

		return err
	}

	if doDebug {

		fmt.Println("*** cleanDestination Ok!")
		fmt.Println("---------------------------------------------------------------------------------")
	}

	return afero.Walk(memFs, "", func(path string, info fs.FileInfo, err error) error {

		if err != nil {

			return err
		}

		if path == "" {

			return nil
		}

		destinationPath := filepath.Join(destinationFolder, path)

		if doDebug {

			fmt.Println("*** destinationFolder:", destinationFolder)
			fmt.Println("*** path             :", path)
			fmt.Println("*** destinationPath  :", destinationPath)
			fmt.Println("*** info.IsDir       :", info.IsDir())
			fmt.Println("*** info.Name        :", info.Name())
			fmt.Println("*** info.Size        :", info.Size())
			fmt.Println("*** info.Mode        :", info.Mode())
		}

		if info.IsDir() {

			return diskFs.MkdirAll(destinationPath, info.Mode())
		}

		destFileInfo, err := diskFs.Stat(destinationPath)

		if err == nil { // The file does exist!!!

			if destFileInfo.Size() == info.Size() {

				ramHash, _ := calculateMD5(memFs, path)
				diskHash, _ := calculateMD5(diskFs, destinationPath)

				if ramHash == diskHash {

					// DATEI & ZEITSTEMPEL bleiebn unverändert

					return nil
				}
			}

			data, err := afero.ReadFile(memFs, path)

			if err != nil {

				return err
			}

			err = afero.WriteFile(diskFs, destinationPath, data, info.Mode())

			if err != nil {

				return nil
			}

			return diskFs.Chtimes(destinationPath, timeStamp, timeStamp)
		} else { // The file does NOT exist!!!

			err = diskFs.MkdirAll(filepath.Dir(destinationPath), os.ModePerm)

			if err != nil {

				return err
			}

			//-----------------------------------------------------------------
			// Create the destination file
			//-----------------------------------------------------------------

			hDestFile, err := diskFs.Create(destinationPath)

			if err != nil {

				return err
			}

			defer hDestFile.Close() // ensure the file ist closed

			//---------------------------------------------------------------------
			// Set the file's mode to the original
			//---------------------------------------------------------------------

			if err = diskFs.Chmod(destinationPath, info.Mode()); err != nil {

				return err
			}

			//---------------------------------------------------------------------

			if doDebug {

				if localInfo, err := diskFs.Stat(destinationPath); err != nil {

					panic(err)
				} else {

					fmt.Println("*** DestFile created:", destinationPath, "with", localInfo.Mode())
				}
			}

			//---------------------------------------------------------------------

			hSourceFile, err := memFs.Open(path)

			if err != nil {

				return err
			}

			defer hSourceFile.Close() // ensure the file ist closed

			if doDebug {

				if localInfo, err := memFs.Stat(path); err != nil {

					panic(err)
				} else {

					fmt.Println("*** SourceFile opened:", path, "with", localInfo.Mode())
				}
			}

			//---------------------------------------------------------------------
			// copy file content from source to dest
			//---------------------------------------------------------------------

			if _, err := io.CopyBuffer(hDestFile, hSourceFile, nil); err != nil {

				return err
			}

			if err = hDestFile.Sync(); err != nil {

				return nil
			}

			//---------------------------------------------------------------------
			/*
				if err = osspecific.PreserveOwner(memFs, path, diskFs, destinationPath, info); err != nil {

					return err
				}
				/*
					if err = PreserveTimes(info, diskFs, destinationPath); err != nil {

						return err
					}
			*/
			if doDebug {

				fmt.Println("*** Copy of filesuccessful")
				fmt.Println("*********************************************************************************")
			}
		}

		return err
	})
}

func (rdr *reader) MachMal() error {

	// var doDebug bool = false

	diskFs := afero.NewOsFs()

	err := copyDir(diskFs, rdr.inputFolder, *rdr.memFs, "")

	if err != nil {

		panic(err)
	}

	return nil
}

func InputFolderRead(localizer *i18n.Localizer, inputfolder string, pureAppName string, memFs *afero.MemMapFs) error {

	return nil // TODO:
}

func InputFolderReadAsync(localizer *i18n.Localizer, inputfolder string, pureAppName string, memFs *afero.Fs, sigCh chan os.Signal) error {

	batzen := newReader(localizer, inputfolder, pureAppName, memFs)

	spinnerSuffix := batzen.getLocalizedMessage(constants.SpinnerSuffixInputFolderRead, "", "")
	spinnerStopMessage := batzen.getLocalizedMessage(constants.SpinnerStopMessage, "", "")

	spinnerConfig := yacspin.Config{
		Frequency:         constants.Spinner_FrequencyMS * time.Millisecond,
		CharSet:           yacspin.CharSets[constants.Spinner_CharSet],
		Suffix:            " " + spinnerSuffix,
		SuffixAutoColon:   true,
		StopCharacter:     constants.Spinner_StopCharacter,
		StopColors:        []string{constants.Spinner_StopColor},
		StopFailCharacter: constants.Spinner_StopFailCharacter,
		StopFailColors:    []string{constants.Spinner_StopFailColor},
		StopMessage:       spinnerStopMessage,
	}

	spinner, err := yacspin.New(spinnerConfig)

	if err != nil {
		panic(fmt.Errorf("spinner init failed: %w", err))
	}

	spinner.Reverse()

	if err := spinner.Start(); err != nil {
		panic(fmt.Errorf("spinner start failed: %w", err))
	}

	defer spinner.Stop()

	//-------------------------------------------------------------------------

	doneChannel := make(chan error, 1)

	go func() {

		doneChannel <- batzen.MachMal()
	}()

	//-------------------------------------------------------------------------

	select {

	case <-sigCh:

		stopFailMessage := batzen.getLocalizedMessage(constants.CheckError_AbortedByUser, "", "")

		spinner.StopFailMessage(stopFailMessage)

		spinner.StopFail()

		return constants.ErrInterrupted

	case err := <-doneChannel:

		if err != nil {

			if err == constants.ErrInputFolderNotCorrect {

				stopFailMessage := batzen.getLocalizedMessage(constants.CheckError_InputfolderIncorrect, "", "")

				spinner.StopFailMessage(stopFailMessage)

				spinner.StopFail()

				// batzen.reportInputFolderErrors()

				return err
			}

			if err == constants.ErrBulkDataNotCorrect {

				stopFailMessage := batzen.getLocalizedMessage(constants.CheckError_BulkdataIncorrect, "", "")

				spinner.StopFailMessage(stopFailMessage)

				spinner.StopFail()

				// batzen.reportBulkDataErrors()

				return err
			}

			spinner.StopFailMessage(err.Error())

			spinner.StopFail()

			// batzen.reportBulkDataErrors()

			return err
		}

		return err
	}
}

package inputFolderRead

import (
	"fmt"
	"io"

	"github.com/ForTheTrashBin/LooLid/helper/osspecific"
	"github.com/spf13/afero"
)

func copyFile(sourceFs afero.Fs, sourceName string, destFs afero.Fs, destName string) error {

	var doDebug bool = false

	if doDebug {

		fmt.Println("*********************************************************************************")
		fmt.Print("*** copyFile source <", sourceName, "> dest <", destName, ">\n")
		fmt.Println("*********************************************************************************")
	}

	//-------------------------------------------------------------------------
	// Open the source file to be copied
	//-------------------------------------------------------------------------

	sourceFile, err := sourceFs.Open(sourceName)

	if err != nil {

		return err
	}

	defer sourceFile.Close() // ensure the file ist closed on exit

	if doDebug {

		fmt.Println("*** sourceFile opened    :", sourceName)
	}

	//-------------------------------------------------------------------------
	// Get the 'fileinfo' of the sourceFile just opened
	//-------------------------------------------------------------------------

	sourceFileInfo, err := sourceFile.Stat()

	if err != nil {

		return err
	}

	if doDebug {

		fmt.Println("*** sourceFileInfo.IsDir :", sourceFileInfo.IsDir())
		fmt.Println("*** sourceFileInfo.Name  :", sourceFileInfo.Name())
		fmt.Println("*** sourceFileInfo.Size  :", sourceFileInfo.Size())
		fmt.Println("*** sourceFileInfo.Mode  :", sourceFileInfo.Mode())
		fmt.Println("---------------------------------------------------------------------------------")
	}

	//-------------------------------------------------------------------------
	// Create the destination file (in memory)
	//-------------------------------------------------------------------------

	destFile, err := destFs.Create(destName)

	if err != nil {

		return err
	}

	defer destFile.Close() // ensure the file ist closed on exit

	//-------------------------------------------------------------------------
	// Set the destFile's mode to the original (mode of sourceFile)
	//-------------------------------------------------------------------------

	if err = destFs.Chmod(destName, sourceFileInfo.Mode()); err != nil {

		return err
	}

	//-------------------------------------------------------------------------

	if doDebug {

		fmt.Println("*** destFile created     :", destName)

		if destFileInfo, err := destFile.Stat(); err == nil {

			fmt.Println("*** destFileInfo.IsDir   :", destFileInfo.IsDir())
			fmt.Println("*** destFileInfo.Name    :", destFileInfo.Name())
			fmt.Println("*** destFileInfo.Size    :", destFileInfo.Size())
			fmt.Println("*** destFileInfo.Mode    :", destFileInfo.Mode())
			fmt.Println("---------------------------------------------------------------------------------")
		} else {

			return err
		}
	}

	//-------------------------------------------------------------------------
	// copy file content from source to dest
	//-------------------------------------------------------------------------

	if _, err = io.CopyBuffer(destFile, sourceFile, nil); err != nil {

		return err
	}

	if err = destFile.Sync(); err != nil {

		return nil
	}

	//---------------------------------------------------------------------

	if err = osspecific.PreserveOwner(sourceFs, sourceName, destFs, destName, sourceFileInfo); err != nil {

		return err
	}

	if err = PreserveTimes(sourceFileInfo, destFs, destName); err != nil {

		return err
	}

	if doDebug {

		fmt.Println("*** copyFile of filesuccessful")
		fmt.Println("*********************************************************************************")
	}

	return nil
}

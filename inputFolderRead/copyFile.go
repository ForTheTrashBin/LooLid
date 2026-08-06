package inputFolderRead

import (
	"context"
	"io"

	"github.com/ForTheTrashBin/LooLid/helper/nutsandbolts"
	"github.com/ForTheTrashBin/LooLid/helper/osspecific"
	"github.com/spf13/afero"
)

//-----------------------------------------------------------------------------

func copyFile(ctx context.Context, sourceFs afero.Fs, sourceName string, destFs afero.Fs, destName string) error {

	//-------------------------------------------------------------------------
	// inner function to manage defered operations properly
	//-------------------------------------------------------------------------

	innerCopyFile := func() error {

		//---------------------------------------------------------------------
		// Check cancellation
		//---------------------------------------------------------------------

		select {

		case <-ctx.Done():

			return ctx.Err()

		default:
		}

		//---------------------------------------------------------------------
		// Open the source file to be copied
		//---------------------------------------------------------------------

		sourceFile, err := sourceFs.Open(sourceName)

		if err != nil {

			return err
		}

		defer sourceFile.Close() // ensure the file is closed on exit

		//---------------------------------------------------------------------
		// Create the destination file
		//---------------------------------------------------------------------

		destFile, err := destFs.Create(destName)

		if err != nil {

			return err
		}

		defer destFile.Close() // ensure the file is closed on exit

		//---------------------------------------------------------------------
		// Check cancellation
		//---------------------------------------------------------------------

		select {

		case <-ctx.Done():

			return ctx.Err()

		default:
		}

		//---------------------------------------------------------------------
		// copy file content from source to dest and do it context-aware!
		//---------------------------------------------------------------------

		if _, err = io.CopyBuffer(destFile, &nutsandbolts.ContextReader{Context: ctx, Reader: sourceFile}, nil); err != nil {

			return err
		}

		if err = destFile.Sync(); err != nil {

			return err
		}

		return nil
	}

	//-------------------------------------------------------------------------
	// Check cancellation
	//-------------------------------------------------------------------------

	select {

	case <-ctx.Done():

		return ctx.Err()

	default:
	}

	//-------------------------------------------------------------------------

	if err := innerCopyFile(); err != nil {

		return err
	}

	//-------------------------------------------------------------------------
	// Get the 'fileinfo' of the sourceFile just opened
	//-------------------------------------------------------------------------

	sourceFileInfo, err := sourceFs.Stat(sourceName)

	if err != nil {

		return err
	}

	//-------------------------------------------------------------------------
	// Set the destFile's mode to the original (Windows: AFTER the file is closed)
	//-------------------------------------------------------------------------

	if err = destFs.Chmod(destName, sourceFileInfo.Mode()); err != nil {

		return err
	}

	//-------------------------------------------------------------------------

	if err = osspecific.PreserveOwner(sourceFs, sourceName, destFs, destName, sourceFileInfo); err != nil {

		return err
	}

	if err = PreserveTimes(sourceFileInfo, destFs, destName); err != nil {

		return err
	}

	return nil
}

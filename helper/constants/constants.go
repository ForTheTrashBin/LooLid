package constants

import "errors"

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

const (
	AppConfig_ConfigFileExtension   = ".config"
	AppConfig_DefaultInputDirectory = "content"

	Spinner_FrequencyMS       = 125
	Spinner_CharSet           = 11
	Spinner_StopCharacter     = "✓"
	Spinner_StopColor         = "fgGreen"
	Spinner_StopFailCharacter = "✖"
	Spinner_StopFailColor     = "fgRed"
)

var ErrInterrupted = errors.New("interrupted by user")
var ErrInputFolderNotCorrect = errors.New("err_inputfolder")
var ErrBulkDataNotCorrect = errors.New("err_bulkdata")

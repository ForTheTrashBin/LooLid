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

const (
	CheckError_NoFileInfo              = "check.no_fileinfo"
	CheckError_NoDirectory             = "check.no_directory"
	CheckError_DotHiddenDirectory      = "check.dot_hidden_directory"
	CheckError_DirectoryHiddenOrSystem = "check.directory_hidden_or_system"
	CheckError_Confusion               = "check.directory_confusion"
	CheckError_DirectoryNoPermission   = "check.directory_no_permission"
	CheckError_DirectoryNotReadable    = "check.directory_not_readable"
	CheckError_NoPathAbs               = "check.no_pathabs"
	CheckError_NoSymLinks              = "check.no_symlinks"
	CheckError_NoWorkingDir            = "check.no_workingdir"
	CheckError_NoPathRel               = "check.no_pathrel"
	CheckError_WorkingDirInInput       = "check.workingdir_in_input"

	CheckError_DirectorySiblingsFound = "check.header.directory_siblings_found"
	CheckError_InvalidWindowsChar     = "check.header.invalid_windows_chars"
	CheckError_TrailingDotSpace       = "check.header.trailing_dot_or_space"
	CheckError_ReservedWindowsName    = "check.header.reserved_windows_name"
	CheckError_FileNameTooLong        = "check.header.filename_too_long"
	CheckError_PathNameTooLong        = "check.header.pathname_too_long"
	CheckError_UnicodeCollision       = "check.header.unicode_collision"
	CheckError_SymbolicLinkDetected   = "check.header.symbolic_link"
	CheckError_ConfigFile             = "check.header.config_file"
	CheckError_DuplicateEntries       = "check.header.duplicate_entries"
	CheckError_InputfolderIncorrect   = "check.header.inputfolder_incorrect"
	CheckError_BulkdataIncorrect      = "check.header.bulkdata_incorrect"
	CheckError_AbortedByUser          = "check.header.aborted_by_user"

	SpinnerSuffixInputfolderCheck = "spinner.suffix.inputfolder_check"
	SpinnerSuffixInputFolderRead  = "spinner.suffix.inputfolder_read"
	SpinnerSuffixXxxxx            = "spinner.suffix.xxxxx"
	SpinnerSuffixYyyyy            = "spinner.suffix.yyyyy"
	SpinnerSuffixZzzzz            = "spinner.suffix.zzzzz"
	SpinnerStopMessage            = "spinner.stopmessage"

	Label_Folder = "label.folder"
	Label_File   = "label.file"
)

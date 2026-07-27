package constants

import "errors"

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

const (
	AppConfig_ConfigFileExtension = ".toml"

	// AppConfig_DefaultInputFolder = "content"
	// AppConfig_DefaultOutputFolder = "target"
	// AppConfig_DefaultTemplatesFolder = "templates"

	AppConfig_DefaultInputFolder     = "/home/u32800/Documents/WebSite/content"
	AppConfig_DefaultOutputFolder    = "/home/u32800/Documents/WebSite/target"
	AppConfig_DefaultTemplatesFolder = "/home/u32800/Documents/WebSite/templates"

	Spinner_FrequencyMS       = 125
	Spinner_CharSet           = 11
	Spinner_StopCharacter     = "✓"
	Spinner_StopColor         = "fgGreen"
	Spinner_StopFailCharacter = "✖"
	Spinner_StopFailColor     = "fgRed"
)

var ErrInterrupted = errors.New("interrupted by user")

var ErrInputFolderIncorrect = errors.New("err_inputfolderincorrect")
var ErrTemplateFolderIncorrect = errors.New("err_templatefolderincorrect")
var ErrBulkDataIncorrect = errors.New("err_bulkdataincorrect")

var ErrRecoveredPanicWithoutType = errors.New("err_recovered_panic_without_type")

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

	CheckError_DirectorySiblingsFound  = "check.header.directory_siblings_found"
	CheckError_InvalidWindowsChar      = "check.header.invalid_windows_chars"
	CheckError_TrailingDotSpace        = "check.header.trailing_dot_or_space"
	CheckError_ReservedWindowsName     = "check.header.reserved_windows_name"
	CheckError_FileNameTooLong         = "check.header.filename_too_long"
	CheckError_PathNameTooLong         = "check.header.pathname_too_long"
	CheckError_UnicodeCollision        = "check.header.unicode_collision"
	CheckError_SymbolicLinkDetected    = "check.header.symbolic_link"
	CheckError_ConfigFile              = "check.header.config_file"
	CheckError_DuplicateEntries        = "check.header.duplicate_entries"
	CheckError_SubfoldersFound         = "check.header.subfolder_found"
	CheckError_NoExtensions            = "check.header.no_extension"
	CheckError_DuplicateExtension      = "check.header.duplicate_extension"
	CheckError_InputfolderIncorrect    = "check.header.inputfolder_incorrect"
	CheckError_TemplatefolderIncorrect = "check.header.templatefolder_incorrect"
	CheckError_BulkdataIncorrect       = "check.header.bulkdata_incorrect"
	CheckError_AbortedByUser           = "check.header.aborted_by_user"

	BlackWhiteError_ReadConfig       = "blackwhite.read_config"
	BlackWhiteError_Blacklist        = "blackwhite.blacklist"
	BlackWhiteError_Whitelist        = "blackwhite.whitelist"
	BlackWhiteError_InvalidScope     = "blackwhite.invalid_scope"
	BlackWhiteError_EmptyPattern     = "blackwhite.empty_pattern"
	BlackWhiteError_EmptyGlobPattern = "blackwhite.empty_glob_pattern"
	BlackWhiteError_InvalidRegex     = "blackwhite.invalid_regex"
	BlackWhiteError_EmptyExactValue  = "blackwhite.empty_exact_value"
	BlackWhiteError_EmptySuffix      = "blackwhite.empty_suffix"
	BlackWhiteError_EmptyPrefix      = "blackwhite.empty_prefix"

	SpinnerSuffixInputfolderCheck    = "spinner.suffix.inputfolder_check"
	SpinnerSuffixInputFolderRead     = "spinner.suffix.inputfolder_read"
	SpinnerSuffixInputFolderProccess = "spinner.suffix.inputfolder_proccess"
	SpinnerSuffixInputFolderWrite    = "spinner.suffix.inputfolder_write"
	SpinnerSuffixZzzzz               = "spinner.suffix.zzzzz"
	SpinnerStopMessageDone           = "spinner.stopmessage.done"
	SpinnerStopMessageError          = "spinner.stopmessage.error"

	Label_Folder    = "label.folder"
	Label_File      = "label.file"
	Label_Extension = "label.extension"
)

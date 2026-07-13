package inputFolderCheck

//-----------------------------------------------------------------------------

type Entry struct {
	filePath  string
	entryName string
	isDir     bool
	isSymlink bool
}

//-----------------------------------------------------------------------------

type DuplicateItem struct {
	itemName string
	isDir    bool
}

type DuplicateGroup struct {
	groupName string
	items     []DuplicateItem
}

//-----------------------------------------------------------------------------

type Issue struct {
	isDir    bool
	filePath string
	info     string
}

//-----------------------------------------------------------------------------

type ErrorType string

const (
	CheckError_NoFileInfo              ErrorType = "check.no_fileinfo"
	CheckError_NoDirectory             ErrorType = "check.no_directory"
	CheckError_DotHiddenDirectory      ErrorType = "check.dot_hidden_directory"
	CheckError_DirectoryHiddenOrSystem ErrorType = "check.directory_hidden_or_system"
	CheckError_Confusion               ErrorType = "check.directory_confusion"
	CheckError_DirectoryNoPermission   ErrorType = "check.directory_no_permission"
	CheckError_DirectoryNotReadable    ErrorType = "check.directory_not_readable"
	CheckError_NoPathAbs               ErrorType = "check.no_pathabs"
	CheckError_NoSymLinks              ErrorType = "check.no_symlinks"
	CheckError_NoWorkingDir            ErrorType = "check.no_workingdir"
	CheckError_NoPathRel               ErrorType = "check.no_pathrel"
	CheckError_WorkingDirInInput       ErrorType = "check.workingdir_in_input"

	CheckError_DirectorySiblingsFound ErrorType = "check.header.directory_siblings_found"
	CheckError_InvalidWindowsChar     ErrorType = "check.header.invalid_windows_chars"
	CheckError_TrailingDotSpace       ErrorType = "check.header.trailing_dot_or_space"
	CheckError_ReservedWindowsName    ErrorType = "check.header.reserved_windows_name"
	CheckError_FileNameTooLong        ErrorType = "check.header.filename_too_long"
	CheckError_PathNameTooLong        ErrorType = "check.header.pathname_too_long"
	CheckError_UnicodeCollision       ErrorType = "check.header.unicode_collision"
	CheckError_SymbolicLinkDetected   ErrorType = "check.header.symbolic_link"
	CheckError_ConfigFile             ErrorType = "check.header.config_file"
	CheckError_DuplicateEntries       ErrorType = "check.header.duplicate_entries"
	CheckError_InputfolderIncorrect   ErrorType = "check.header.inputfolder_incorrect"
	CheckError_BulkdataIncorrect      ErrorType = "check.header.bulkdata_incorrect"
	CheckError_AbortedByUser          ErrorType = "check.header.aborted_by_user"

	SpinnerSuffixCheckInputfolder ErrorType = "spinner.suffix.check_inputfolder"
	SpinnerSuffixReadInputFolder  ErrorType = "spinner.suffix.read_inputfolder"
	SpinnerSuffixXxxxx            ErrorType = "spinner.suffix.xxxxx"
	SpinnerSuffixYyyyy            ErrorType = "spinner.suffix.yyyyy"
	SpinnerSuffixZzzzz            ErrorType = "spinner.suffix.zzzzz"
	SpinnerStopMessage            ErrorType = "spinner.stopmessage"

	Label_Folder ErrorType = "label.folder"
	Label_File   ErrorType = "label.file"
)

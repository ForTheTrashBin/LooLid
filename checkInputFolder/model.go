package checkInputFolder

//-----------------------------------------------------------------------------

type Entry struct {
	Path string
	Name string

	IsDir     bool
	IsSymlink bool
}

//-----------------------------------------------------------------------------

type ErrorType string

const (
	CheckError_NoFileInfo              ErrorType = "check.no_fileinfo"
	CheckError_NoDirectory             ErrorType = "check.no_directory"
	CheckError_DotHiddenDirectory      ErrorType = "check.dot_hidden_directory"
	CheckError_DirectoryHiddenOrSystem ErrorType = "check.directory_hidden_or_system"
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

	xCaseCollision ErrorType = "check.case_collision"
)

type Issue struct {
	Path string
	Info string
}

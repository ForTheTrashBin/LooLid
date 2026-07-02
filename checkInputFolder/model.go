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
	NoFileInfo              ErrorType = "check.no_fileinfo"
	NoDirectory             ErrorType = "check.no_directory"
	DotHiddenDirectory      ErrorType = "check.dot_hidden_directory"
	DirectoryHiddenOrSystem ErrorType = "check.directory_hidden_or_system"
	DirectoryNoPermission   ErrorType = "check.directory_no_permission"
	DirectoryNotReadable    ErrorType = "check.directory_not_readable"
	NoPathAbs               ErrorType = "check.no_pathabs"
	NoSymLinks              ErrorType = "check.no_symlinks"
	NoWorkingDir            ErrorType = "check.no_workingdir"
	NoPathRel               ErrorType = "check.no_pathrel"
	WorkingDirInInput       ErrorType = "check.workingdir_in_input"

	DirectorySiblingsFound ErrorType = "check.header.directory_siblings_found"
	InvalidWindowsChar     ErrorType = "check.header.invalid_windows_chars"
	TrailingDotSpace       ErrorType = "check.header.trailing_dot_or_space"
	ReservedWindowsName    ErrorType = "check.header.reserved_windows_name"
	NameTooLong            ErrorType = "check.header.name_too_long"
	PathTooLong            ErrorType = "check.header.path_too_long"
	UnicodeCollision       ErrorType = "check.header.unicode_collision"
	SymbolicLinkDetected   ErrorType = "check.header.symbolic_link"

	CaseCollision ErrorType = "check.case_collision"
)

type Issue struct {
	Type ErrorType
	Path string
	Info string
}

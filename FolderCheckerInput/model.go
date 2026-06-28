package FolderCheckerInput

//-----------------------------------------------------------------------------

type Statistics struct {
	Files       int
	Directories int
}

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
	NoFileInfo              ErrorType = "no_fileinfo"
	NoDirectory             ErrorType = "no_directory"
	DirectoryNotReadable    ErrorType = "directory_not_readable"
	DirectoryNoPermission   ErrorType = "girectory_no_permission"
	DirectoryHiddenOrSystem ErrorType = "directory_hidden_or_system"
	DotHiddenDirectory      ErrorType = "dot_hidden_directory"
	NoPathAbs               ErrorType = "no_pathabs"
	NoPathRel               ErrorType = "no_pathrel"
	NoSymLinks              ErrorType = "no_symlinks"
	NoWorkingDir            ErrorType = "no_workingdir"
	WorkingDirInInput       ErrorType = "workingdir_in_input"
	CaseCollision           ErrorType = "case_collision"
	UnicodeCollision        ErrorType = "unicode_collision"
	InvalidWindowsChar      ErrorType = "invalid_windows_character"
	ReservedWindowsName     ErrorType = "reserved_windows_name"
	TrailingDotSpace        ErrorType = "trailing_dot_or_space"
	PathTooLong             ErrorType = "path_too_long"
	NameTooLong             ErrorType = "name_too_long"
	SymbolicLinkDetected    ErrorType = "symbolic_link"
)

type Issue struct {
	Type ErrorType
	Path string
	Info string
}

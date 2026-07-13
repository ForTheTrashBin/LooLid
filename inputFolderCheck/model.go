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

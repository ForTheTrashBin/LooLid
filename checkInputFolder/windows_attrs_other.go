//go:build !windows

package FolderCheckerInput

func isHiddenOrSystemOnWindows(path string) bool {
	return false
}

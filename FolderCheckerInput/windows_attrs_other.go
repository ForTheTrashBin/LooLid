//go:build !windows

package FolderCheckerInput

func hasHiddenOrSystemOnWindows(path string) bool {
	return false
}

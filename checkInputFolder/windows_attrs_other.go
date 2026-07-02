//go:build !windows

package checkInputFolder

func isHiddenOrSystemOnWindows(path string) bool {
	return false
}

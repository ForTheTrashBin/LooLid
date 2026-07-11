//go:build !windows

package osspecific

func IsHiddenOrSystemOnWindows(path string) bool {
	return false
}

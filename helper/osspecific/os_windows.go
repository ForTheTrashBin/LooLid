//go:build windows

package osspecific

import (
	"errors"
	"os"
	"syscall"
	"time"

	"github.com/spf13/afero"
)

//-----------------------------------------------------------------------------

func fileInfoStat(v interface{}) *syscall.Win32FileAttributeData {

	s, ok := v.(*syscall.Win32FileAttributeData)

	if !ok {

		panic(errors.New("WindowsCast: not a *syscall.Win32FileAttributeData"))
	}

	return s
}

//-----------------------------------------------------------------------------

func IsHiddenOrSystem(path string) bool {

	ptr, err := syscall.UTF16PtrFromString(path)

	if err != nil {

		return false
	}

	attr, err := syscall.GetFileAttributes(ptr)

	if err != nil {

		return false
	}

	const (
		HIDDEN = 0x2
		SYSTEM = 0x4
	)

	return ((attr&HIDDEN != 0) || (attr&SYSTEM != 0))
}

//-----------------------------------------------------------------------------

func PreserveOwner(sourceFs afero.Fs, source string, destFs afero.Fs, dest string, info os.FileInfo) (err error) {

	return nil
}

//-----------------------------------------------------------------------------

func GetTimeSpec(info os.FileInfo) TimeSpec {

	stat := fileInfoStat(info.Sys())

	return TimeSpec{

		TimeModify: time.Unix(0, stat.LastWriteTime.Nanoseconds()),
		TimeAccess: time.Unix(0, stat.LastAccessTime.Nanoseconds()),
		TimeCreate: time.Unix(0, stat.CreationTime.Nanoseconds()),
	}
}

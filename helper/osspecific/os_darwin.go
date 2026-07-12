//go:build darwin

package osspecific

import (
	"errors"
	"os"
	"syscall"
	"time"

	"github.com/spf13/afero"
)

//-----------------------------------------------------------------------------

func fileInfoStat(v interface{}) *syscall.Stat_t {

	s, ok := v.(*syscall.Stat_t)

	if !ok {

		panic(errors.New("DarwinCast: not a *syscall.Stat_t"))
	}

	return s
}

//-----------------------------------------------------------------------------

func IsHiddenOrSystem(path string) bool {

	return false
}

//-----------------------------------------------------------------------------

func PreserveOwner(sourceFs afero.Fs, source string, destFs afero.Fs, dest string, info os.FileInfo) (err error) {

	if info == nil {

		if info, err = sourceFs.Stat(source); err != nil {

			return err
		}
	}

	if stat, ok := info.Sys().(*syscall.Stat_t); ok {

		if err := destFs.Chown(dest, int(stat.Uid), int(stat.Gid)); err != nil {

			return err
		}
	}

	return nil
}

//-----------------------------------------------------------------------------

func GetTimeSpec(info os.FileInfo) TimeSpec {

	stat := fileInfoStat(info.Sys())

	return TimeSpec{

		TimeModify: info.ModTime(),
		TimeAccess: time.Unix(stat.Atimespec.Sec, stat.Atimespec.Nsec),
		TimeCreate: time.Unix(stat.Ctimespec.Sec, stat.Ctimespec.Nsec),
	}
}

//-----------------------------------------------------------------------------

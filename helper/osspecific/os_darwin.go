//go:build darwin

package osspecific

import (
	"errors"
	"os"
	"syscall"
	"time"

	"github.com/spf13/afero"
	"golang.org/x/sys/unix"
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
// Set the terminal echo on or off.
//
// If the teminal is redirected "IoctlGetTermios" will return "unix.ENOTTY".
// In that case, the caller should ignore the error and continue. This is
// useful for example when running in a docker container
//-----------------------------------------------------------------------------

func SetEcho(enable bool) error {

	fd := int(os.Stdin.Fd())

	termios, err := unix.IoctlGetTermios(fd, unix.TIOCGETA)

	if err != nil {

		return err
	}

	if enable {

		termios.Lflag |= unix.ECHO

	} else {

		termios.Lflag &^= unix.ECHO
	}

	return unix.IoctlSetTermios(fd, unix.TIOCSETA, termios)
}

func FlushStdin() error {

	return unix.IoctlSetPointerInt(int(os.Stdin.Fd()), unix.TIOCFLUSH, 1)
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

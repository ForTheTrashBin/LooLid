//go:build darwin && !ios && cgo

//-----------------------------------------------------------------------------
// THIS FILE IS TAKEN FROM "https://github.com/jeandeaual/go-locale"
//-----------------------------------------------------------------------------

package locales

// Non-CGO implementation is in locale_darwin_nocgo.go

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation

#include <AppKit/AppKit.h>

const char * preferredLocalization();
const char * preferredLocalizations();
*/
import "C"
import (
	"strings"
)

// GetLocale retrieves the IETF BCP 47 language tag set on the system.
func GetLocale() (string, error) {
	str := C.preferredLocalization()
	if output := C.GoString(str); output != "" {
		return strings.Replace(output, "_", "-", 1), nil
	}

	return getLocaleCli()
}

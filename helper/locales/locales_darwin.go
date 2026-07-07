//go:build darwin && !ios

//-----------------------------------------------------------------------------
// THIS FILE IS TAKEN FROM "https://github.com/jeandeaual/go-locale"
//-----------------------------------------------------------------------------

package locales

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

// Common code between locale_darwin_cgo.go and locale_darwin_nocgo.go
// GetLocale exported functions are in those files

func execCommand(cmd string, args ...string) (status int, out string, err error) {
	var bytesOut []byte
	status = -1
	command := exec.Command(cmd, args...)

	// Execute the command and get the standard and error outputs
	bytesOut, err = command.CombinedOutput()
	out = string(bytesOut)
	if err != nil {
		return
	}

	// Check the status code
	if w, ok := command.ProcessState.Sys().(syscall.WaitStatus); ok {
		status = w.ExitStatus()
	}

	return
}

// getLocaleCli retrieves the IETF BCP 47 language tag set on the system without using CGO to call OS APIs.
func getLocaleCli() (string, error) {
	_, output, err := execCommand("defaults", "read", "-g", "AppleLocale")
	if err != nil {
		return "", fmt.Errorf("cannot determine locale: %v (output: %s)", err, output)
	}

	// defaults read -g AppleLocale can return a string containing additional
	// information after the locale, e.g. "en_US@currency=USD"
	if idx := strings.Index(output, "@"); idx != -1 {
		output = output[:idx]
	}

	return strings.TrimRight(strings.Replace(output, "_", "-", 1), "\n"), nil
}

// GetLanguage retrieves the IETF BCP 47 language tag set on the system and
// returns the language part of the tag.
func GetLanguage() (string, error) {
	language := ""

	locale, err := GetLocale()
	if err == nil {
		language, _ = splitLocale(locale)
	}

	return language, err
}

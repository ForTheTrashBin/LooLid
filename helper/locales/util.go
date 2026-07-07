//go:build !android

//-----------------------------------------------------------------------------
// THIS FILE IS TAKEN FROM "https://github.com/jeandeaual/go-locale"
//-----------------------------------------------------------------------------

package locales

import (
	"strings"
)

func splitLocale(locale string) (string, string) {
	// Remove the encoding, if present.
	formattedLocale, _, _ := strings.Cut(locale, ".")

	// Normalize by replacing the hyphens with underscores
	formattedLocale = strings.ReplaceAll(formattedLocale, "-", "_")

	// Split at the underscore.
	language, territory, _ := strings.Cut(formattedLocale, "_")
	return language, territory
}

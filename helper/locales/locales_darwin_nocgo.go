//go:build darwin && !ios && !cgo

//-----------------------------------------------------------------------------
// THIS FILE IS TAKEN FROM "https://github.com/jeandeaual/go-locale"
//-----------------------------------------------------------------------------

package locales

// GetLocale retrieves the IETF BCP 47 language tag set on the system.
func GetLocale() (string, error) {
	return getLocaleCli()
}

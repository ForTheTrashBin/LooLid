package checkInputFolder

import (
	"strings"
	"unicode/utf16"

	"golang.org/x/text/unicode/norm"
)

//-----------------------------------------------------------------------------
// Windows does not allow these characters in file and directory names.
//-----------------------------------------------------------------------------

var invalidWindowsCharacters = []rune{'<', '>', ':', '"', '/', '\\', '|', '?', '*'}

func (chk *checker) checkInvalidWindowsCharacters(entry Entry) {

	for _, character := range entry.Name {

		for _, invalid := range invalidWindowsCharacters {

			if character == invalid {

				chk.issuesInvalidWindowsCharacters = append(
					chk.issuesInvalidWindowsCharacters, Issue{
						Type: InvalidWindowsChar,
						Path: entry.Path,
						Info: "invalid Windows character: " + string(character)})
			}
		}
	}
}

//-----------------------------------------------------------------------------
// Windows does not allow filenames that end with a period or a space
//-----------------------------------------------------------------------------

func (chk *checker) checkTrailingCharacters(entry Entry) {

	if len(entry.Name) != 0 {

		last := entry.Name[len(entry.Name)-1]

		if last == '.' || last == ' ' {

			chk.issuesTrailingCharacters = append(
				chk.issuesTrailingCharacters, Issue{
					Type: TrailingDotSpace,
					Path: entry.Path,
					Info: "name ends with dot or space"})
		}
	}
}

//-----------------------------------------------------------------------------
// Windows does not allow filenames that might conflict with device names
//-----------------------------------------------------------------------------

var reservedNames = map[string]bool{

	"CON": true,
	"PRN": true,
	"AUX": true,
	"NUL": true,

	"COM1": true,
	"COM2": true,
	"COM3": true,
	"COM4": true,
	"COM5": true,
	"COM6": true,
	"COM7": true,
	"COM8": true,
	"COM9": true,

	"LPT1": true,
	"LPT2": true,
	"LPT3": true,
	"LPT4": true,
	"LPT5": true,
	"LPT6": true,
	"LPT7": true,
	"LPT8": true,
	"LPT9": true,
}

func (chk *checker) checkReservedWindowsName(entry Entry) {

	name := entry.Name

	if index := strings.IndexByte(name, '.'); index >= 0 {
		name = name[:index]
	}

	name = strings.ToUpper(name)

	if reservedNames[name] {

		chk.issuesReservedWindowsName = append(
			chk.issuesReservedWindowsName, Issue{
				Type: ReservedWindowsName,
				Path: entry.Path,
				Info: "reserved Windows filename"})
	}
}

//-----------------------------------------------------------------------------
// Windows does not allow filenames longer than 255 UTF-16 units
//-----------------------------------------------------------------------------

func utf16Length(
	value string,
) int {

	return len(
		utf16.Encode(
			[]rune(value),
		),
	)
}

func (chk *checker) checkLength(entry Entry) {

	if utf16Length(entry.Name) > 255 {

		chk.issuesNameLength = append(
			chk.issuesNameLength,
			Issue{
				Type: NameTooLong,
				Path: entry.Name,
				Info: "filename exceeds 255 UTF-16 units",
			},
		)
	}

	if utf16Length(entry.Path) > 240 {

		chk.issuesPathLength = append(
			chk.issuesPathLength,
			Issue{
				Type: PathTooLong,
				Path: entry.Path,
				Info: "path exceeds 240 UTF-16 units",
			},
		)
	}
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

func (chk *checker) checkUnicodeNormalization(entry Entry) {

	if norm.NFC.String(entry.Path) != entry.Path {

		chk.issuesUnicodeNormalization = append(
			chk.issuesUnicodeNormalization,
			Issue{

				Type: UnicodeCollision,
				Path: entry.Path,
				Info: "path is not NFC normalized",
			},
		)
	}
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

func (chk *checker) checkSymLink(entry Entry) {

	if entry.IsSymlink {
		chk.issuesSymLink = append(
			chk.issuesSymLink,
			Issue{

				Type: SymbolicLinkDetected,
				Path: entry.Path,
				Info: "symbolic link detected",
			},
		)
	}
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

func (chk *checker) validateEntries(entries []Entry) {

	for _, entry := range entries {

		chk.checkInvalidWindowsCharacters(entry)
		chk.checkTrailingCharacters(entry)
		chk.checkReservedWindowsName(entry)
		chk.checkLength(entry)
		chk.checkUnicodeNormalization(entry)
		chk.checkSymLink(entry)
	}
}

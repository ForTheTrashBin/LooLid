package FolderCheckerInput

import (
	"strings"
	"unicode/utf16"

	"golang.org/x/text/unicode/norm"
)

//-----------------------------------------------------------------------------
// Windows does not allow these characters in file and directory names.
//-----------------------------------------------------------------------------

var invalidWindowsCharacters = []rune{
	'<',
	'>',
	':',
	'"',
	'/',
	'\\',
	'|',
	'?',
	'*',
}

func checkInvalidWindowsCharacters(entry Entry) []Issue {

	var issues []Issue

	for _, character := range entry.Name {

		for _, invalid := range invalidWindowsCharacters {

			if character == invalid {

				issues = append(
					issues,
					Issue{
						Type: InvalidWindowsChar,
						Path: entry.Path,
						Info: "invalid Windows character: " +
							string(character),
					},
				)

				return issues
			}
		}
	}

	return issues
}

//-----------------------------------------------------------------------------
// Windows does not allow filenames that end with a period or a space
//-----------------------------------------------------------------------------

func checkTrailingCharacters(entry Entry) []Issue {

	if len(entry.Name) == 0 {
		return nil
	}

	last := entry.Name[len(entry.Name)-1]

	if last == '.' || last == ' ' {

		return []Issue{
			{
				Type: TrailingDotSpace,
				Path: entry.Path,
				Info: "name ends with dot or space",
			},
		}
	}

	return nil
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

func checkReservedWindowsName(entry Entry) []Issue {

	name := entry.Name

	if index := strings.IndexByte(name, '.'); index >= 0 {
		name = name[:index]
	}

	name = strings.ToUpper(name)

	if reservedNames[name] {

		return []Issue{
			{
				Type: ReservedWindowsName,
				Path: entry.Path,
				Info: "reserved Windows filename",
			},
		}
	}

	return nil
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

func checkLength(entry Entry) []Issue {

	var issues []Issue

	if utf16Length(entry.Name) > 255 {

		issues = append(
			issues,
			Issue{
				Type: NameTooLong,
				Path: entry.Path,
				Info: "filename exceeds 255 UTF-16 units",
			},
		)
	}

	if utf16Length(entry.Path) > 240 {

		issues = append(
			issues,
			Issue{
				Type: PathTooLong,
				Path: entry.Path,
				Info: "path exceeds 240 UTF-16 units",
			},
		)
	}

	return issues
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

func checkUnicodeNormalization(entry Entry) []Issue {

	if norm.NFC.String(entry.Path) != entry.Path {

		return []Issue{
			{
				Type: UnicodeCollision,
				Path: entry.Path,
				Info: "path is not NFC normalized",
			},
		}
	}

	return nil
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

func validateEntries(entries []Entry) []Issue {

	var issues []Issue

	for _, entry := range entries {

		issues = append(
			issues,
			checkInvalidWindowsCharacters(entry)...,
		)

		issues = append(
			issues,
			checkTrailingCharacters(entry)...,
		)

		issues = append(
			issues,
			checkReservedWindowsName(entry)...,
		)

		issues = append(
			issues,
			checkLength(entry)...,
		)

		issues = append(
			issues,
			checkUnicodeNormalization(entry)...,
		)

	}

	return issues
}

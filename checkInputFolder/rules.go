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

func (chk *checker) checkInvalidWindowsChar(entry Entry) {

	for _, character := range entry.entryName {

		for _, invalid := range invalidWindowsCharacters {

			if character == invalid {

				chk.issuesInvalidWindowsChar = append(
					chk.issuesInvalidWindowsChar, Issue{
						isDir:    entry.isDir,
						filePath: entry.filePath,
						info:     string(character),
					})
			}
		}
	}
}

//-----------------------------------------------------------------------------
// Windows does not allow filenames that end with a period or a space
//-----------------------------------------------------------------------------

func (chk *checker) checkTrailingDotSpace(entry Entry) {

	if len(entry.entryName) != 0 {

		last := entry.entryName[len(entry.entryName)-1]

		if last == '.' || last == ' ' {

			chk.issuesTrailingDotSpace = append(
				chk.issuesTrailingDotSpace, Issue{
					isDir:    entry.isDir,
					filePath: entry.filePath,
				})
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

	name := entry.entryName

	if index := strings.IndexByte(name, '.'); index >= 0 {
		name = name[:index]
	}

	name = strings.ToUpper(name)

	if reservedNames[name] {

		chk.issuesReservedWindowsName = append(
			chk.issuesReservedWindowsName, Issue{
				isDir:    entry.isDir,
				filePath: entry.filePath,
			})
	}
}

//-----------------------------------------------------------------------------
// Windows does not allow filenames-length longer than 255 UTF-16 units
// Windows does not allow pathnames-length longer than 240 UTF-16 units
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

func (chk *checker) checkFileNameAndPathNameLength(entry Entry) {

	if utf16Length(entry.entryName) > 255 {

		chk.issuesFileNameLength = append(
			chk.issuesFileNameLength,
			Issue{
				isDir:    entry.isDir,
				filePath: entry.entryName,
			},
		)
	}

	if utf16Length(entry.filePath) > 240 {

		chk.issuesPathNameLength = append(
			chk.issuesPathNameLength,
			Issue{
				isDir:    entry.isDir,
				filePath: entry.filePath,
			},
		)
	}
}

//-----------------------------------------------------------------------------
// Check UTF16-normalisation for MacOs
//-----------------------------------------------------------------------------

func (chk *checker) checkUnicodeNormalization(entry Entry) {

	if norm.NFC.String(entry.filePath) != entry.filePath {

		chk.issuesUnicodeNormalization = append(
			chk.issuesUnicodeNormalization,
			Issue{
				isDir:    entry.isDir,
				filePath: entry.filePath,
			},
		)
	}
}

//-----------------------------------------------------------------------------
// Do not allow any symbolic links
//-----------------------------------------------------------------------------

func (chk *checker) checkSymLink(entry Entry) {

	if entry.isSymlink {
		chk.issuesSymLink = append(
			chk.issuesSymLink,
			Issue{
				isDir:    entry.isDir,
				filePath: entry.filePath,
			},
		)
	}
}

//-----------------------------------------------------------------------------
// Check all entries with all methods
//-----------------------------------------------------------------------------

func (chk *checker) validateEntries(entries []Entry) {

	for _, entry := range entries {

		chk.checkInvalidWindowsChar(entry)
		chk.checkTrailingDotSpace(entry)
		chk.checkReservedWindowsName(entry)
		chk.checkFileNameAndPathNameLength(entry)
		chk.checkUnicodeNormalization(entry)
		chk.checkSymLink(entry)
	}
}

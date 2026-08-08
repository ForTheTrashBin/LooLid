package inputFolderCheck

import (
	"strings"
	"unicode/utf16"

	"golang.org/x/text/unicode/norm"
)

//-----------------------------------------------------------------------------
// Windows does not allow these characters in file and directory names.
//-----------------------------------------------------------------------------

var invalidWindowsCharacters = []rune{'<', '>', ':', '"', '/', '\\', '|', '?', '*'}

func (chk *checker) checkInvalidWindowsChar(filePath string, entryName string, isDir bool, isSymlink bool) {

	for _, character := range entryName {

		for _, invalid := range invalidWindowsCharacters {

			if character == invalid {

				chk.issuesInvalidWindowsChar = append(chk.issuesInvalidWindowsChar,

					Issue{

						isDir:    isDir,
						filePath: filePath,
						info:     string(character),
					})
			}
		}
	}
}

//-----------------------------------------------------------------------------
// Windows does not allow filenames that end with a period or a space
//-----------------------------------------------------------------------------

func (chk *checker) checkTrailingDotSpace(filePath string, entryName string, isDir bool, isSymlink bool) {

	if len(entryName) != 0 {

		last := entryName[len(entryName)-1]

		if last == '.' || last == ' ' {

			chk.issuesTrailingDotSpace = append(chk.issuesTrailingDotSpace,

				Issue{

					isDir:    isDir,
					filePath: filePath,
				})
		}
	}
}

//-----------------------------------------------------------------------------
// Windows does not allow filenames that might conflict with device names
//-----------------------------------------------------------------------------

var reservedNames = map[string]bool{

	"CON": true, "PRN": true, "AUX": true, "NUL": true,

	"COM1": true, "COM2": true, "COM3": true,
	"COM4": true, "COM5": true, "COM6": true,
	"COM7": true, "COM8": true, "COM9": true,

	"LPT1": true, "LPT2": true, "LPT3": true,
	"LPT4": true, "LPT5": true, "LPT6": true,
	"LPT7": true, "LPT8": true, "LPT9": true,
}

func (chk *checker) checkReservedWindowsName(filePath string, entryName string, isDir bool, isSymlink bool) {

	name := entryName

	if index := strings.IndexByte(name, '.'); index >= 0 {

		name = name[:index]
	}

	name = strings.ToUpper(name)

	if reservedNames[name] {

		chk.issuesReservedWindowsName = append(chk.issuesReservedWindowsName,

			Issue{

				isDir:    isDir,
				filePath: filePath,
			})
	}
}

//-----------------------------------------------------------------------------
// Windows does not allow filenames-length longer than 255 UTF-16 units
// Windows does not allow pathnames-length longer than 240 UTF-16 units
//-----------------------------------------------------------------------------

func utf16Length(value string) int {

	return len(utf16.Encode([]rune(value)))
}

func (chk *checker) checkFileNameAndPathNameLength(filePath string, entryName string, isDir bool, isSymlink bool) {

	if utf16Length(entryName) > 255 {

		chk.issuesFileNameLength = append(chk.issuesFileNameLength,

			Issue{

				isDir:    isDir,
				filePath: entryName,
			})
	}

	if utf16Length(filePath) > 240 {

		chk.issuesPathNameLength = append(chk.issuesPathNameLength,

			Issue{

				isDir:    isDir,
				filePath: filePath,
			})
	}
}

//-----------------------------------------------------------------------------
// Check UTF16-normalisation for MacOs
//-----------------------------------------------------------------------------

func (chk *checker) checkUnicodeNormalization(filePath string, entryName string, isDir bool, isSymlink bool) {

	if norm.NFC.String(filePath) != filePath {

		chk.issuesUnicodeNormalization = append(chk.issuesUnicodeNormalization,

			Issue{

				isDir:    isDir,
				filePath: filePath,
			})
	}
}

//-----------------------------------------------------------------------------
// Do not allow any symbolic links
//-----------------------------------------------------------------------------

func (chk *checker) checkSymLink(filePath string, entryName string, isDir bool, isSymlink bool) {

	if isSymlink {

		chk.issuesSymLink = append(chk.issuesSymLink,

			Issue{

				isDir:    isDir,
				filePath: filePath,
			})
	}
}

//-----------------------------------------------------------------------------
// Check all entries with all methods
//-----------------------------------------------------------------------------

func (chk *checker) validateEntry(filePath string, entryName string, isDir bool, isSymlink bool) {

	chk.checkInvalidWindowsChar(filePath, entryName, isDir, isSymlink)
	chk.checkTrailingDotSpace(filePath, entryName, isDir, isSymlink)
	chk.checkReservedWindowsName(filePath, entryName, isDir, isSymlink)
	chk.checkFileNameAndPathNameLength(filePath, entryName, isDir, isSymlink)
	chk.checkUnicodeNormalization(filePath, entryName, isDir, isSymlink)
	chk.checkSymLink(filePath, entryName, isDir, isSymlink)
}

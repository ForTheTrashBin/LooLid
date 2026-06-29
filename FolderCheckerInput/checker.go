package FolderCheckerInput

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Check(inputPath string) ([]Issue, Statistics, error) {

	var stats Statistics

	//-------------------------------------------------------------------------
	// get fileinfo of given input-path
	//-------------------------------------------------------------------------

	fileInfo, err := os.Stat(inputPath)

	if err != nil {
		return []Issue{{Type: NoFileInfo, Path: inputPath, Info: err.Error()}}, stats, nil
	}

	//-------------------------------------------------------------------------
	// Check if input-path points to a directory
	//-------------------------------------------------------------------------

	if !fileInfo.IsDir() {
		return []Issue{{Type: NoDirectory, Path: inputPath}}, stats, nil
	}

	//-------------------------------------------------------------------------
	// Check, if directory-name starts with a dot (.)
	//-------------------------------------------------------------------------

	baseName := filepath.Base(inputPath)

	// Do not treat the "." and ".." directories as hidden directories.

	if (baseName != ".") && (baseName != "..") {
		if strings.HasPrefix(baseName, ".") {
			return []Issue{{Type: DotHiddenDirectory, Path: inputPath}}, stats, nil
		}
	}

	//-------------------------------------------------------------------------
	// On windows directory must not be hidden or system
	//-------------------------------------------------------------------------

	if hasHiddenOrSystemOnWindows(inputPath) {
		return []Issue{{Type: DirectoryHiddenOrSystem, Path: inputPath}}, stats, nil
	}

	//-------------------------------------------------------------------------
	// Check, if directory is readable (Just try and check the result)
	//-------------------------------------------------------------------------

	_, err = os.ReadDir(inputPath)

	if err != nil {
		if os.IsPermission(err) {
			return []Issue{{Type: DirectoryNoPermission, Path: inputPath}}, stats, nil
		} else {
			return []Issue{{Type: DirectoryNotReadable, Path: inputPath, Info: err.Error()}}, stats, nil
		}
	}

	//-------------------------------------------------------------------------
	// Make input-path absolute and evaluate sym-links
	//-------------------------------------------------------------------------

	inputPathAbs, err := filepath.Abs(inputPath)

	if err != nil {
		return []Issue{{Type: NoPathAbs, Path: inputPath, Info: err.Error()}}, stats, nil
	}

	inputPathAbs, err = filepath.EvalSymlinks(inputPathAbs)

	if err != nil {
		return []Issue{{Type: NoSymLinks, Path: inputPath, Info: err.Error()}}, stats, nil
	}

	//-------------------------------------------------------------------------
	// Get working-directory (which is an absolute path) and evaluate sym-links
	//-------------------------------------------------------------------------

	workingDirAbs, err := os.Getwd()

	if err != nil {
		return []Issue{{Type: NoWorkingDir, Info: err.Error()}}, stats, nil
	}

	workingDirAbs, err = filepath.EvalSymlinks(workingDirAbs)

	if err != nil {
		return []Issue{{Type: NoSymLinks, Path: workingDirAbs, Info: err.Error()}}, stats, nil
	}

	//-------------------------------------------------------------------------
	// 'Compare' and evaluate if workingDirAbs is inside inputPathAbs
	//-------------------------------------------------------------------------

	relativePath, err := filepath.Rel(inputPathAbs, workingDirAbs)

	if err != nil {
		return []Issue{{Type: NoPathRel, Path: inputPath, Info: err.Error()}}, stats, nil
	}

	if (relativePath != "..") && ((len(relativePath) < 3) || (relativePath[:3] != ".."+string(filepath.Separator))) {
		return []Issue{{Type: WorkingDirInInput, Path: inputPath}}, stats, nil
	}

	//-------------------------------------------------------------------------
	// Try to read siblings to detect directories and files with ambiguous names
	//-------------------------------------------------------------------------

	parentPath := filepath.Dir(inputPathAbs)

	dirEntries, err := os.ReadDir(parentPath)

	if err != nil {
		return []Issue{{Type: DirectoryNotReadable, Path: inputPath, Info: err.Error()}}, stats, nil
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			if dirEntry.Name() != baseName { // it's NOT me!
				if strings.ToLower(dirEntry.Name()) == strings.ToLower(baseName) {
					fmt.Println("Directory-Sibling found: ", dirEntry.Name())
				}
			}
		} else {
			if strings.Contains(strings.ToLower(dirEntry.Name()), strings.ToLower(baseName)) {
				fmt.Println("File-Sibling found: ", dirEntry.Name())
			}
		}
	}

	//-------------------------------------------------------------------------

	entries, scanIssues, stats, err := scanDirectory(inputPath)

	if err != nil {
		return nil, stats, err
	}

	var issues []Issue

	issues = append(issues, scanIssues...)

	issues = append(
		issues,
		validateEntries(entries)...,
	)

	issues = append(
		issues,
		analyzeCollisions(entries)...,
	)

	return issues, stats, nil
}

package helper

import (
	"io/fs"
	"os"
	"path/filepath"
)

type WalkDirFunc func(depth int, path string, d fs.DirEntry, err error) error

func walkDir(depth int, path string, d fs.DirEntry, walkDirFn WalkDirFunc) error {
	if err := walkDirFn(depth, path, d, nil); err != nil || !d.IsDir() {
		if err == filepath.SkipDir && d.IsDir() {
			// Successfully skipped directory.
			err = nil
		}
		return err
	}

	dirs, err := os.ReadDir(path)
	if err != nil {
		// Second call, to report ReadDir error.
		err = walkDirFn(depth+1, path, d, err)
		if err != nil {
			if err == filepath.SkipDir && d.IsDir() {
				err = nil
			}
			return err
		}
	}

	for _, d1 := range dirs {
		path1 := filepath.Join(path, d1.Name())
		if err := walkDir(depth+1, path1, d1, walkDirFn); err != nil {
			if err == filepath.SkipDir {
				break
			}
			return err
		}
	}
	return nil
}

func WalkDir(root string, fn WalkDirFunc) error {
	info, err := os.Lstat(root)
	if err != nil {
		err = fn(0, root, nil, err)
	} else {
		err = walkDir(0, root, fs.FileInfoToDirEntry(info), fn)
	}
	if err == filepath.SkipDir || err == filepath.SkipAll {
		return nil
	}
	return err
}

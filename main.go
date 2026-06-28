package main

import (
	"LooLid/FolderCheckerInput"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:")
		fmt.Println("  portable-checker <directory>")
		os.Exit(1)
	}

	root := os.Args[1]
	/*
		fmt.Println("root       : ", root)

		rootAbs, err := filepath.Abs(root)

		if err != nil {

			panic(err)
		}

		fmt.Println("rootAbs    : ", rootAbs)

		rootAbs, err = filepath.EvalSymlinks(rootAbs)

		if err != nil {

			panic(err)
		}

		fmt.Println("rootAbsSym : ", rootAbs)

		// Aktuelles Arbeitsverzeichnis bestimmen
		cwd, err := os.Getwd()

		if err != nil {
			panic(err)
		}

		fmt.Println("cwd        : ", cwd)

		cwd, err = filepath.EvalSymlinks(cwd)

		if err != nil {
			panic(err)
		}

		fmt.Println("cwdSym     : ", cwd)

		rel, err := filepath.Rel(rootAbs, cwd)

		if err != nil {
			panic(err)
		}

		fmt.Println("rel        : ", rel)

		if rel == ".." || len(rel) >= 3 && rel[:3] == ".."+string(filepath.Separator) {
			fmt.Println("WD liegt NICHT drin!")
		} else {
			fmt.Println("WD liegt drin!")
		}
	*/
	//--------------------------------------------

	issues, stats, err := FolderCheckerInput.Check(root)

	if err != nil {
		fmt.Println("Fatal:", err)
		os.Exit(2)
	}

	FolderCheckerInput.Print(issues, stats)

	if len(issues) > 0 {
		os.Exit(1)
	}
}

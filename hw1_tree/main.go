package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func buildTree(out io.Writer, path string, printFiles bool, tabString string, depth int, lastFile bool) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
	}

	var bufSlice []os.FileInfo

	if printFiles == false {
		for _, file := range files {
			if file.IsDir() {
				bufSlice = append(bufSlice, file)
			}
		}
		files = bufSlice
	}

	if depth != 0 {
		if !lastFile {
			tabString += "│" + "\t"
		} else {
			tabString += "\t"
		}
	}
	for idx, file := range files {
		if idx != len(files)-1 {
			_, err := fmt.Fprintf(out, "%s├───%s", tabString, file.Name())
			if err != nil {
				return err
			}
			lastFile = false
		} else {
			_, err := fmt.Fprintf(out, "%s└───%s", tabString, file.Name())
			if err != nil {
				return err
			}
			lastFile = true
		}

		if !file.IsDir() {
			if file.Size() != 0 {
				_, err := fmt.Fprintf(out, " (%db)", file.Size())
				if err != nil {
					return err
				}
			} else {
				_, err := fmt.Fprintf(out, " (empty)")
				if err != nil {
					return err
				}
			}
			_, err = fmt.Fprintf(out, "\n")
			if err != nil {
				return err
			}

		} else {
			_, err = fmt.Fprintf(out, "\n")
			if err != nil {
				return err
			}
			err := buildTree(out, path+string(os.PathSeparator)+file.Name(), printFiles, tabString, depth+1, lastFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return buildTree(out, path, printFiles, "", 0, false)
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

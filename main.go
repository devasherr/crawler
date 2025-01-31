package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

func traverse(folderPath, target, searchType string, avoidFilesMap map[string]bool) string {
	directories, err := os.ReadDir(folderPath)
	if err != nil {
		return ""
	}

	result := ""
	for _, dir := range directories {
		nextDir := filepath.Join(folderPath, dir.Name())

		if dir.IsDir() {
			if avoidFilesMap[dir.Name()] {
				return result
			}
			if dir.Name() == target {
				if searchType == "any" || searchType == "folder" {
					return nextDir
				}
			}
			curRes := traverse(nextDir, target, searchType, avoidFilesMap)
			if len(curRes) > 1 {
				result = curRes
			}
		} else {
			if dir.Name() == target && !avoidFilesMap[dir.Name()] {
				if searchType == "any" || searchType == "file" {
					return nextDir
				}
			}
		}
	}

	return result
}

func main() {
	var searchType string
	var avoidFiles []string
	pflag.StringVarP(&searchType, "type", "t", "any", "operate on folder or file")
	pflag.StringSliceVarP(&avoidFiles, "exclude", "x", []string{""}, "folder or file to exclude from the search")
	pflag.Parse()

	var avoidFilesMap = make(map[string]bool)
	for _, file := range avoidFiles {
		avoidFilesMap[file] = true
	}

	if len(pflag.Args()) == 0 {
		fmt.Println("no target provided")
		return
	}

	if len(pflag.Args()) == 0 {
		fmt.Println("too many targets provided")
		return
	}

	folderPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Printf(">> %s", traverse(folderPath, pflag.Args()[0], searchType, avoidFilesMap))
}

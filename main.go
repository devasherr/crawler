package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

type searchInfo struct {
}

func traverse(folderPath, target, searchType string, prefix bool, avoidFilesMap map[string]bool, prefixMatch *[]string) string {
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
			if searchType == "any" || searchType == "folder" {
				if dir.Name() == target {
					return nextDir
				} else if prefix && checkPrefix(target, dir.Name()) {
					*prefixMatch = append(*prefixMatch, nextDir)
				}
			}

			curRes := traverse(nextDir, target, searchType, prefix, avoidFilesMap, prefixMatch)
			if len(curRes) > 1 {
				result = curRes
			}
		} else {
			if avoidFilesMap[dir.Name()] {
				return result
			}

			if searchType == "any" || searchType == "file" {
				if dir.Name() == target {
					return nextDir
				} else if prefix && checkPrefix(target, dir.Name()) {
					*prefixMatch = append(*prefixMatch, nextDir)
				}
			}

			if dir.Name() == target && !avoidFilesMap[dir.Name()] {
				if searchType == "any" || searchType == "file" {
					return nextDir
				}
			}
		}
	}

	return result
}

func checkPrefix(word, target string) bool {
	if len(word) > len(target) {
		return false
	}
	for i := 0; i < len(word); i++ {
		if word[i] != target[i] {
			return false
		}
	}
	return true
}

func main() {
	var searchType string
	var avoidFiles []string
	var matchPrefix bool
	pflag.StringVarP(&searchType, "type", "t", "any", "operate on folder or file")
	pflag.StringSliceVarP(&avoidFiles, "exclude", "x", []string{""}, "folder or file to exclude from the search")
	pflag.BoolVarP(&matchPrefix, "prefix", "p", false, "return closest matching prefix if target not found")
	pflag.Parse()

	var avoidFilesMap = make(map[string]bool)
	for _, file := range avoidFiles {
		avoidFilesMap[file] = true
	}

	if len(pflag.Args()) == 0 {
		fmt.Println("no target provided")
		return
	}

	if len(pflag.Args()) > 1 {
		fmt.Println("too many targets provided")
		return
	}

	folderPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	prefixMatch := []string{}
	result := traverse(folderPath, pflag.Args()[0], searchType, matchPrefix, avoidFilesMap, &prefixMatch)
	if result == "" {
		fmt.Println(">> No match found !!")

		if matchPrefix {
			if len(prefixMatch) == 0 {
				fmt.Println(">> >> No prefix found")
			} else {
				fmt.Println("Closest matching prefixes found: ")
				for i, prefix := range prefixMatch {
					fmt.Printf("%d) %s\n", i+1, prefix)
				}
			}
		}
	}
}

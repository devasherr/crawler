package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
)

type Config struct {
	target      string
	searchType  string
	avoidFiles  map[string]bool
	matchPrefix bool
}

func traverse(folderPath string, config Config, prefixMatch *[]string) string {
	directories, err := os.ReadDir(folderPath)
	if err != nil {
		return ""
	}

	result := ""
	for _, dir := range directories {
		nextDir := filepath.Join(folderPath, dir.Name())

		if dir.IsDir() {
			if config.avoidFiles[dir.Name()] {
				continue
			}
			if config.searchType == "any" || config.searchType == "folder" {
				if dir.Name() == config.target {
					return nextDir
				} else if config.matchPrefix && checkPrefix(dir.Name(), config.target) {
					*prefixMatch = append(*prefixMatch, nextDir)
				}
			}

			curRes := traverse(nextDir, config, prefixMatch)
			if len(curRes) > 1 {
				result = curRes
			}
		} else {
			if config.avoidFiles[dir.Name()] {
				return result
			}

			if config.searchType == "any" || config.searchType == "file" {
				if dir.Name() == config.target {
					return nextDir
				} else if config.matchPrefix && checkPrefix(dir.Name(), config.target) {
					*prefixMatch = append(*prefixMatch, nextDir)
				}
			}

			if dir.Name() == config.target && !config.avoidFiles[dir.Name()] {
				if config.searchType == "any" || config.searchType == "file" {
					return nextDir
				}
			}
		}
	}

	return result
}

func checkPrefix(word, target string) bool {
	return strings.HasPrefix(word, target)
}

func parseFlag() Config {
	c := Config{}
	var avoidFiles []string
	pflag.StringVarP(&c.searchType, "type", "t", "any", "operate on folder or file")
	pflag.StringSliceVarP(&avoidFiles, "exclude", "x", []string{""}, "folder or file to exclude from the search")
	pflag.BoolVarP(&c.matchPrefix, "prefix", "p", false, "return closest matching prefix if target not found")
	pflag.Parse()

	c.avoidFiles = make(map[string]bool)
	for _, file := range avoidFiles {
		c.avoidFiles[file] = true
	}

	if len(pflag.Args()) == 0 {
		panic("no target provided")
	}
	if len(pflag.Args()) > 1 {
		panic("too many targets provided")
	}
	c.target = pflag.Args()[0]

	return c
}

func main() {
	config := parseFlag()

	folderPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	prefixMatch := []string{}
	result := traverse(folderPath, config, &prefixMatch)
	if result == "" {
		fmt.Println(">> No match found !!")

		if config.matchPrefix {
			if len(prefixMatch) == 0 {
				fmt.Println(">> >> No prefix found")
			} else {
				fmt.Println("Closest matching prefixes found: ")
				for i, prefix := range prefixMatch {
					fmt.Printf("%d) %s\n", i+1, prefix)
				}
			}
		}
	} else {
		fmt.Println(">> ", result)
	}
}

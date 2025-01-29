package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

func traverse(folderPath, target string) string {
	directories, err := os.ReadDir(folderPath)
	if err != nil {
		return ""
	}

	result := ""
	for _, dir := range directories {
		nextDir := filepath.Join(folderPath, dir.Name())

		if dir.IsDir() {
			if dir.Name() == target {
				return nextDir
			}
			curRes := traverse(nextDir, target)
			if len(curRes) > 1 {
				result = curRes
			}
		} else {
			if dir.Name() == target {
				return nextDir
			}
		}
	}

	return result
}

func main() {
	var flatType string
	pflag.StringVarP(&flatType, "type", "t", "", "operate on folder or file")
	pflag.Parse()

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

	fmt.Println(traverse(folderPath, pflag.Args()[0]))
}

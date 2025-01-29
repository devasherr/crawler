package main

import (
	"fmt"
	"os"
	"path/filepath"
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
	folderName, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Println(traverse(folderName, "abs"))
}

package fileTools

import (
	"os"
	"path/filepath"
)

//TODO: Can this be made concurrent??
var fileList []string

func visit(path string, f os.FileInfo, err error) error {
	if f.IsDir() {
		fileList = append(fileList, path)
	}
	return nil
}

//FindFiles searches through a folder for all files matching a pattern
func FindFiles(root, pattern string) []string {
	fileList = []string{}
	filepath.Walk(root, visit)
	var matchingFiles []string
	for _, path := range fileList {
		files, _ := filepath.Glob(path + "/" + pattern)
		for _, file := range files {
			matchingFiles = append(matchingFiles, file)
		}
	}
	return matchingFiles
}

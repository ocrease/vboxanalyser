package file

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type Explorer struct {
}

func (fe *Explorer) GetDirectory(path string) (dirs []Directory) {
	var subDirs []string
	if len(path) == 0 {
		subDirs = getRootPaths()
	} else {
		subDirs = subdirectoriesOf(path)
	}
	for _, v := range subDirs {
		dirs = append(dirs, Directory{Parent: path, Name: v, SubDirectories: subdirectoriesOf(path + v + "/")})
	}
	return dirs
}

func subdirectoriesOf(path string) []string {
	dirs := make([]string, 0)
	fmt.Printf("Looking for directories under %v\n", path)
	files, err := ioutil.ReadDir(filepath.Clean(path))
	if err != nil {
		fmt.Printf("Error occured %v\n", err)
		return dirs
	}
	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, f.Name())
		}
	}
	return dirs
}

package file

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func bitsToDrives(bitMap uint32) (drives []string) {
	availableDrives := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	for i := range availableDrives {
		if bitMap&1 == 1 {
			drives = append(drives, availableDrives[i]+":")
		}
		bitMap >>= 1
	}

	return
}

type Explorer struct {
}

// func (fe *Explorer) GetDirectoryContents(path string) ([]FileDescriptor, error) {
// 	var dirs []FileDescriptor
// 	if len(path) == 0 {
// 		for _, v := range getRootPaths() {
// 			dirs = append(dirs, FileDescriptor{Path: v, IsDir: true})
// 		}
// 	} else {
// 		files, err := ioutil.ReadDir(path)
// 		if err != nil {
// 			return nil, err
// 		}

// 		for _, f := range files {
// 			if f.IsDir() {
// 				dirs = append(dirs, FileDescriptor{BasePath: path, Path: f.Name(), IsDir: true})
// 			}
// 		}
// 	}
// 	return dirs, nil
// }

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

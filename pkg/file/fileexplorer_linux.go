// +build linux darwin

package file

func getRootPaths() (drives []string) {
	drives = append(drives, subdirectoriesOf("/")...)
	return
}

// +build windows

package file

import (
	"runtime"
	"syscall"
)

func getRootPaths() (drives []string) {
	switch runtime.GOOS {
	case "windows":
		kernel32, _ := syscall.LoadLibrary("kernel32.dll")
		getLogicalDrivesHandle, _ := syscall.GetProcAddress(kernel32, "GetLogicalDrives")

		if ret, _, callErr := syscall.Syscall(uintptr(getLogicalDrivesHandle), 0, 0, 0, 0); callErr != 0 {
			drives = append(drives, "C:")
		} else {
			drives = bitsToDrives(uint32(ret))
		}
	default:
		drives = append(drives, subdirectoriesOf("/")...)
	}
	return
}

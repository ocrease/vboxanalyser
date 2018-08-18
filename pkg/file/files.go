package file

type Service interface {
	GetDirectory(dirPath string) []Directory
}

type Directory struct {
	Parent         string   `json:"parent"`
	Name           string   `json:"name"`
	SubDirectories []string `json:"subdirectories"`
}

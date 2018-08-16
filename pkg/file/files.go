package file

type FileDescriptor struct {
	BasePath string `json:"basePath"`
	Path     string `json:"path"`
	IsDir    bool   `json:"isDir"`
}

type Service interface {
	GetDirectory(dirPath string) []Directory
}

type Directory struct {
	Parent         string   `json:"parent"`
	Name           string   `json:"name"`
	SubDirectories []string `json:"subdirectories"`
}

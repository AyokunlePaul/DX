package file

import (
	"mime/multipart"
)

var (
	ReqTypes = []string{"errand", "profile"}
)

func New() *File {
	return new(File)
}

func NewListRequest(folder string, headers []*multipart.FileHeader) []*File {
	files := make([]*File, len(headers))
	for index, header := range headers {
		currentFile := new(File)

		currentFile.Folder = folder
		currentFile.Header = header
		files[index] = currentFile
	}
	return files
}

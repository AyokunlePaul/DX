package file

import (
	"mime/multipart"
)

type Type string

type File struct {
	Folder      string
	Header      *multipart.FileHeader
	UploadedUrl string
}

package file

import (
	"DX/src/domain/entity/file"
)

type UseCase interface {
	UploadFile(string, string, *file.File) (string, error)
	UploadFiles(string, string, []*file.File) ([]string, error)
}

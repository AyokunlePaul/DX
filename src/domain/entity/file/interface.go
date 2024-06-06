package file

type writer interface {
	Create(string, string, *File) error
	CreateList(string, string, []*File) error
	UploadCategoryIcon(string, *File) error
}

type Repository interface {
	writer
}

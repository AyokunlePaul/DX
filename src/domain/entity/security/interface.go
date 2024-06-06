package security

type reader interface {
	Get(string) (*Security, error)
}

type writer interface {
	Create(*Security) error
}

type Repository interface {
	reader
	writer
}

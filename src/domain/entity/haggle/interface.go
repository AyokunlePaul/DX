package haggle

type reader interface {
	GetHagglesForBid(string)
}

type writer interface {
	CreateHaggle(*Haggle) error
}

type Repository interface {
	reader
	writer
}

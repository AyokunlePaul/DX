package errand

import "go.mongodb.org/mongo-driver/bson"

type NewWriter interface {
	Create(*Errand) error
	Update(bson.D) (Errand, error)
	Delete(string) error
}

type NewReader interface {
	Get(string) (*Errand, error)
	GetDraft(string) (*Errand, error)
	GetFor(string) ([]Errand, error)
	GetAll() ([]Errand, error)
	Search(string) ([]string, error)
}

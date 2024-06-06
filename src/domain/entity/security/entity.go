package security

import (
	"errors"
)

type Security struct {
	Question    string `json:"question,omitempty" bson:"question"`
	Answer      string `json:"-" bson:"answer"`
	UserId      string `json:"-" bson:"user_id"`
	PhoneNumber string `json:"-" bson:"phone_number"`
}

func Create(data map[string]interface{}) (*Security, error) {
	var question, answer string
	var ok bool

	if question, ok = data["question"].(string); !ok {
		return nil, errors.New("question is required")
	}
	if answer, ok = data["answer"].(string); !ok {
		return nil, errors.New("answer is required")
	}
	return &Security{
		Question: question,
		Answer:   answer,
	}, nil
}

func GetAnswer(data map[string]interface{}) (string, string, error) {
	var phone, answer string
	var ok bool
	if phone, ok = data["phone_number"].(string); !ok {
		return "", "", errors.New("phone number is required")
	}
	if answer, ok = data["answer"].(string); !ok {
		return "", "", errors.New("answer is required")
	}
	return phone, answer, nil
}

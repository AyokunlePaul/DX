package entity

import "errors"

type Source int

const (
	Sender Source = iota
	Runner
	Admin
)

func GetSource(value string) (Source, error) {
	if value == Sender.Id() {
		return Sender, nil
	}
	if value == Runner.Id() {
		return Runner, nil
	}
	if value == Admin.Id() {
		return Admin, nil
	}
	return -1, errors.New("invalid bid source")
}

func (s Source) String() string {
	if s == Sender {
		return "Sender"
	}
	if s == Runner {
		return "Runner"
	}
	if s == Admin {
		return "Admin"
	}
	return ""
}

func (s Source) Id() string {
	if s == Sender {
		return "sender"
	}
	if s == Runner {
		return "runner"
	}
	if s == Admin {
		return "admin"
	}
	return ""
}

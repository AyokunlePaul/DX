package bid

const (
	Open State = iota
	Accepted
	Rejected
)

type State int

func (s State) String() string {
	if s == Open {
		return "Open"
	}
	if s == Accepted {
		return "Accepted"
	}
	if s == Rejected {
		return "Rejected"
	}
	return ""
}

func (s State) Id() string {
	if s == Open {
		return "open"
	}
	if s == Accepted {
		return "accepted"
	}
	if s == Rejected {
		return "rejected"
	}
	return ""
}

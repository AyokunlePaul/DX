package errand

type State int

const (
	Draft State = iota
	Open
	Pending //Bid has been accepted (pending)
	Active  // Runner has started errand (active)
	Completed
	Review
	EditMode
	Cancelled
	RunnerCompleted
	Abandoned
)

func (s State) String() string {
	if s == Draft {
		return "Draft"
	}
	if s == Open {
		return "Open"
	}
	if s == Pending {
		return "Pending"
	}
	if s == Active {
		return "Active"
	}
	if s == Completed {
		return "Completed"
	}
	if s == RunnerCompleted {
		return "Runner Completed"
	}
	if s == Review {
		return "Review"
	}
	if s == EditMode {
		return "Edit-Mode"
	}
	if s == Cancelled {
		return "Cancelled"
	}
	if s == Abandoned {
		return "Abandoned"
	}
	return ""
}

func (s State) Id() string {
	if s == Draft {
		return "draft"
	}
	if s == Open {
		return "open"
	}
	if s == Pending {
		return "pending"
	}
	if s == Active {
		return "active"
	}
	if s == Completed {
		return "completed"
	}
	if s == RunnerCompleted {
		return "runner-completed"
	}
	if s == Review {
		return "review"
	}
	if s == EditMode {
		return "edit-mode"
	}
	if s == Cancelled {
		return "cancelled"
	}
	if s == Abandoned {
		return "abandoned"
	}
	return ""
}

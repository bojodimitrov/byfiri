package cli

// Command enum
type Command int

// Command enum values
const (
	Open Command = iota
	Make
	Delete
	Edit
	Rename
	List
	Tree
	Exit
)

func (s Command) String() string {
	return [...]string{"open", "make", "del", "edit", "ren", "list", "tree", "exit"}[s]
}

// MakeParameters enum
type MakeParameters int

// MakeParameters enum values
const (
	File MakeParameters = iota
	Directory
)

func (s MakeParameters) String() string {
	return [...]string{"f", "d"}[s]
}

package constants

type RoleType uint8

func (d RoleType) String() string {
	strings := [...]string{
		"member",
		"admin",
	}

	if d < Member || d > Admin {
		return "Unknown"
	}
	return strings[d-1]
}

func (d RoleType) EnumIndex() int {
	return int(d)
}

const (
	Member RoleType = iota + 1 // EnumIndex = 1
	Admin                      // EnumIndex = 2
)

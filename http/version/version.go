package version

import (
	"fmt"
)

type Version struct {
	Major, Minor int
}

func (v Version) String() string {
	if v.Major > 1 && v.Minor == 0 {
		return fmt.Sprintf("HTTP/%d", v.Major)
	}
	return fmt.Sprintf("HTTP/%d.%d", v.Major, v.Minor)
}

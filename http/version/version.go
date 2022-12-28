package version

import (
	"fmt"

	"go.x2ox.com/sorbifolia/pyrokinesis"
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

var versions = [4][10][]byte{
	{
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
		[]byte(Version{0, 9}.String()),
	},
	{
		[]byte(Version{1, 0}.String()),
		[]byte(Version{1, 1}.String()),
	},
	{
		[]byte(Version{2, 0}.String()),
	},
	{
		[]byte(Version{3, 0}.String()),
	},
}

func (v Version) Bytes() []byte {
	if b := versions[v.Major][v.Minor]; len(b) != 0 {
		return b
	}
	return pyrokinesis.String.ToBytes(v.String())
}

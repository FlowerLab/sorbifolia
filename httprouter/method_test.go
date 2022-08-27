package httprouter

import (
	"testing"
)

func TestNewMethod(t *testing.T) {
	arr := []struct {
		m Method
		v string
	}{
		{GET, "GET"}, {GET, "GeT"},
		{HEAD, "HEAD"}, {HEAD, "HEaD"},
		{POST, "POST"}, {POST, "POsT"},
		{PUT, "PUT"}, {PUT, "PUt"},
		{PATCH, "PATCH"}, {PATCH, "PaTCH"},
		{DELETE, "DELETE"}, {DELETE, "DElETE"},
		{CONNECT, "CONNECT"}, {CONNECT, "CoNNECT"},
		{OPTIONS, "OPTIONS"}, {OPTIONS, "OPTIonS"},
		{TRACE, "TRACE"}, {TRACE, "TRacE"},
		{255, ""}, {255, "psot"},
	}

	for _, v := range arr {
		if NewMethod(v.v) != v.m {
			t.Error("fail", v)
		}
	}
}

package reflectype

import (
	"encoding/json"
	"reflect"
)

var JSONRawMessage = reflect.TypeOf((*json.RawMessage)(nil)).Elem()

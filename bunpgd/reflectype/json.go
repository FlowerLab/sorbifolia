package reflectype

import (
	"encoding/json"
	"reflect"
)

var JSONRawMessage = reflect.TypeFor[json.RawMessage]()

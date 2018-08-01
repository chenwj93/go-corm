package struct_utils

import (
	"reflect"
	"utils"
	"encoding/json"
)

func Asset(key *reflect.Value, value interface{}) {
	switch key.Type().Kind() {
	case reflect.String:
		key.SetString(utils.ParseString(value))
	default:
		valJ, _ := json.Marshal(value)
		json.Unmarshal(valJ, key.Addr().Interface())
	}
}
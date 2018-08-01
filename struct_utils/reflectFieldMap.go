package struct_utils

import (
	"reflect"
	"go-corm/logs"
	"fmt"
)

type ReflectFieldMap map[string]reflect.Value

func NewReflectFieldMap() ReflectFieldMap{
	return make(map[string]reflect.Value)
}

func (m1 ReflectFieldMap) mergeMap(m2 ReflectFieldMap) {
	for k, v := range m2 {
		if _, ok := m1[k]; ok {
			logs.Warn(fmt.Sprintf("duplicate key [%s]", k))
		} else {
			m1[k] = v
		}
	}
	return
}

// merge m2 to m1
func MergeStringMap(m1, m2 map[string]string) {
	for k, v := range m2 {
		if _, ok := m1[k]; ok {
			logs.Warn(fmt.Sprintf("duplicate key [%s]", k))
		} else {
			m1[k] = v
		}
	}
	return
}

package corm

import (
	"fmt"
	"go-corm/errorHandle"
	"go-corm/logs"
	"reflect"
	"sync"
	"go-corm/struct_utils"
)

// 每个运行单元可能包含多个struct
// 所以一个运行单元对应的field map 是一个双层的map，即每个struct对应一个map
// 每个运行单元最多只对应一个表
type ExecuteUnit struct {
	Structs    []string
	Table      string
	PK         []string
	StructPtr  map[string]interface{}
	FieldMap   map[string]struct_utils.ReflectFieldMap
	TagMap     map[string]struct_utils.ReflectFieldMap
	FieldToTag map[string]map[string]string
	TagToField map[string]map[string]string
	IfCatch    bool
	CatchMap   Cache
	mutex      sync.RWMutex
}

func RegisterStructs(ifCatch bool, table string, structs ...interface{}) {
	defer errorHandle.CatchError()
	exut := &ExecuteUnit{IfCatch:ifCatch,
						Table:table,
						FieldMap:struct_utils.NewMapRef(),
						TagMap:struct_utils.NewMapRef(),
						FieldToTag:struct_utils.NewMapMapString(),
						TagToField:struct_utils.NewMapMapString(),
						StructPtr:make(map[string]interface{}),
						CatchMap:NewCache()}
	for _, stru := range structs {
		structName := reflect.ValueOf(stru).Elem().Type().Name()
		f, t, ft, pk, err := struct_utils.Analysis(stru)
		if err != nil {
			logs.Fatal(fmt.Sprintf("struct analysis error [%s]", err.Error()))
			return
		}
		tf := make(map[string]string)
		for f, t := range ft {
			if _, ok := tf[t]; ok {
				logs.Fatal(fmt.Sprintf("a column not allow reflect to multiple field [%s]", t))
			}
			tf[t] = f
		}
		exut.StructPtr[structName] = stru
		exut.FieldMap[structName] = f
		exut.TagMap[structName] = t
		exut.FieldToTag[structName] = ft
		exut.TagToField[structName] = tf
		exut.Structs = append(exut.Structs, structName)
		exut.PK = append(exut.PK, pk...)
	}
	for _, structName := range exut.Structs {
		Orm.ExecuteUnitMap[structName] = exut
	}
	Orm.ExecuteUnitMapByTable[table] = exut
}

func (c *ExecuteUnit) IsThisUnit(structName string) bool {
	for _, s := range c.Structs {
		if s == structName {
			return true
		}
	}
	return false
}

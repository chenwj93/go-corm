package struct_utils

import (
	"reflect"

	"go-corm/errorHandle"
	"strings"
	"time"
)

func Analysis(o interface{}) (mapField, mapTag ReflectFieldMap, mapFieldToTag map[string]string, pk []string, err error) {
	defer errorHandle.CatchLoadDataError(&err)

	mapField, mapTag, mapFieldToTag = NewReflectFieldMap(), NewReflectFieldMap(), make(map[string]string)
	refE := reflect.ValueOf(o).Elem()
	for i := 0; i < refE.NumField(); i++ {
		if refE.Field(i).Kind() == reflect.Struct && !refE.Field(i).Type().AssignableTo(reflect.TypeOf(time.Time{})) {
			mField, mTag, mFtoT, pkRet, err := Analysis(refE.Field(i).Addr().Interface())
			if err != nil {
				return nil, nil, nil, nil, err
			}
			mapField.mergeMap(mField)
			mapTag.mergeMap(mTag)
			MergeStringMap(mapFieldToTag, mFtoT)
			pk = append(pk, pkRet...)
		} else if refE.Field(i).CanSet() {
			mapField[refE.Type().Field(i).Name] = refE.Field(i)
			if fieldTag := refE.Type().Field(i).Tag.Get("corm"); fieldTag != "" { // 有corm 标签
				// fieldTag: user_id[,pk]
				tagArr := strings.Split(fieldTag, ",")
				if len(tagArr) > 1 && (strings.ToUpper(strings.TrimSpace(tagArr[1])) == "PK") { // 判断是否主键标识
					pk = append(pk, strings.TrimSpace(tagArr[0]))
				}
				mapTag[strings.TrimSpace(tagArr[0])] = refE.Field(i)
				mapFieldToTag[refE.Type().Field(i).Name] = strings.TrimSpace(tagArr[0])
			}
		}
	}
	return
}

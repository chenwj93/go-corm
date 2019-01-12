package corm

import (
	"database/sql"
	"reflect"
	"github.com/chenwj93/utils"

	"fmt"
	"go-corm/errorHandle"
	"go-corm/logs"
	"go-corm/struct_utils"
	"encoding/json"
	"errors"
)

type Row struct {
	slt 	string
	args	[]interface{}
	Rows 	*sql.Rows
}

func NewRow(slt string, args ...interface{}) *Row{
	return &Row{slt:slt, args:args}
}

func (r *Row) ToMap(result *[]map[string]interface{}) error {
	_, _, err := r.query(utils.EMPTY_STRING, nil)
	if err != nil{
		return err
	}

	//返回所有列
	cols, err := r.Rows.Columns()
	if err != nil {
		return err
	}
	//这里表示一行所有列的值，用[]byte表示
	vals := make([]interface{}, len(cols))
	//这里表示一行填充数据
	scans := make([]interface{}, len(cols))
	//这里scans引用vals，把数据填充到[]byte里
	for k := range vals {
		scans[k] = &vals[k]
	}
	for r.Rows.Next() {
		//填充数据
		r.Rows.Scan(scans...) //因为Scan的参数是地址数组，所以才需要scans作为中间变量
		//每行数据
		row := make(map[string]interface{})
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := utils.ParseCamelCase(cols[k])
			//这里把[]byte数据转成string
			row[key] = utils.ParseString(v)
		}
		//放入结果集
		*result = append(*result, row)
	}

	return nil
}

func (r *Row) ToRow(scanner interface{}) (err error) {
	refScan := reflect.ValueOf(scanner)
	if refScan.Type().Kind() != reflect.Ptr{
		err = errors.New("receiver is not a ptr")
		logs.Error(err.Error())
		return
	}
	ok, cache, err := r.query(refScan.Elem().Type().Name(), scanner)
	if err != nil{
		return err
	}
	if ok{
		return
	}

	var cols []string
	cols, err = r.Rows.Columns()
	if err == nil{
		if r.Rows.Next() {
			typ := refScan.Elem().Type()
			object := reflect.New(typ)
			var result interface{}
			result, err = resultToStructs(cols, r.Rows, object.Interface())
			resJ, _ := json.Marshal(result)
			err = json.Unmarshal(resJ, scanner)
			if cache != nil{
				go cache.SetCache(GenKey(r.slt, r.args), scanner)
			}
		}
	} else {
		logs.Error(err)
	}
	return
}

func (r *Row) ToRows(scanner interface{}) (err error) {
	refScan := reflect.ValueOf(scanner)
	ind := reflect.Indirect(refScan)
	if ind.Type().Kind() != reflect.Slice {
		logs.Error("receiver is not a slice address")
		return
	}
	typ := ind.Type().Elem()
	ok, cache, err := r.query(typ.Name(), scanner)
	if err != nil{
		return err
	}
	if ok{
		return
	}
	//返回所有列
	cols, err := r.Rows.Columns()
	if err != nil {
		logs.Error(err)
	}
	var result []interface{}
	for r.Rows.Next() {
		object := reflect.New(typ)
		res, err := resultToStructs(cols, r.Rows, object.Interface())
		if err != nil {
			return err
		}
		result = append(result, res)
	}
	resJ, _ := json.Marshal(result)
	json.Unmarshal(resJ, scanner)
	if cache != nil{
		go cache.SetCache(GenKey(r.slt, r.args), scanner)
	}
	return
}

func resultToStructs(cols []string, row *sql.Rows, object interface{}) (res interface{}, err error) {
	defer errorHandle.CatchLoadDataError(&err)

	scans := make([]interface{}, len(cols))
	vals := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i ++{
		scans[i] = &vals[i]
	}

	if warn := row.Scan(scans...); warn != nil{
		logs.Warn(warn)
	}

	fieldMap, tagMap, _, _, err := struct_utils.Analysis(object)
	if err != nil {
		return nil, err
	}

	for k, v := range cols {
		var field reflect.Value
		var ok bool
		if field, ok = tagMap[v]; !ok {
			if field, ok = fieldMap[utils.FirstCapital(utils.ParseCamelCase(v))]; !ok {
				logs.Warn(fmt.Sprintf("column find variable failed [%s]", v))
				scans[k] = new(interface{})
			}
		}
		if ok && field.CanAddr() && field.Addr().CanInterface() {
			struct_utils.Asset(&field, vals[k])
		}else {
			logs.Warn(fmt.Sprintf("field[%s] can not set value", v))
		}
	}

	return object, err
}


//
func (r *Row)query(structName string, scanner interface{}) (ok bool, cache *Cache, err error){
	exut, ok := Orm.ExecuteUnitMap[structName]
	if ok && exut.IfCatch{
		if val, ok := exut.CatchMap.GetCache(GenKey(r.slt, r.args)); ok {
			err = json.Unmarshal(val, scanner)
			if err == nil{
				return ok, nil, nil
			}
			logs.Info(err)
		}
	}
	logs.Info(r.slt, r.args)
	r.Rows, err = Orm.DB.Query(r.slt, r.args...)
	if nil != err {
		logs.Error(err)
		return false, nil, err
	}
	if ok {
		return false, &exut.CatchMap, nil
	}
	return false, nil, nil
}

func marshalVal(src interface{}, desc interface{}) (err error){
	srcJ, _ := json.Marshal(src)
	err = json.Unmarshal(srcJ, &desc)
	if err != nil{
		logs.Error(err)
	}
	return
}
package exec

import (
	"database/sql"
	"fmt"
	"go-corm/builder"
	"go-corm/corm"
	"go-corm/errorHandle"
	"go-corm/logs"
	"log"
	"reflect"
	"github.com/chenwj93/utils"
	"go-corm/struct_utils"
)

type Ormer struct {
	tx      *sql.Tx
	f 		func(string, ...interface{}) (sql.Result, error)
	ifTrans bool
	exut 	*corm.ExecuteUnit
}

func NewOrm() *Ormer {
	o := Ormer{ifTrans:false, f:corm.Orm.DB.Exec}
	return &o
}

func (o *Ormer) Begin() {
	var e1 error
	o.ifTrans = true
	o.tx, e1 = corm.Orm.DB.Begin()
	o.f = o.tx.Exec
	if e1 != nil {
		log.Println(e1)
	}
}

func (o *Ormer) Commit() error {
	o.ifTrans = false
	o.f = corm.Orm.DB.Exec
	return o.tx.Commit()
}

func (o *Ormer) Rollback() error {
	o.ifTrans = false
	o.f = corm.Orm.DB.Exec
	return o.tx.Rollback()
}

func (o *Ormer) Query(sqlCon string, args ...interface{}) *corm.Row {

	return corm.NewRow(sqlCon, args...)
}

func (o *Ormer) Exec(sqlCon string, args ...interface{}) (sql.Result, error) {
	exut := o.getExut(sqlCon)
	if exut != nil{
		exut.CatchMap.EmptyCache()
	}
	logs.Info(sqlCon, args)
	res, err := o.f(sqlCon, args...)
	if err != nil {
		fmt.Println(err)
	}
	return res, err
}

func (o *Ormer) Insert(objectAddr interface{}) (sql.Result, error) {
	defer errorHandle.CatchError()
	refObj, beanName, exut, err := GetExut(objectAddr)
	if err != nil {
		return nil, err
	}
	ftoT, cols, args := exut.FieldToTag[beanName], []string{}, []interface{}{}

	for f, t := range ftoT {
		if field := refObj.FieldByName(f); field.CanSet() {
			if utils.IsContain(exut.PK, t) && utils.IsEmpty(field.Interface()) { // 主键为空
				continue
			}
			cols = append(cols, t)
			args = append(args, field.Interface())
		} else {
			logs.Error(fmt.Sprintf("field [%s] can't access", f))
		}
	}
	var insertStat builder.Insert
	insertStat.Tb(exut.Table).Cols(cols...).Args(args...)
	o.exut = exut
	return o.Exec(insertStat.GenStat(), insertStat.GenArgs()...)
}

func (o *Ormer) Update(objectAddr interface{}, columns ...string) (sql.Result, error) {
	defer errorHandle.CatchError()
	refObj, beanName, exut, err := GetExut(objectAddr)
	if err != nil {
		return nil, err
	}
	tToF, cols, args := exut.TagToField[beanName], []string{}, []interface{}{}
	updateStat := builder.NewUpdate()

	for _, pk := range exut.PK {
		if field := refObj.FieldByName(tToF[pk]); field.CanSet() {
			updateStat.Eq(pk, field.Interface())
		} else {
			logs.Error(fmt.Sprintf("primary key field [%s] value can not access", pk))
		}
	}
	for _, column := range columns {
		if utils.IsContain(exut.PK, column) { // is 主键
			logs.Warn(fmt.Sprintf("primary key field [%s] can't modified at auto model", column))
			continue
		}
		if field := refObj.FieldByName(tToF[column]); field.CanSet() {
			cols = append(cols, column)
			args = append(args, field.Interface())
		} else {
			logs.Error(fmt.Sprintf("field [%s/%s] can't access", column, tToF[column]))
		}
	}

	updateStat.Tb(exut.Table).SetCols(cols...).SetArgs(args...)
	o.exut = exut
	return o.Exec(updateStat.GenStat(), updateStat.GenArgs()...)
}

func (o *Ormer) Read(objectAddr interface{}) (err error) {
	defer errorHandle.CatchError()
	refObj, beanName, exut, err := GetExut(objectAddr)
	if err != nil {
		return err
	}
	tToF := exut.TagToField[beanName]
	selectStat := builder.NewSelect()
	for k := range tToF {
		selectStat.Slt(k)
	}
	selectStat.Tb(exut.Table)
	for _, pk := range exut.PK {
		if field := refObj.FieldByName(tToF[pk]); field.CanSet() {
			selectStat.Eq(pk, field.Interface())
		} else {
			logs.Error(fmt.Sprintf("primary key field [%s] value can not access", pk))
		}
	}
	err = o.Query(selectStat.GenCom().ToString(), selectStat.GetParamWhere()...).ToRow(objectAddr)

	return
}

func (o *Ormer) Delete(objectAddr interface{}) (sql.Result, error) {
	defer errorHandle.CatchError()
	refObj, beanName, exut, err := GetExut(objectAddr)
	if err != nil {
		return nil, err
	}
	tToF := exut.TagToField[beanName]
	deleteStat := builder.NewDelete()

	for _, pk := range exut.PK {
		if field := refObj.FieldByName(tToF[pk]); field.CanSet() {
			deleteStat.Eq(pk, field.Interface())
		} else {
			logs.Error(fmt.Sprintf("primary key field [%s] value can not access", pk))
		}
	}

	deleteStat.Tb(exut.Table)
	o.exut = exut
	return o.Exec(deleteStat.GenStat(), deleteStat.GenArgs()...)
}

func (o *Ormer)getExut(sql string) (e *corm.ExecuteUnit){
	if o.exut == nil{
		table := struct_utils.SelectTable(sql)
		if table != utils.EMPTY_STRING {
			o.exut = corm.Orm.ExecuteUnitMapByTable[table]
		}
	}
	if o.exut != nil && o.exut.IfCatch{
		e = o.exut
	}
	return
}

func GetExut(objectAddr interface{}) (refObj reflect.Value, beanName string, exut *corm.ExecuteUnit, err error) {
	refObj = reflect.ValueOf(objectAddr).Elem()
	beanName = refObj.Type().Name()
	var ok bool
	exut, ok = corm.Orm.ExecuteUnitMap[beanName]
	if !ok || exut.Table == utils.EMPTY_STRING {
		return refObj, utils.EMPTY_STRING, nil, fmt.Errorf("unregister table")
	}
	return
}
package corm

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"time"

	"go-corm/logs"
)

var Orm Corm

// 每个struct都对应一个执行单元，而一个执行单元中可能包含了若干个struct的结构信息
type Corm struct {
	ExecuteUnitMap map[string]*ExecuteUnit
	ExecuteUnitMapByTable map[string]*ExecuteUnit
	DB             *sql.DB
}

func Register(dataSourceName string, maxOpenConns, maxIdleConns int, connMaxLifeTime int64) (err error) {
	Orm.DB, err = sql.Open("mysql", dataSourceName)
	if nil != err {
		log.Printf("[corm]open database error: %s", err.Error())
	}
	Orm.ExecuteUnitMap, Orm.ExecuteUnitMapByTable = make(map[string]*ExecuteUnit), make(map[string]*ExecuteUnit)

	err = Orm.DB.Ping()
	if nil != err {
		log.Println("[corm]connect database failed")
	}
	if maxIdleConns > 0 {
		Orm.DB.SetMaxIdleConns(maxIdleConns)
	}
	if maxOpenConns > 0 {
		Orm.DB.SetMaxOpenConns(maxOpenConns)
	}
	if connMaxLifeTime > 0 {
		Orm.DB.SetConnMaxLifetime(time.Duration(connMaxLifeTime))
	}

	logs.DebugLevel = logs.ERROR
	return
}

func Debug(level logs.DEBUGLEVEL) {
	logs.DebugLevel = level
}

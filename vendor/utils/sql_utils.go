package utils

import (
	"strconv"
	"strings"
)

type SqlSelect struct {
	table      string
	slt        string
	where      string
	lmt        string
	orderBy    string
	groupBy    string
	sqlRet     string
	paramWhere []interface{}
	paramIn		[]interface{}
}

func (t *SqlSelect) GetWhere() *SqlSelect {
	if strings.TrimSpace(t.where) != EMPTY_STRING {
		t.sqlRet += " where " + t.where[4:]
	}
	return t
}

func (t *SqlSelect) GetLmt() *SqlSelect {
	t.sqlRet += t.lmt
	return t
}

func (t *SqlSelect) GetOrderBy() *SqlSelect {
	t.sqlRet += t.orderBy
	return t
}

func (t *SqlSelect) GetGroupBy() *SqlSelect {
	t.sqlRet += t.groupBy
	return t
}

func (t *SqlSelect) GetParamWhere() []interface{} {
	return t.paramWhere
}

func (t *SqlSelect) ToString() string {
	var s = t.sqlRet
	t.sqlRet = EMPTY_STRING
	return s
}

func (t *SqlSelect) New(table string) *SqlSelect {
	t.table = table
	return t
}

func (t *SqlSelect) Slt(col ...string) *SqlSelect{
	t.sqlRet = "select "
	if len(col) == 0{
		t.sqlRet += "*"
	}else {
		t.sqlRet += ParseStringFromArray(col, COMMA, "''")
	}
	return t
}

func (t *SqlSelect) From(table string) *SqlSelect{
	t.sqlRet += " from" + table
	return t
}

func (t *SqlSelect) Equal(key string, v interface{}) *SqlSelect {
	t.where += " and " + key + " = ?"
	t.paramWhere = append(t.paramWhere, v)
	return t
}

func (t *SqlSelect) Like(key string, v interface{}) *SqlSelect {
	t.where += " and " + key + " like ?"
	t.paramWhere = append(t.paramWhere, Wrap(v, "%"))
	return t
}

func (t *SqlSelect) Greater(key string, v interface{}) *SqlSelect {
	t.where += " and " + key + " > ?"
	t.paramWhere = append(t.paramWhere, v)
	return t
}

func (t *SqlSelect) GreaterAndEqual(key string, v interface{}) *SqlSelect {
	t.where += " and " + key + " >= ?"
	t.paramWhere = append(t.paramWhere, v)
	return t
}

func (t *SqlSelect) Less(key string, v interface{}) *SqlSelect {
	t.where += " and " + key + " < ?"
	t.paramWhere = append(t.paramWhere, v)
	return t
}

func (t *SqlSelect) LessAndEqual(key string, v interface{}) *SqlSelect {
	t.where += " and " + key + " <= ?"
	t.paramWhere = append(t.paramWhere, v)
	return t
}

func (t *SqlSelect) Limit(Page interface{}, Rows interface{}) *SqlSelect {
	page := ParseInt(Page)
	rows := ParseInt(Rows)
	if (page == 0 || rows == 0) && page != -1 {
		page = 1
		rows = 10
	}
	if page != -1 {
		t.lmt = " limit " + strconv.Itoa(ParseInt((page-1)*rows)) + ", " + strconv.Itoa(ParseInt(rows))
	}
	return t
}

func (t *SqlSelect) OrderBy(orderBy ...string) *SqlSelect {
	s := ParseStringFromArray(orderBy, ", ", EMPTY_STRING)
	if s != EMPTY_STRING {
		t.orderBy = " order by " + s
	}
	return t
}

func (t *SqlSelect) GroupBy(groupBy ...string) *SqlSelect {
	s := ParseStringFromArray(groupBy, ", ", EMPTY_STRING)
	if s != EMPTY_STRING {
		t.groupBy = " group by " + s
	}
	return t
}

// and col in (?, ?, ? ...), args...
func (t *SqlSelect) In(col string, args ...interface{}) *SqlSelect {
	if len(args) <= 0 {
		t.where += " and 1 = 0 "
	} else {
		t.where += " and " + col + " in ("
		for i := 0; i < len(args); i++ {
			t.where += "?, "
		}
		t.where = t.where[:len(t.where)-2] + ")"
		t.paramWhere = append(t.paramWhere, args...)
	}
	return t
}

func (t *SqlSelect) Custom(s string, v ...interface{}) *SqlSelect {
	t.where += " " + s
	t.paramWhere = append(t.paramWhere, v...)
	return t
}


func FormatIn(col string, length int) (ret string) {
	if length <= 0 {
		ret = " 1 = 0 "
	} else {
		ret = col + " in ("
		for i := 0; i < length; i++ {
			ret += "?, "
		}
		ret = ret[:len(ret)-2] + ")"
	}
	return
}

package builder

import (
	"strconv"
	"utils"
	"go-corm/logs"
)

type SqlSelect struct {
	table      string
	slt        string
	lmt        string
	orderBy    string
	groupBy    string
	sqlRet     string
	Where
}

func NewSelect() *SqlSelect{
	return &SqlSelect{}
}

func (t *SqlSelect) GenCom() *SqlSelect {
	if t.table == utils.EMPTY_STRING {
		logs.Error("undefined table")
	}
	if t.slt == utils.EMPTY_STRING {
		t.slt = "select *"
	}
	t.sqlRet += t.slt + t.table + t.GetWhere().ToString()
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

func (t *SqlSelect) ToString() string {
	var s = t.sqlRet
	t.sqlRet = utils.EMPTY_STRING
	return s
}


func (t *SqlSelect) Slt(col ...string) *SqlSelect{
	if t.slt == utils.EMPTY_STRING{
		t.slt = "select " + utils.ParseStringFromArray(col, utils.COMMA, "''")
	} else {
		t.slt += "," + utils.ParseStringFromArray(col, utils.COMMA, "''")
	}
	return t
}

func (t *SqlSelect) Tb(table string) *SqlSelect{
	t.table = " from " + table
	return t
}

func (t *SqlSelect) Limit(Page interface{}, Rows interface{}) *SqlSelect {
	page := utils.ParseInt(Page)
	rows := utils.ParseInt(Rows)
	if (page == 0 || rows == 0) && page != -1 {
		page = 1
		rows = 10
	}
	if page != -1 {
		t.lmt = " limit " + strconv.Itoa((page-1)*rows) + ", " + strconv.Itoa(rows)
	}
	return t
}

func (t *SqlSelect) OrderBy(orderBy ...string) *SqlSelect {
	s := utils.ParseStringFromArray(orderBy, ", ", utils.EMPTY_STRING)
	if s != utils.EMPTY_STRING {
		t.orderBy = " order by " + s
	}
	return t
}

func (t *SqlSelect) GroupBy(groupBy ...string) *SqlSelect {
	s := utils.ParseStringFromArray(groupBy, ", ", utils.EMPTY_STRING)
	if s != utils.EMPTY_STRING {
		t.groupBy = " group by " + s
	}
	return t
}

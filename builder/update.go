package builder

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strings"
	"github.com/chenwj93/utils"
)

type Update struct {
	table string
	cols  []string
	args  []interface{}
	Where
}

func NewUpdate() *Update {
	return &Update{}
}

func (o *Update) Tb(tb string) *Update {
	if o.table == utils.EMPTY_STRING {
		o.table = tb
	}
	return o
}

func (o *Update) SetCols(cols ...string) *Update {
	if len(o.cols) == 0 {
		o.cols = cols
	}
	return o
}

func (o *Update) SetArgs(args ...interface{}) *Update {
	if len(o.cols) != 0 && len(args) == len(o.cols) {
		o.args = append(o.args, args...)
	}
	return o
}

func (o *Update) GenStat() string {
	if len(o.cols) == 0 || len(o.args) == 0 || o.table == utils.EMPTY_STRING {
		return utils.EMPTY_STRING
	}
	var stat strings.Builder

	stat.WriteString(fmt.Sprintf("update %s set ", o.table))
	cols := strings.Join(o.cols, " = ?, ")
	stat.WriteString(cols)
	stat.WriteString(" = ? ")
	stat.WriteString(o.GetWhere().ToString())
	return stat.String()
}

func (o *Update) GenArgs() (param []interface{}) {
	param = append(param, o.args...)
	param = append(param, o.GetParamWhere()...)
	return
}

func (o *Update) Exec(or orm.Ormer) (sql.Result, error) {
	return or.Raw(o.GenStat(), o.GenArgs()...).Exec()
}

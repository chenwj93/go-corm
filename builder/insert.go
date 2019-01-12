package builder

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strings"
	"github.com/chenwj93/utils"
)

type Insert struct {
	table string
	cols  []string
	args  [][]interface{}
}

func NewInsert() *Insert {
	return &Insert{}
}

func (o *Insert) Tb(tb string) *Insert {
	if o.table == utils.EMPTY_STRING {
		o.table = tb
	}
	return o
}

func (o *Insert) Cols(cols ...string) *Insert {
	if len(o.cols) == 0 {
		o.cols = cols
	}
	return o
}

func (o *Insert) Args(args ...interface{}) *Insert {
	if len(o.cols) != 0 && len(args) == len(o.cols) {
		o.args = append(o.args, args)
	}
	return o
}

func (o *Insert) GenStat() string {
	if len(o.cols) == 0 || len(o.args) == 0 || o.table == utils.EMPTY_STRING {
		return utils.EMPTY_STRING
	}
	var stat strings.Builder
	var wildcards strings.Builder
	for i := 0; i < len(o.cols)-1; i++ {
		wildcards.WriteString("?, ")
	}
	wildcards.WriteString("?")
	cols := strings.Join(o.cols, ", ")
	stat.WriteString(fmt.Sprintf("insert into %s (%s) values ", o.table, cols))
	for i := 0; i < len(o.args)-1; i++ {
		stat.WriteString("(" + wildcards.String() + "), ")
	}
	stat.WriteString("(" + wildcards.String() + ")")
	return stat.String()
}

func (o *Insert) GenArgs() (param []interface{}) {
	for _, arg := range o.args {
		param = append(param, arg...)
	}
	return
}

func (o *Insert) Exec(or orm.Ormer) (sql.Result, error) {
	return or.Raw(o.GenStat(), o.GenArgs()...).Exec()
}

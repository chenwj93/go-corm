package builder

import (
	"fmt"
	"strings"
	"utils"
)

type Delete struct {
	table string
	Where
}

func NewDelete() *Delete {
	return &Delete{}
}

func (o *Delete) Tb(tb string) *Delete {
	if o.table == utils.EMPTY_STRING {
		o.table = tb
	}
	return o
}

func (o *Delete) GenStat() string {
	if o.table == utils.EMPTY_STRING {
		return utils.EMPTY_STRING
	}
	var stat strings.Builder

	stat.WriteString(fmt.Sprintf("delete from %s", o.table))
	stat.WriteString(o.GetWhere().ToString())
	return stat.String()
}

func (o *Delete) GenArgs() (param []interface{}) {
	return o.GetParamWhere()
}

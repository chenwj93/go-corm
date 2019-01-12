package builder

import (
	"strings"

	"github.com/chenwj93/utils"
)

type Where struct {
	where      string
	sqlRet     string
	paramWhere []interface{}
	paramIn    []interface{}
}

func (t *Where) GetWhere() *Where {
	if strings.TrimSpace(t.where) != utils.EMPTY_STRING {
		t.sqlRet += " where " + t.where[4:]
	}
	return t
}

func (t *Where) GetParamWhere() []interface{} {
	return t.paramWhere
}

func (t *Where) ToString() string {
	var s = t.sqlRet
	t.sqlRet = utils.EMPTY_STRING
	return s
}

// @param ifCheckNil 是否对v值判空， 无输入=false
func (t *Where) Eq(key string, v interface{}, ifCheckNil ...bool) *Where {
	if privateCheckParam(key, v, ifCheckNil) {
		t.where += " and " + key + " = ?"
		t.paramWhere = append(t.paramWhere, v)
	}
	return t
}

func (t *Where) Like(key string, v interface{}, ifCheckNil ...bool) *Where {
	if privateCheckParam(key, v, ifCheckNil) {
		t.where += " and " + key + " like ?"
		t.paramWhere = append(t.paramWhere, utils.Wrap(v, "%"))
	}
	return t
}

func (t *Where) LL(key string, v interface{}, ifCheckNil ...bool) *Where {
	if privateCheckParam(key, v, ifCheckNil) {
		t.where += " and " + key + " like ?"
		t.paramWhere = append(t.paramWhere, "%" + utils.ParseString(v))
	}
	return t
}

func (t *Where) RL(key string, v interface{}, ifCheckNil ...bool) *Where {
	if privateCheckParam(key, v, ifCheckNil) {
		t.where += " and " + key + " like ?"
		t.paramWhere = append(t.paramWhere, utils.ParseString(v) + "%")
	}
	return t
}

func (t *Where) Gt(key string, v interface{}, ifCheckNil ...bool) *Where {
	if privateCheckParam(key, v, ifCheckNil) {
		t.where += " and " + key + " > ?"
		t.paramWhere = append(t.paramWhere, v)
	}
	return t
}

func (t *Where) GtAndEq(key string, v interface{}, ifCheckNil ...bool) *Where {
	if privateCheckParam(key, v, ifCheckNil) {
		t.where += " and " + key + " >= ?"
		t.paramWhere = append(t.paramWhere, v)
	}
	return t
}

func (t *Where) Lt(key string, v interface{}, ifCheckNil ...bool) *Where {
	if privateCheckParam(key, v, ifCheckNil) {
		t.where += " and " + key + " < ?"
		t.paramWhere = append(t.paramWhere, v)
	}
	return t
}

func (t *Where) LtAndEq(key string, v interface{}, ifCheckNil ...bool) *Where {
	if privateCheckParam(key, v, ifCheckNil) {
		t.where += " and " + key + " <= ?"
		t.paramWhere = append(t.paramWhere, v)
	}
	return t
}

// and col in (?, ?, ? ...), args...
func (t *Where) In(col string, args ...interface{}) *Where {
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

func (t *Where) Custom(s string, v ...interface{}) *Where {
	if s != "" {
		t.where += " " + s
		t.paramWhere = append(t.paramWhere, v...)
	}
	return t
}

func privateCheckParam(key string, v interface{}, ifCheckNil []bool) bool {
	return key != "" && (len(ifCheckNil) == 0 ||!ifCheckNil[0] ||!utils.IsEmpty(v))
}

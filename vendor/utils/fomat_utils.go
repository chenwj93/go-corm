package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//create by cwj on 2017-09-16
func ParseString(val interface{}) string {
	if val != nil {
		switch v := val.(type) {
		case bool:
			if v {
				return "true"
			} else {
				return "false"
			}
		case string:
			return v
		case int64:
			return strconv.FormatInt(v, 10)
		case int8:
			return strconv.Itoa(int(v))
		case int32:
			return strconv.Itoa(int(v))
		case int:
			return strconv.Itoa(v)
		case float64:
			return strconv.FormatFloat(v, 'f', -1, 64)
		case float32:
			return strconv.FormatFloat(float64(v), 'f', -1, 32)
		case time.Time:
			return v.Format(TIME_FORMAT_1)
		case map[string]interface{}:
			s, _ := json.Marshal(val)
			return string(s)
		case []byte:
			return string(v)
		}
		return EMPTY_STRING
	}
	return EMPTY_STRING
}

//create by cwj on 2017-09-16
func ParseInt(val interface{}) int {
	if val != nil {
		switch v := val.(type) {
		case bool:
			if v {
				return 1
			} else {
				return 0
			}
		case string:
			i, _ := strconv.ParseInt(v, 10, 64)
			return int(i)
		case int64:
			return int(v)
		case int8:
			return int(v)
		case int32:
			return int(v)
		case int:
			return v
		case float64:
			return int(v)
		case float32:
			return int(v)
		case time.Time:
			return int(v.UnixNano())
		}
		return 0
	}
	return 0
}

//create by cwj on 2017-09-16
func ParseInt64(val interface{}) int64 {
	if val != nil {
		switch v := val.(type) {
		case bool:
			if v {
				return 1
			} else {
				return 0
			}
		case string:
			i, _ := strconv.ParseInt(v, 10, 64)
			return i
		case int64:
			return v
		case int8:
			return int64(v)
		case int32:
			return int64(v)
		case int:
			return int64(v)
		case float64:
			return int64(v)
		case float32:
			return int64(v)
		case time.Time:
			return v.UnixNano()
		}
		return 0
	}
	return 0
}

//create by cwj on 2017-11-25
func ParseFloat64(val interface{}) float64 {
	if val != nil {
		switch v := val.(type) {
		case bool:
			if v {
				return 1
			} else {
				return 0
			}
		case string:
			i, _ := strconv.ParseFloat(v, 64)
			return i
		case int64:
			return float64(v)
		case int8:
			return float64(v)
		case int32:
			return float64(v)
		case int:
			return float64(v)
		case float64:
			return v
		case float32:
			return float64(v)
		case time.Time:
			return float64(v.UnixNano())
		}
		return 0
	}
	return 0
}

//create by cwj on 2017-09-16
// format time from string in map
func ParseTimeForMap(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		ret, err := regexp.MatchString(TIME_REG_1, ParseString(v))
		if ret == true && err == nil {
			t, errParse := time.Parse(TIME_FORMAT_1, ParseString(v))
			if errParse == nil {
				m[k] = t
			}
			continue
		}
		ret, err = regexp.MatchString(TIME_REG_2, ParseString(v))
		if ret == true && err == nil {
			t, errParse := time.Parse(TIME_FORMAT_2, ParseString(v))
			if errParse == nil {
				m[k] = t
			}
			continue
		}
		ret, err = regexp.MatchString(TIME_REG_3, ParseString(v))
		if ret == true && err == nil {
			t, errParse := time.Parse(TIME_FORMAT_3, ParseString(v))
			if errParse == nil {
				m[k] = t
			}
			continue
		}
		ret, err = regexp.MatchString(TIME_REG_4, ParseString(v))
		if ret == true && err == nil {
			t, errParse := time.Parse(TIME_FORMAT_4, ParseString(v))
			if errParse == nil {
				m[k] = t
			}
		}
	}
	return m
}

//create by cwj on 2017-10-17
// parse string from array
func ParseStringFromArray(a []string, character string, empty string) string {
	var buff bytes.Buffer
	for i, e := range a {
		if e == EMPTY_STRING {
			e = empty
		}
		if i < len(a)-1 {
			buff.WriteString(ParseString(e) + character)
		} else {
			buff.WriteString(ParseString(e))
		}
	}
	return buff.String()
}

//create by cwj on 2017-10-17
// parse string from array
func ParseStringFromArrayWithQuote(a []string, character string) string {
	var buff bytes.Buffer
	for i, e := range a {
		if i < len(a)-1 {
			buff.WriteString("'" + ParseString(e) + "'" + character)
		} else {
			buff.WriteString("'" + ParseString(e) + "'")
		}
	}
	return buff.String()
}

func LinkString(character string, empty string, str ...string) string {
	return ParseStringFromArray(str, character, empty)
}

func ParseStringFromMap(m map[string]interface{}) string {
	s, _ := json.Marshal(m)
	return string(s)
}

func ParseMapFromString(s string) (map[string]interface{}, error) {
	var m = make(map[string]interface{})
	err := json.Unmarshal([]byte(s), &m)
	return m, err
}

func ParseMap(o interface{}) (map[string]interface{}, error) {
	var m = make(map[string]interface{})
	s, err := json.Marshal(o)
	if err != nil {
		return m, err
	}
	err = json.Unmarshal(s, &m)
	return m, err
}

func ParseStruct(val interface{}, object interface{}) (e error) {
	dataJson, _ := json.Marshal(val)
	e = json.Unmarshal(dataJson, object)
	return
}

func SubString(str string, head int, tail int, empty string) string {
	if head < 0 {
		head = 0
	}
	if tail < 0 {
		return str[head:]
	}
	if len(str) < head {
		return empty
	}
	if len(str) < tail {
		return str[head:]
	}
	if head > tail {
		return str[head:] + str[:tail]
	}
	return str[head:tail]
}

//create by cwj on 2017-10-17
//check errArray
//if there are error, return it
func ChangeFromCamelCase(s string) string {
	var snakeString strings.Builder
	for index, char := range s {
		if char >= 'A' && char <= 'Z' {
			if index > 0{
				snakeString.WriteByte('_')
			}
			snakeString.WriteByte(byte(char + 32))
		}else {
			snakeString.WriteByte(byte(char))
		}
	}
	return snakeString.String()
}

// create by cwj on 2018-03-02
// eg: order_user_log_  ->  orderUserLog
func ParseCamelCase(s string) string {
	ifToUp := false
	var camelString strings.Builder
	for _, char := range s {
		if char == '_' {
			ifToUp = true
		}else{
			if ifToUp && char <= 122 && char >= 97{
				camelString.WriteByte(byte(char - 32))
			}else {
				camelString.WriteByte(byte(char))
			}
			ifToUp = false
		}

	}
	return camelString.String()
}

//created by cwj 2018-03-16
// first character upper case
func FirstCapital(s string) string {
	if len(s) == 1 {
		s = strings.ToUpper(s[:1])
	}else if len(s) > 1{
		s = strings.ToUpper(s[:1]) + s[1:]
	}

	return s
}

func FirstLower(s string) string {
	if len(s) == 1 {
		s = strings.ToLower(s[:1])
	}else if len(s) > 1{
		s = strings.ToLower(s[:1]) + s[1:]
	}
	return strings.ToLower(s[:1] + s[1:])
}

//created by cwj 2018-03-16
// wrap string with character
func Wrap(str interface{}, character string) string {
	return character + ParseString(str) + character
}

//created by cwj 2018-03-16
// batch judge values if empty value
func IsEmptyBat(val ...interface{}) bool {
	for _, ele := range val {
		if IsEmpty(ele) {
			return true
		}
	}
	return false
}

//created by cwj 2018-03-16
// judge the value if empty value
func IsEmpty(val interface{}) bool {
	if val != nil {
		t := reflect.ValueOf(val)
		switch t.Kind() {
		case reflect.Bool:
			return false
		case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
			return t.Int() == 0
		case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
			return t.Uint() == 0
		case reflect.Float32, reflect.Float64:
			return t.Float() == 0
		case reflect.String:
			return t.String() == EMPTY_STRING
		case reflect.Slice:
			return t.Len() == 0
		case reflect.Map:
			return t.Len() == 0
		case reflect.Struct:
			if t.Type() == reflect.ValueOf(time.Now()).Type(){
				return t.Interface().(time.Time).IsZero()
			}
			return t.Interface() == nil
		}

	}
	return true
}

// created by cwj 2018-03-16
// judge whether there is an empty value in some fields of the object
func IsEmptyObject(object interface{}, key ...interface{}) bool {
	if len(key) == 0 {
		return false
	}
	reflectValue := reflect.ValueOf(object)

	for i := 0; i < reflectValue.NumField(); i++ {
		fieldName := reflectValue.Type().Field(i).Name
		fieldValue := reflectValue.Field(i).Interface()
		if IsContain(key, fieldName) && IsEmpty(fieldValue) {
			fmt.Println("检测到空========", key, fieldName, fieldValue)
			return true
		}
	}
	return false
}

// created by cwj 2018-03-16
// judge a array if contain the value
func IsContain(array interface{}, val interface{}) bool {
	if array == nil {
		return false
	}
	switch v := array.(type) {
	case map[string]interface{}:
		_, ok := v[ParseString(val)]
		return ok
	case []interface{}:
		for _, ele := range v {
			if ele == val {
				return true
			}
		}
	case []string:
		for _, ele := range v {
			if ele == val {
				return true
			}
		}
	// 有需要可以继续添加实现
	}
	return false
}

//created by cwj 2018-03-16
// O2 = O1
func Assignment(O1 interface{}, O2 interface{}) {
	dataJson, err := json.Marshal(O1)
	if err != nil {
		log.Println(err)
		return
	}
	json.Unmarshal(dataJson, O2)
}


func SimplifyFmtMap(m map[string]interface{}, l int){
	fmt.Print("map[")
	for key, value := range m{
		vb, _ := json.Marshal(value)
		vs := string(vb)
		if len(vs) > l{
			fmt.Printf("%s:%s +++ too long ", key, vs[:l])
		} else {
			fmt.Printf("%s:%s  ", key, vs)
		}
	}
	fmt.Println("]")
}


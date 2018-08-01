package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	//"utils_back"
)

/**
create by cwj on 2018-03-31
功能：1、自动扫描注册的controller文件，并自动注册router
	 2、自动urlmapping
*/

type function struct {
	requestFunc    string
	mainParamType  string
	channelDefault string
	channelType    string
}

type controller struct {
	object   interface{}
	conUrl   string
	fileName string
	funcMap  map[string]function
}

type Router struct {
	routeMap map[string]controller
	files    []interface{}
}

func (r *Router) GetUrl() (m map[string]controller){
	clone, _ := json.Marshal(r.routeMap)
	json.Unmarshal(clone, &m)
	return
}

func (r *Router) PrintURL() {
	for con, e := range r.routeMap {
		for funMap, fun := range e.funcMap {
			fmt.Println("-->", con+URL_SEPARATOR+funMap, ":\t", fun)
		}
	}
}

func (r *Router) URLMapping(operate string, paramMap map[string]interface{}) (res *Response, err error) {
	opList := strings.Split(operate, URL_SEPARATOR)
	var ok bool
	var con controller
	if con, ok = r.routeMap[opList[0]]; len(opList) != 2 || !ok {
		return ConcatNotFound()
	}

	reC := reflect.ValueOf(con.object)

	var fun function
	if fun, ok = con.funcMap[opList[1]]; !ok {
		return ConcatNotFound()
	}
	param := make([]reflect.Value, 0)
	if fun.mainParamType != EMPTY_STRING {
		switch fun.mainParamType {
		case "[]byte":
			paramJson, _ := json.Marshal(paramMap)
			param = append(param, reflect.ValueOf(paramJson))
		case "map[string]interface{}":
			param = append(param, reflect.ValueOf(paramMap))
		}
	}
	if fun.channelType != EMPTY_STRING {
		switch fun.channelType {
		case "string":
			param = append(param, reflect.ValueOf(fun.channelDefault))
		case "int":
			param = append(param, reflect.ValueOf(ParseInt(fun.channelDefault)))
		case "int8":
			param = append(param, reflect.ValueOf(int8(ParseInt(fun.channelDefault))))
		}
	}
	reFunc := reC.MethodByName(fun.requestFunc)
	if !reFunc.IsValid() {
		return ConcatNotFound()
	}
	ret := reFunc.Call(param)
	resI := ret[0].Interface()
	errI := ret[1].Interface()
	res, ok1 := resI.(*Response)
	err, _ = errI.(error)

	if !ok1 {
		res, err = ConcatNotFound()
	}

	fmt.Println(SubString(string(res.Json), 0, 500, ""))
	return res, err
}

func (r *Router) init(ControllerMap map[string]controller) error {
	filesCon, _ := ioutil.ReadDir("controllers")
	for _, file := range filesCon {
		conName, conUrl, funcMap, err := getFuncMap(file.Name())
		if err != nil {
			fmt.Println(err)
		}
		if len(funcMap) > 0 {
			if conUrl == EMPTY_STRING {
				conUrl = FirstLower(ParseCamelCase(SubString(conName, 0, strings.Index(conName, "Control"), EMPTY_STRING)))
			}
			ControllerMap[conName] = controller{nil, conUrl, file.Name(), funcMap}
		}
	}
	return nil
}

func (r *Router) AddRouters(cons ...interface{}) {
	var controllerMap = make(map[string]controller)
	r.init(controllerMap)
	if r.routeMap == nil {
		r.routeMap = make(map[string]controller)
	}
	for _, con := range cons {
		reC := reflect.TypeOf(con).Elem()
		if cont, ok := controllerMap[reC.Name()]; ok {
			if _, ok := r.routeMap[cont.conUrl]; ok {
				panic("controller注解重复")
			}
			cont.object = con
			r.routeMap[cont.conUrl] = cont
			r.files = append(r.files, cont.fileName)
		} else {
			panic(fmt.Errorf("未扫描到该controller：%s", reC.Name()))
		}

	}

}

func (r *Router) AddRouter(fileName string, con interface{}) error {
	fileName += ".go"
	if IsContain(r.files, fileName) {
		panic(fmt.Errorf("文件重复: %s", fileName))
	}
	r.files = append(r.files, fileName)

	conName, conUrl, funcMap, err := getFuncMap(fileName)
	if err != nil {
		return err
	}
	reC := reflect.TypeOf(con).Elem()
	if reC.Name() != conName {
		panic(fmt.Errorf("注册ctl[%s]与文件中ctl[%s]不一致！", reC.Name(), conName))
	}
	if conUrl == EMPTY_STRING {
		conUrl = FirstLower(ParseCamelCase(SubString(conName, 0, strings.Index(conName, "Control"), EMPTY_STRING)))
	}

	if r.routeMap == nil {
		r.routeMap = make(map[string]controller)
	}
	if _, ok := r.routeMap[conUrl]; ok {
		panic("controller重复注册:" + conName)
	}
	r.routeMap[conUrl] = controller{con, "", "", funcMap}

	return nil
}

func getFuncMap(fileName string) (controllerName string, conUrl string, funcMap map[string]function, err error) {

	fn := "controllers//" + fileName
	f, err := ioutil.ReadFile(fn)
	if err != nil {
		return EMPTY_STRING, EMPTY_STRING, nil, err
	}
	reader := bytes.NewReader(f)
	utilsReader := NewReader(bufio.NewReader(reader))

	funcMap = make(map[string]function)

	ifCon := true
	ifValidate := false
	line, lineStr, e := utilsReader.ReadLine()
	for line != -1 && e == nil {
		if ifCon && strings.Contains(lineStr, "@controller") {
			conUrl = lineStr[strings.Index(lineStr, "(")+1 : strings.Index(lineStr, ")")]
			ifCon = false
			line, lineStr, e = utilsReader.ReadLine()
			if line == -1 || !strings.Contains(lineStr, "type ") || !strings.Contains(lineStr, " struct") {
				return EMPTY_STRING, "", nil, errors.New("controller注解有误")
			}
			controllerName = strings.TrimSpace(lineStr[5:strings.Index(lineStr, " struct")])
			ifValidate = true
			continue
		}

		//说明没有@controller注解，采用默认的controller
		if ifCon && strings.Contains(lineStr, "type ") && strings.Contains(lineStr, " struct") {
			controllerName = strings.TrimSpace(lineStr[5:strings.Index(lineStr, " struct")])
			ifValidate = true
			ifCon = false
			line, lineStr, e = utilsReader.ReadLine()
			continue
		}

		if ifValidate && strings.Contains(lineStr, "@router") {
			err = getRouter(fn, line, lineStr, &utilsReader, funcMap, controllerName)
			if err != nil {
				return EMPTY_STRING, EMPTY_STRING, nil, err
			}
		}

		line, lineStr, e = utilsReader.ReadLine()
	}

	return
}

func getRouter(fileName string, line int, lineStr string, reader *Reader, funcMap map[string]function, con string) (err error) {
	var mc [][2]string
	for line != -1 && strings.Contains(lineStr, "@router") && strings.Contains(lineStr, "(") && strings.Contains(lineStr, ")") {
		m, c, e := handelString(fileName, line, lineStr[strings.Index(lineStr, "(")+1:strings.Index(lineStr, ")")], ";", "=", false)
		if e != nil {
			return e
		}
		if m == EMPTY_STRING {
			return fmt.Errorf("@router注解缺少method参数[%s：line %d]", fileName, line)
		}
		mc = append(mc, [2]string{m, c})
		line, lineStr, e = reader.ReadLine()
	}
	// lineStr = func (c *xxController) FuncName(param mainParamType |{channel channelType})(  ){}
	lineList := strings.Split(lineStr, ")")
	if len(lineList) != 4 || !strings.Contains(lineStr, "func") || !strings.Contains(lineStr, con) {
		return fmt.Errorf("@router注解无对应方法[%s：line %d]", fileName, line-1)
	}
	mainLine := lineList[1] //  FuncName(param mainParamType |{channel channelType}
	requestFunc := strings.TrimSpace(mainLine[:strings.Index(mainLine, "(")])
	var mainParamType, channelType string
	if !strings.HasSuffix(strings.TrimSpace(mainLine), "("){
		mainParamType, channelType, err = handelString(fileName, line, mainLine[strings.Index(mainLine, "(")+1:], ",", " ", true)
		if err != nil {
			return err
		}
	}

	for _, m_c := range mc {
		funcMap[m_c[0]] = function{requestFunc, mainParamType, m_c[1], channelType}
	}
	return nil
}

func handelString(fileName string, line int, s string, op1 string, op2 string, ifOrder bool) (m string, c string, err error) {
	kvList := strings.Split(s, op1)
	if !ifOrder {
		for _, kv := range kvList {
			if strings.TrimSpace(kv) == EMPTY_STRING {
				err = fmt.Errorf("@router注解缺少参数[%s：line %d]", fileName, line)
			} else {
				kvG := strings.Split(strings.TrimSpace(kv), op2)
				if len(kvG) == 2 && strings.TrimSpace(kvG[0]) != EMPTY_STRING && strings.TrimSpace(kvG[1]) != EMPTY_STRING {
					switch strings.TrimSpace(kvG[0]) {
					case "method":
						m = strings.TrimSpace(kvG[1])
					case "channel":
						c = strings.TrimSpace(kvG[1])
					}
				} else {
					err = fmt.Errorf("键值对不匹配[%s: line %d]", fileName, line)
				}
			}
		}
	} else {
		mainParam := strings.TrimSpace(kvList[0])
		if strings.TrimSpace(mainParam) != EMPTY_STRING {
			m = strings.TrimSpace(mainParam[strings.Index(mainParam, " ")+1:])
			if len(kvList) == 2 {
				channelParam := strings.TrimSpace(kvList[1])
				c = strings.TrimSpace(channelParam[strings.Index(channelParam, " ")+1:])
			}
		}
	}
	return
}

package utils

import (
	"io/ioutil"
	"net/http"
	"fmt"
	"encoding/json"
	"runtime/debug"
	"strings"
)

/**
create by cwj on j2018-05-30
used by router
 */
type RouterInterface interface {
	Handle(operation string, paramInput map[string]interface{}) (*Response, error)
}

type Handler struct {
	RootPath string
	F func() RouterInterface
}

type Response struct {
	Code int
	Json []byte
}

func NewResponse() *Response{
	return &Response{}
}

func ErrHandle()  {
	if err := recover(); err != nil {
		//v := fmt.Sprintf("ERROR!!\n%s--\n  stack \n%s", err,string(debug.Stack()))

		fmt.Printf("ERROR : %s",err)
		debug.PrintStack()

	}
}

func (c *Handler)Handle(w http.ResponseWriter, r *http.Request){
	fmt.Println(r.URL.Path)
	defer ErrHandle()
	r.ParseForm()
	var e error
	var paramInput  = make(map[string]interface{})
	if r.Method == "POST" {
		body,_ := ioutil.ReadAll(r.Body)
		e = json.Unmarshal(body, &paramInput)
		if e != nil {
			fmt.Printf("post data format error : %s", e)
			return
		}
	} else if r.Method == "GET" {
		for k, v := range r.Form{
			paramInput[k] = v[0]
		}
	} else if r.Method == "OPTIONS"{
		w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		w.Header().Set("content-type", "application/json;charset=utf-8")             //返回数据格式是json
		w.WriteHeader(200)
		paramOutput := map[string]interface{}{"status":200}
		res,_ := json.Marshal(paramOutput)
		w.Write(res)
		return
	}
	operation := strings.TrimPrefix(r.URL.Path, c.RootPath)
	var res *Response
	res, e = c.F().Handle(operation, paramInput)
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json;charset=utf-8")             //返回数据格式是json
	w.WriteHeader(res.Code)
	w.Write(res.Json)
}
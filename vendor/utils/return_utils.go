package utils

import (
	"encoding/json"
	"net/http"
)

//至多接受一个data
func ConcatSuccess(data ...interface{}) (r *Response, err error) {
	if len(data) == 0{
		return ConcatErr(OPERATE_SUCCESS, http.StatusOK, http.StatusOK, nil)
	}
	return ConcatErr(OPERATE_SUCCESS, http.StatusOK, http.StatusOK, data[0])
}

//至多接受一条msg
func ConcatFailed(data interface{}, msg ...string) (r *Response, err error) {
	switch NilCase(data, msg) {
	case 0:
		r, err = ConcatErr(OPERATE_FAILED, http.StatusInternalServerError, http.StatusInternalServerError, nil)
	case 1:
		r, err = ConcatErr(OPERATE_FAILED, http.StatusInternalServerError, http.StatusInternalServerError, data)
	case 2:
		r, err = ConcatErr(ThreeElementExpression(msg[0] != EMPTY_STRING, msg[0], OPERATE_FAILED), http.StatusInternalServerError, http.StatusInternalServerError, nil)
	case 3:
		r, err = ConcatErr(ThreeElementExpression(msg[0] != EMPTY_STRING, msg[0], OPERATE_FAILED), http.StatusInternalServerError, http.StatusInternalServerError, data)
	}
	return
}

func ConcatDeny() (r *Response, err error) {
	return ConcatErr(OPERATE_DENY, http.StatusInternalServerError, http.StatusInternalServerError, nil)
}

func ConcatNotFound() (r *Response, err error) {
	return ConcatErr(NOT_FOUND, http.StatusNotFound, http.StatusNotFound, nil)
}

func ConcatFormat() (r *Response, err error) {
	return ConcatErr(DATA_MISTAKE, http.StatusInternalServerError, http.StatusInternalServerError, nil)
}

func PackageReturn(correctData interface{}, incorrectData interface{}, errs ...error) (r *Response, err error) {
	if err = CheckError(errs...); err != nil {
		r, err = ConcatFailed(incorrectData)
	} else {
		r, err = ConcatSuccess(correctData)
	}
	return
}

func PackageReturnMsg(correctData interface{}, incorrectData interface{}, incorrectMsg string, errs ...error) (r *Response, err error) {
	if err = CheckError(errs...); err != nil {
		r, err = ConcatFailed(incorrectData, incorrectMsg)
	} else {
		r, err = ConcatSuccess(correctData)
	}
	return
}

// create by cwj 2017.8.18
func ConcatErr(msg interface{}, status int, code int, data interface{}) (r *Response, err error) {
	r = NewResponse()
	paramOutput := make(map[string]interface{})
	paramOutput["msg"] = msg
	paramOutput["status"] = status
	if data != nil {
		paramOutput["data"] = data
	}
	result, err := json.Marshal(paramOutput)
	r.Json = result
	r.Code = code

	return r, err
}

func PackageResult(res interface{}, total int, e error) (data map[string]interface{}, err error) {
	if err == nil {
		data = make(map[string]interface{})
		data["rows"] = res
		data["total"] = total
	} else {
		err = e
	}
	return
}

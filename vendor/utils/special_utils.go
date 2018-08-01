package utils

import (
	"bytes"
	"github.com/satori/go.uuid"
	"log"
	"strconv"
	"strings"
	"net"
	"fmt"
)

type Base struct {
	Page    string `json:"page"`
	Rows    string `json:"rows"`
	OrderBy string `json:"orderBy"`
}

//create by cwj on 2017-08-24
func GenerateUuid() string {
	return strings.Replace(uuid.NewV4().String(), "-", EMPTY_STRING, -1)
}

//create by cwj on 2017-10-17
//check errArray
//if there are error, return it
func CheckError(errArray ...error) error {
	for _, e := range errArray {
		if e != nil {
			log.Println(e)
			return e
		}
	}
	return nil
}

//create by cwj on 2017-10-17
//check string
//if there are Single quotation marks('),transfer it for prevent injection
func QuotationTransferred(s interface{}) string {
	str := ParseString(s)
	str = strings.Replace(str, "'", "''", -1)
	return str
}

func QuotationTransferredForLike(s interface{}) string {
	str := ParseString(s)
	str = strings.Replace(str, "'", "''", -1)
	str = strings.Replace(str, "%", "[%]", -1)
	str = strings.Replace(str, "_", "[_]", -1)
	return str
}

// 生成16位随机号
// create by cwj 2017.10.31
func GenerateNumCode16() string {
	chars := "0123456789"

	u := strings.Replace(uuid.NewV4().String(), "-", "", -1)
	var code bytes.Buffer
	for i := 0; i < 16; i++ {
		str := u[(i * 2):(i*2 + 2)]
		i1, _ := strconv.ParseInt(str, 16, 64)
		code.WriteString(string(chars[i1%10]))
	}
	return code.String()
}

// 生成16位随机号
// create by cwj 2017.10.31
func GenerateNumCode8() string {
	chars := "0123456789"

	u := strings.Replace(uuid.NewV4().String(), "-", "", -1)
	var code bytes.Buffer
	for i := 0; i < 8; i++ {
		str := u[(i * 4):(i*4 + 4)]
		i1, _ := strconv.ParseInt(str, 16, 64)
		code.WriteString(string(chars[i1%10]))
	}
	return code.String()
}

func ThreeElementExpression(isTrue bool, ele1 interface{}, ele2 interface{}) interface{}{
	if isTrue{
		return ele1
	} else {
		return ele2
	}
}

func GetLocalIp() string{
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		return "127.0.0.1"
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}
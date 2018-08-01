package main

import (
	"go-corm/corm"
	"flag"
	"os"
	"fmt"
	"go-corm/exec"
	"log"
	"time"
)


func init() {
	err := corm.Register("root:admin@tcp(localhost:3306)/test?charset=utf8", 512, 10, 8*60*60*1000)
	if nil != err {
		panic("数据库打开失败:" + err.Error())
	}

	corm.RegisterStructs(true, "order_user", new(OrderUser))

	//LogInit("log/log1")

}

func LogInit(Name string){
	logFileName := flag.String("log", Name + ".log", "Log file name")
	flag.Parse()

	//set logfile Stdout
	logFile, logErr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println(logErr)
		fmt.Println("Fail to find", *logFile, "cServer start Failed")
		os.Exit(1)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}


func main()  {
	//defer destroy()
	selectEqual("李满玲")
}

func selectEqual(args string) error{

	exe := exec.NewOrm()

	err := exe.Query("select * from order_user where enc_user_name = ?", args).ToRows(new(OrderUser))

	fmt.Println("error:" ,err)



	//var mm = []map[string]interface{}{{"AA":"test"}}
	//
	//err := exe.Query("select * from order_user where enc_user_name = ?", args).ToMap(&mm)
	//fmt.Println(err)
	//fmt.Println(mm)
	//exe.Begin()
	//_, err := exe.Exec("update order_user set enc_user_sex = ? where enc_order_user_id = ?", 5, "ididid")
	//fmt.Println(exe.Rollback())

	//o := OrderUser{OrderUser:"aaaaaaaaaa"}
	//o.UserName = "李满玲"
	//_, err := exe.Insert(&o)
	fmt.Println(err)
	return nil
}


type OrderUser struct {
	OrderUser       string `json:"id" corm:"enc_order_user_id,pk"`
	EncOrderId         string `json:"orderId"`
	EncUserPhone       string `json:"phone"`
	EncCertificateType string
	CTime			 	time.Time `json:"time" corm:"enc_time"`
	BaseOrder
}

type BaseOrder struct{
	UserName        string `corm:"enc_user_name"`
	EncCertificateNo   string
	Aaa
}

type Aaa struct {
	EncUserSex         int8
	//EncOrder           int
}
func (a *Aaa) Te(){

}

//func (c *OrderUser)New() interface{}{
//	return new(OrderUser)
//}


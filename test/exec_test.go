package main

import (
	"testing"
	"go-corm/exec"
	"fmt"
	"go-corm/corm"
	"go-corm/logs"
	"time"
)

func init() {
	err := corm.Register("root:admin@tcp(localhost:3306)/test?charset=utf8&parseTime=true", 512, 10, 8*60*60*1000)
	if nil != err {
		panic("数据库打开失败:" + err.Error())
	}

	corm.RegisterStructs(true, "order_user", new(OrderUser))
	corm.Debug(logs.WARN)

	//LogInit("log/log1")
}

func TestQueryRows(t *testing.T)  {
	exe := exec.NewOrm()

	var object []OrderUser
	err := exe.Query("select * from order_user where enc_user_name = ?", "张三").ToRows(&object)
	fmt.Println(object)
	err = exe.Query("select * from order_user where enc_user_name = ?", "张三").ToRows(&object)
	fmt.Println(object)
	_, err = exe.Exec("update order_user set enc_user_sex = ? where enc_order_user_id = ?", 5, "ididid")
	err = exe.Query("select * from order_user where enc_user_name = ?", "张三").ToRows(&object)
	fmt.Println(object)
	fmt.Println("error:" ,err)

}

func TestQueryRow(t *testing.T)  {
	exe := exec.NewOrm()

	var object OrderUser
	err := exe.Query("select * from order_user where enc_user_name = ?", "张三").ToRow(&object)

	fmt.Println("error:" ,err)
	fmt.Println(object)

}

func TestQueryMap(t *testing.T)  {
	exe := exec.NewOrm()
	var mm = []map[string]interface{}{{"AA":"test"}}

	err := exe.Query("select * from order_user where enc_user_name = ?", "张三").ToMap(&mm)
	fmt.Println(err)
	fmt.Println(mm)
}

func TestExecInsert(t *testing.T)  {
	exe := exec.NewOrm()
	sqlcon := ` INSERT INTO card_no_pool (card_no, prefix, suffix, used)
    VALUE
      (
        concat(
          ?,
          LPAD(?, 8, '0'),
          ?
        ),
        ?,
        ?,
        0
      ) ;`
      now := time.Now()
      exe.Begin()
	for i := 20001; i < 30000; i++{
		_, err := exe.Exec(sqlcon, 7788, i, 66, 7788, 66)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println(exe.Commit(), time.Now().Sub(now))
}

func TestExecUpdate(t *testing.T)  {
	exe := exec.NewOrm()
	_, err := exe.Exec("update order_user set enc_user_sex = ? where enc_order_user_id = ?", 5, "ididid")
	fmt.Println(err)
}

func TestTransaction(t *testing.T)  {
	exe := exec.NewOrm()
	exe.Begin()
	_, err := exe.Exec("update order_user set enc_user_sex = ? where enc_order_user_id = ?", 5, "ididid")
	fmt.Println(err)
	fmt.Println(exe.Commit())

}

func TestInsertObject(t *testing.T){
	exe := exec.NewOrm()
	o := OrderUser{OrderUser:"aaaaaac"}
	o.UserName = "张三"
	o.CTime = time.Now()
	_, err := exe.Insert(&o)
	fmt.Println(err)
}

func TestUpdateObject(t *testing.T){
	exe := exec.NewOrm()
	o := OrderUser{OrderUser:"aaaaaac"}
	o.UserName = "张三1"
	_, err := exe.Update(&o, "enc_user_name", "enc_")
	fmt.Println(err)
}

func TestSelectObject(t *testing.T){
	exe := exec.NewOrm()
	o := OrderUser{OrderUser:"aaaaaaaab"}
	err := exe.Read(&o)
	fmt.Println(err)
	fmt.Println(o)
}

func TestDeleteObject(t *testing.T){
	exe := exec.NewOrm()
	o := OrderUser{OrderUser:"aaaaaac"}
	o.UserName = "张三1"
	_, err := exe.Delete(&o)
	fmt.Println(err)
}

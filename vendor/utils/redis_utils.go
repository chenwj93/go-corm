package utils

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

//redis 实例
var Client *redis.Client

//create by cwj on 2017-11-11
//get lock from redis
func GetLock(key string, wait *chan bool, sec int) {
	for {
		ret := Client.SetNX(key, time.Now().UnixNano(), time.Duration(ZERO))
		if ret.Val() { //获取锁成功，直接返回
			*wait <- true
			return
		}
		value := Client.Get(key)
		if time.Now().Sub(time.Unix(ZERO_B, ParseInt64(value.Val()))) > time.Second*time.Duration(sec) { //前一个获取锁进程处理时间已经超时
			value2 := Client.GetSet(key, time.Now().UnixNano())
			if value.Val() == value2.Val() { //说明获取锁成功
				*wait <- true
				return
			}
		}
		time.Sleep(time.Microsecond * 100)
	}
}

//create by cwj on 2017-11-11
//release lock from redis
func ReleaseLock(key string, wait *chan bool) {
	<-*wait
	Client.Del(key)
}

func CheckDuplicate(key string) error {
	now := time.Now().Unix()
	fmt.Println(now)
	v := Client.Get(key)
	fmt.Println(v.Val())
	if now-ParseInt64(v.Val()) <= 10 {
		return errors.New("操作频繁！")
	}
	Client.Set(key, now, 0)
	return nil
}

func GetSerialNo(key string, expiry time.Time, length int) string {
	no := Client.Incr(key)
	if no.Val()%10 == 1 {
		Client.ExpireAt(key, expiry)
	}
	return fmt.Sprintf("%."+strconv.Itoa(length)+"d", no.Val())
}

func GetDateSerialNo(key string, length int) string {
	return GetSerialNo(key, GetToday24(), length)
}

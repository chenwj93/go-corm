package utils

import (
	"code.aliyun.com/airportcloud/common-utils/business"
	"code.aliyun.com/wyunshare/wyun-zookeeper/go-client/src/conf_center"
	"os"
	"sync"
)

var m conf_center.AppProperties
var once sync.Once

func GetConfigCenterInstance() conf_center.AppProperties {
	once.Do(func() {
		envName := GetEnvName("local_env")
		if len(envName) > 0 {
			m = business.LoadYmlDefault()
		} else {
			m = conf_center.New("airport-product")
			m.Init()
		}
	})
	return m
}

func GetEnvName(env string) string {
	return os.Getenv(env)
}

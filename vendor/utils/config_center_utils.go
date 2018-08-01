package utils

import "code.aliyun.com/wyunshare/wyun-zookeeper/go-client/src/conf_center"

var appProperties conf_center.AppProperties

func GetConfigCenterValues(data string) map[string]string {
	if appProperties.AppName == EMPTY_STRING {
		appProperties = GetConfigCenterInstance()
	}
	//conf := GetConfigCenterInstance()
	// get config
	dataMap := appProperties.ConfProperties[data]
	return dataMap
}

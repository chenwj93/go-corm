package business

import (
    "os"
    "io/ioutil"
    "strings"
    "code.aliyun.com/wyunshare/wyun-zookeeper/go-client/src/conf_center"
    "runtime"
)

func LoadYmlDefault() conf_center.AppProperties{

    file, err := os.Open("conf.yml")
    if err != nil {
        panic(err)
    }

    defer file.Close()
    content, err := ioutil.ReadAll(file)
    if err != nil {
        panic(err)
    }

    sysName := runtime.GOOS

    sep := "\n"
    if ("windows" == sysName){
        sep = "\r\n"
    }

    s :=strings.Split(string(content),sep)

    m := conf_center.New("local")

    datasource_params := make(map[string]string)
    extention_params := make(map[string]string)
    custom_params := make(map[string]string)
    protocol_params := make(map[string]string)

    for _,v := range s {
        if (len(v) > 0 && !strings.HasPrefix(v,"#") && strings.Contains(v, "=")){
            _s := strings.Split(v, "=")
            key := strings.TrimSpace(_s[0])
            va := strings.TrimSpace(_s[1])
            if(strings.HasPrefix(key, "data.source")){
                datasource_params[strings.TrimPrefix(key, "data.source.")] = va
            }
            if(strings.HasPrefix(key, "extention")){
                extention_params[strings.TrimPrefix(key, "extention.")] = va
            }
            if(strings.HasPrefix(key, "custom")){
                custom_params[strings.TrimPrefix(key, "custom.")] = va
            }
            if(strings.HasPrefix(key, "protocol")){
                protocol_params[strings.TrimPrefix(key, "protocol.")] = va
            }
        }
    }

    prop := make(map[string]map[string]string)
    prop["data.source"] = datasource_params
    prop["extention"] = extention_params
    prop["custom"] = custom_params
    prop["protocol"] = protocol_params

    m.ConfProperties = prop
    return m
}
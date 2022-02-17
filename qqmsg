package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "-------------received")
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	var mp = make(map[string]interface{})
	json.Unmarshal(body, &mp)
	//fmt.Println(string(body))
	var raw_msg string
	var IsGroup = false
	for key := range mp {
		switch key {
		case "raw_message":
			//raw_msg = *mp[key].(*string)
			//fmt.Println(mp[key])
			//fmt.Println(reflect.TypeOf(mp[key]), reflect.ValueOf(mp[key]))
			raw_msg = reflect.ValueOf(mp[key]).String()
			//fmt.Println(raw_msg)
			//LogToQQ("收到")
		case "message_type":
			if reflect.ValueOf(mp[key]).String() == "group" {
				IsGroup = true
			}
		}

	}
	//fmt.Println(raw_msg)
	if IsGroup {
		fmt.Println(raw_msg)
		go LogToQQ("我收到了你说的话 ： " + raw_msg)
	}

}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe("10.225.22.191:5701", nil)
}
func LogToQQ(info string) {
	//qq机器人 - 参考go-cqhttp框架 https://docs.go-cqhttp.org/
	//将日志 post 至 云服务器上的QQClient ，由 QQClient 发送至QQ群
	//764601511  M6项目组玩家内测QQ群号
	//696331693  Test go-cqhttp
	url := "http://127.0.0.1:5700/send_group_msg"
	s := "{\"group_id\":\"696331693\",\"message\":\"" + info + "\"}"
	resp, err := http.Post(url, "application/json; charset=utf-8", strings.NewReader(s))

	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//handle error
	}
	fmt.Println("body:" + string(body))
}

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	. "wbproject/chufangrefresh/envbuild"
	. "wbproject/chufangrefresh/server"
	. "wbproject/chufangrefresh/structure"
)

func main() {

	flag.Parse()

	db1, db2, err := EnvBuild()
	if err != nil {
		panic(err.Error())
	}

	go post("post-json")
	WebServerBase(db1, db2)

}

func get() {

	time.Sleep(1 * time.Second)

	response, err := http.Get("http://192.168.101.127:8082/getstatprogress?seq=1000")
	if err != nil {

		fmt.Println("err happens ", err)
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("in here , ", string(body))

	if response.StatusCode == 100 {
		fmt.Println("ok")
	} else {
		fmt.Println("error")
	}
}
func post(t string) {
	time.Sleep(1 * time.Second)

	switch t {
	case "post":

		data := url.Values{}
		data.Set("firstname", "foo")
		data.Add("lastname", "bar")
		//两种情况的post ..
		//1.普通的post表单请求，Content-Type=application/x-www-form-urlencoded
		//2.有文件上传的表单，Content-Type=multipart/form-data
		resp, err := http.Post("http://localhost:8082/login", "application/x-www-form-urlencoded", bytes.NewBufferString(data.Encode()))
		defer resp.Body.Close()
		if err != nil {
			fmt.Println(err)
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println("POST OK: ", string(body))
		}

	case "postform":

		resp, err := http.PostForm("http://127.0.0.1:8082/refresh",
			url.Values{"firstname": {"ruifengyun"}, "lastname": {"johnnythree"}})
		defer resp.Body.Close()
		if err != nil {
			fmt.Println(err)
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println("POST OK: ", string(body))
		}

	case "post-json":

		/*
			{
				"type":  --0 //aid  --1//gid  --2//uid
				"usrinfo":[
				{"uid":100,"st":1453046400,"et":1453046400}
				{"uid":101,"st":1453046400,"et":1453046400}
				{"uid":102,"st":1453046400,"et":1453046400}
				{"uid":103,"st":1453046400,"et":1453046400}
				]
			}
		*/
		for i := 0; i < 5; i++ {

			var s Refresh

			//s.Groupinfo = Groupinfo{Gid: 11004, St: 1452873600, Et: 1500000600}
			s.Et = 1454601600
			s.St = 1454515200
			s.Id = 220
			s.Type = 2
			//s.Id = 138
			//s.Type = 1

			b, err := json.Marshal(s)

			body := bytes.NewBuffer(b)

			resp, err := http.Post("http://192.168.101.127:8082/refresh", "application/json", body)

			defer resp.Body.Close()
			if err != nil {
				fmt.Println(err)
			} else {
				body, _ := ioutil.ReadAll(resp.Body)
				fmt.Println("POST OK: ", string(body))
			}

		}
	}

}

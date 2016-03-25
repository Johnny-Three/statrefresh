package server

import (
	//"bytes"
	"encoding/json"
	"fmt"
	//"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
)

type Refresh struct {
	Type int   `json:"type"`
	Id   int   `json:"uid"`
	St   int64 `json:"st"`
	Et   int64 `json:"et"`
}

var request_seq int
var refresh_request_map map[string]int

func init() {

	refresh_request_map = make(map[string]int)
}

func WebServerBase() {
	fmt.Println("This is webserver base!")

	//第一个参数为客户端发起http请求时的接口名，第二个参数是一个func，负责处理这个请求。
	http.HandleFunc("/refresh", refresh)

	//第一个参数为客户端发起http请求时的接口名，第二个参数是一个func，负责处理这个请求。
	http.HandleFunc("/getstatprogress", getprogress)

	//服务器要监听的主机地址和端口号
	err := http.ListenAndServe("192.168.101.127:8081", nil)

	if err != nil {
		fmt.Println("ListenAndServe error: ", err.Error())
	}
}

func getprogress(w http.ResponseWriter, r *http.Request) {

	fmt.Println("getstatprogress is running...")

	result := NewBaseJsonBean()

	//获取客户端通过GET/POST方式传递的参数
	r.ParseForm()

	if r.Method == "GET" {

		seq, found := r.Form["seq"]

		if !found {

			fmt.Println("seq is ", seq)

			result.Code = 102
			result.Message = "Get 方法访问出错，请注意参数和URL拼写.."
			bytes, _ := json.Marshal(result)
			fmt.Fprint(w, string(bytes))
			return
		}

		//todo .. 找到对应的batch_seq,从统计中发现现在对应的这个batch_seq统计量现在是什么情况，计算进度。。
		result.Code = 100
		result.Data = 90
		result.Message = "OK"
		bytes, _ := json.Marshal(result)
		fmt.Fprint(w, string(bytes))

	}

}

func refresh(w http.ResponseWriter, r *http.Request) {

	fmt.Println("refresh is running...")
	//获取客户端通过GET/POST方式传递的参数
	r.ParseForm()

	if r.Method == "POST" {

		//test json style ..
		postinfo, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			fmt.Println("in r post ", err)
			return
		}
		fmt.Printf("parse post result is %s\n", postinfo)

		var rf Refresh
		json.Unmarshal(postinfo, &rf)

		fmt.Println("=============================")
		fmt.Println(rf)

		//处理完上述操作后，将[]byte(postinfo)放至到MAP中,统计完成后清除之..
		for k, v := range refresh_request_map {

			fmt.Printf("%s ==== %s\n", k, v)

		}

		// 检查键值是否存在，如果存在则打印

		if v, ok := refresh_request_map[string(postinfo)]; ok {

			fmt.Println(v)

		} else {

			fmt.Println("Key Not Found")

		}

		request_seq += 1
		refresh_request_map[string(postinfo)] = request_seq

		result := NewBaseJsonBean()

		switch rf.Type {

		case 1:
			fmt.Println("gettype is ", rf.Type)
			//aid操作
			result.Code = 101
			result.Message = "wrong username or password .."
			break
		case 2:
			fmt.Println("gettype is ", rf.Type)
			//todo .. gid 解析
			result.Code = 100
			result.Message = "login success .."
			break

		default:
			fmt.Println(rf.Type, " is of a type I don't know how to handle")

		}

		bytes, _ := json.Marshal(result)
		fmt.Fprint(w, string(bytes))

	} else {

		//todo wrong type ..
	}
}

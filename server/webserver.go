package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	//"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync/atomic"
	. "wbproject/chufangrefresh/dbop"
	. "wbproject/chufangrefresh/logs"
	. "wbproject/chufangrefresh/structure"
	. "wbproject/chufangrefresh/util"
)

var db1, db2 *sql.DB
var request_seq int32
var refresh_request_map *BeeMap

func init() {

	refresh_request_map = NewBeeMap()
}

func WebServerBase(db01, db02 *sql.DB) {

	db1, db2 = db01, db02

	fmt.Println("This is webserver base!")
	//第一个参数为客户端发起http请求时的接口名，第二个参数是一个func，负责处理这个请求。
	http.HandleFunc("/refresh", refresh)

	//第一个参数为客户端发起http请求时的接口名，第二个参数是一个func，负责处理这个请求。
	http.HandleFunc("/getstatprogress", getprogress)

	//服务器要监听的主机地址和端口号
	err := http.ListenAndServe("192.168.101.127:8082", nil)

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
		fmt.Println("seq is ", seq)
		fmt.Println("refresh_request_map", refresh_request_map)

		var seqin int
		seqin, err := strconv.Atoi(seq[0])

		if err != nil {

			//fmt.Println("seq is ", seq)
			result.Code = 102
			result.Message = "seq转换int出错，请注意拼写"
			bytes, _ := json.Marshal(result)
			fmt.Fprint(w, string(bytes))
			return

		}

		//find uploadid in Deal_status_map ..
		key := refresh_request_map.GetByValue(seqin)
		if key == nil {

			result.Code = 100
			result.Message = "统计完毕"
			result.Data = 100
			bytes, _ := json.Marshal(result)
			fmt.Fprint(w, string(bytes))
			return
		}

		fmt.Println("Deal_status_map ", Deal_status_map)

		uploadid := Deal_status_map.Get(key)

		//找到对应的key
		if uploadid != nil {

			ifexist := SelectUploadid(db2, uploadid.(int))
			//找到，说明尚未处理；未找到，说明已经处理完毕
			if ifexist == true {

				result.Code = 100
				result.Message = "统计中"
				result.Data = 0
				bytes, _ := json.Marshal(result)
				fmt.Fprint(w, string(bytes))

				return

			} else {

				refresh_request_map.Delete(key)

				result.Code = 100
				result.Message = "统计完毕"
				result.Data = 100
				bytes, _ := json.Marshal(result)
				fmt.Fprint(w, string(bytes))

				return

			}

		} else {
			//todo.. impossible happens ..
		}

		return

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

		result := NewBaseJsonBean()

		var rf Refresh
		json.Unmarshal(postinfo, &rf)

		if rf.Et < rf.St {

			result.Code = 102
			result.Message = "结束时间不可能小于开始时间"
			bytes, _ := json.Marshal(result)
			fmt.Fprint(w, string(bytes))

			return

		}

		// 检查键值是否存在，如果存在则打印
		ifexist := refresh_request_map.Check(string(postinfo))

		if true == ifexist {

			result.Code = 103
			result.Data = refresh_request_map.Get(string(postinfo))
			result.Message = "任务已经在处理，请勿重复发起"

			bytes, _ := json.Marshal(result)
			fmt.Fprint(w, string(bytes))

			return
		}

		//键值不存在，说明是新的任务，需要受理
		atomic.AddInt32(&request_seq, 1)
		ifsuc := refresh_request_map.Set(string(postinfo), int(request_seq))

		if !ifsuc {

			Logger.Critical("chufangRefresh break down cause refresh_request_map.Set wrong ..")
			panic("refresh_request_map.Set wrong ..")
		}

		result.Code = 100
		result.Data = request_seq
		result.Message = "OK"

		bytes, _ := json.Marshal(result)
		fmt.Fprint(w, string(bytes))

		go InsertQueue(&rf, string(postinfo), db1, db2)
	}
}

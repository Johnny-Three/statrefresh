package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	//"net/url"
	//"strings"
	"bytes"
)

func main() {

	url := "http://192.168.101.127:8082/refresh"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"xx":"xxxx","title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	fmt.Println(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}

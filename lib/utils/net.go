package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// 发送 json格式请求
func SendJsonPost(requestBody string, url string, method string) (body []byte, isOk bool) {
	var jsonStr = []byte(requestBody)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	isOk = true
	if err != nil {
		isOk = false
		return
	}
	defer resp.Body.Close()
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	//fmt.Println("response Body:", string(body))
	body, _ = ioutil.ReadAll(resp.Body)
	return
}
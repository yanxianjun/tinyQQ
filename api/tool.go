package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func TikTok(url string) (msg string) {
	if "" == url {
		msg = "请输入要解析的url！"
		return
	}

	res, err := sendGetReq(url)
	if err != nil {
		msg = "解析失败！"
		return
	}
	id := strings.Split(res.Header.Get("Location"), "/")[5]
	// 通过接口获得视频的详细内容
	relUrl := "https://www.iesdouyin.com/web/api/v2/aweme/iteminfo/?item_ids=" + id
	// 获取到视频无水印的播放地址
	res, err = sendGetReq(relUrl)
	body, err := ioutil.ReadAll(res.Body)
	var tempMap map[string]interface{}
	_ = json.Unmarshal(body, &tempMap)
	a := tempMap["item_list"].([]interface{})[0].(map[string]interface{})["video"].(map[string]interface{})["play_addr"].(map[string]interface{})["url_list"].([]interface{})[0].(string)
	relUrl = strings.Replace(a, "playwm", "play", -1)
	//获取真实的视频url
	res, err = sendGetReq(relUrl)
	relUrl = res.Header.Get("location")
	msg = relUrl
	return
}

func sendGetReq(url string) (res *http.Response, err error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Referer", "[url=https://v.douyin.com/]https://v.douyin.com/[/url]")
	req.Header.Set("User-Agent","Mozilla/5.0 (Linux; Android 8.0.0; Pixel 2 XL Build/OPD1.170816.004) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Mobile Safari/537.36")
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	res, err = client.Do(req)
	return
}
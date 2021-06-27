package help

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"tinyQQ/lib/bot"
	"tinyQQ/lib/message"
	"tinyQQ/lib/utils"
)

var (
	matchReg = regexp.MustCompile(`\[CQ:\w+?.*?]`)
	typeReg  = regexp.MustCompile(`\[CQ:(\w+)`)
	paramReg = regexp.MustCompile(`,([\w\-.]+?)=([^,\]]+)`)
)

/**
 * 解析消息
 */
func ParseMessage(str string, qq int64, group int64) (im *message.SendingMessage, err error) {
	var r []message.IMessageElement
	i := matchReg.FindAllStringSubmatchIndex(str, -1)
	if len(i) == 0 {
		im = message.NewSendingMessage().Append(message.NewText(str))
		return
	}
	si := 0
	for _, idx := range i {
		// 第一个CQ出现前的文本
		if idx[0] > si {
			text := str[si:idx[0]]
			r = append(r, message.NewText(CQCodeUnescapeText(text)))
		}
		code := str[idx[0]:idx[1]]
		si = idx[1]
		// CQ的类型
		t := typeReg.FindAllStringSubmatch(code, -1)[0][1]
		// CQ的键值
		ps := paramReg.FindAllStringSubmatch(code, -1)
		// 用于存放键值
		d := make(map[string]string)
		for _, p := range ps {
			d[p[1]] = CQCodeUnescapeValue(p[2])
		}
		// 组建消息元素
		element, errB := buildElement(t, d, qq, group)
		err = errB
		// 组装到桶中
		r = append(r, element)
	}
	// 最后一个CQ出现前的文本
	if len(i) > 0 {
		endText := str[i[len(i)-1][1]:]
		r = append(r, message.NewText(CQCodeUnescapeText(endText)))
	}
	// 声明消息体
	im = message.NewSendingMessage()
	// 将桶中消息元素方式消息体中
	for _, v := range r {
		im.Append(v)
	}
	return
}

func buildElement(t string, d map[string]string, qq int64, group int64) (r message.IMessageElement, err error) {
	switch t {
	case "image":
		if _, ok := d["url"]; !ok {

			return r, err
		}
		bytes, err := utils.HttpGetBytes(d["url"], "")
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		if qq != 0 {
			r, err = bot.Instance.UploadPrivateImage(qq, bytes)
		} else {
			r, err = bot.Instance.UploadGroupImage(group, bytes)
		}
		return r, err
	case "face":
		if _, ok := d["id"]; !ok {
			return r, err
		}
		faceId, err := strconv.ParseInt(d["id"], 10, 32)
		if err != nil {
			return nil, err
		}
		r = message.NewFace(int32(faceId))
		return r, err
	case "at":
		if _, ok := d["qq"]; !ok {
			return r, err
		}
		qq, err := strconv.ParseInt(d["qq"], 10, 64)
		if qq != 0 {
			r := message.NewAt(qq, d["content"])
			return r, err
		} else {
			r := message.AtAll()
			return r, err
		}
	default:
		return nil, nil
	}
	return
}

func CQCodeUnescapeText(content string) string {
	ret := content
	ret = strings.ReplaceAll(ret, "&#91;", "[")
	ret = strings.ReplaceAll(ret, "&#93;", "]")
	ret = strings.ReplaceAll(ret, "&amp;", "&")
	return ret
}

func CQCodeUnescapeValue(content string) string {
	ret := strings.ReplaceAll(content, "&#44;", ",")
	ret = CQCodeUnescapeText(ret)
	return ret
}

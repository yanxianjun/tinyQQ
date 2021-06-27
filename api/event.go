package api

import (
	"encoding/json"
	"fmt"
	"tinyQQ/cache"
	"tinyQQ/config"
	"tinyQQ/lib/client"
	"tinyQQ/lib/message"
	"tinyQQ/lib/utils"
	"tinyQQ/logs"
)

type (
	PrivateMessage struct {
		Id         int32
		InternalId int32
		Target     int64
		Time       int32
		Sender     Sender
		Elements   []IMessageElement
	}

	GroupMessage struct {
		Id         int32
		InternalId int32
		GroupCode  int64
		GroupName  string
		Sender     Sender
		Time       int32
		Content    string
	}

	GroupMemberJoined struct {
		Uin                    int64
		Gender                 byte
		Nickname               string
		CardName               string
		Level                  uint16
		JoinTime               int64
		LastSpeakTime          int64
		SpecialTitle           string
		SpecialTitleExpireTime int64
	}

	EventType struct {
		Type string
		Msg interface{} // GroupMessage
	}

	Sender struct {
		Uin      int64
		Nickname string
		CardName string
		IsFriend bool
	}

	IMessageElement interface {
		Type() ElementType
	}

	ElementType int
)

// 监听私聊消息
func OnPrivateMessage(client *client.QQClient, privateMessage *message.PrivateMessage){
	logs.Logger.Infof("收到好友[%v]新消息,内容:%s", privateMessage.Sender.Nickname, privateMessage.ToString())

	text := privateMessage.ToString()
	if len(text) > 21 && text[0:21] == "解析抖音水印：" {
		tok := TikTok(text[21:])
		go client.SendPrivateMessage(privateMessage.Sender.Uin, message.NewSendingMessage().Append(message.NewText(tok)))
	}
	if len(text) > 19 && text[0:19] == "解析抖音水印:" {
		tok := TikTok(text[19:])
		go client.SendPrivateMessage(privateMessage.Sender.Uin, message.NewSendingMessage().Append(message.NewText(tok)))
	}

	var msg PrivateMessage
	msg.Id = privateMessage.Id
	msg.InternalId = privateMessage.InternalId
	msg.Target = privateMessage.Target
	msg.Time = privateMessage.Time
	msg.Sender.Uin = privateMessage.Sender.Uin
	msg.Sender.Nickname = privateMessage.Sender.Nickname
	msg.Sender.CardName = privateMessage.Sender.CardName
	msg.Sender.Uin = privateMessage.Sender.Uin
	msg.Sender.IsFriend = privateMessage.Sender.IsFriend

	buf, _ := json.Marshal(msg)
	url := config.GlobalConfig.GetString("event.url")
	_, ok := utils.SendJsonPost(string(buf), url, "POST")
	logs.Logger.Infof("私聊事件上报%v", ok)
}

// 监听群聊消息
func OnGroupMessage(client *client.QQClient, message *message.GroupMessage){
	logs.Logger.Infof("收到群[%v]新消息,内容:%s", message.GroupName, message.ToString())

	var msg GroupMessage
	var typ EventType
	msg.Id = message.Id
	msg.Time = message.Time
	msg.InternalId = message.InternalId
	msg.GroupCode = message.GroupCode
	msg.GroupName = message.GroupName
	msg.Sender.CardName = message.Sender.CardName
	msg.Sender.Nickname = message.Sender.Nickname
	msg.Sender.IsFriend = message.Sender.IsFriend
	msg.Sender.Uin = message.Sender.Uin
	msg.Content = message.ToString()
	typ.Type = "GroupMessage"
	typ.Msg = msg
	//group := client.FindGroupByUin(message.GroupCode)
	group := client.FindGroup(message.GroupCode)
	if group == nil {
		return
	}
	buf, _ := json.Marshal(typ)
	var userInfo struct{
		Callback string `gorm:"column:callback"`
	}
	cache.DB.Debug().Table("tb_users").Where("`qq` = ?", group.OwnerUin).First(&userInfo)
	if userInfo.Callback != "" {
		utils.SendJsonPost(string(buf), userInfo.Callback, "POST")
	}

	//url := config.GlobalConfig.GetString("event.url")
	//_, ok := utils.SendJsonPost(string(buf), url, "POST")
	//logs.Logger.Infof("群消息事件上报%v", ok)
}

func OnNewFriendAdded(client *client.QQClient, event *client.NewFriendEvent) {
	//fmt.Println(event.Friend)
}

var NewFriendRequest map[int64]interface{}
var NewGroupInvited map[int64]interface{}

func OnNewFriendRequest(client *client.QQClient, request *client.NewFriendRequest) {
	logs.Logger.Infof("收到新好友[%v]添加请求,请求内容:%s", request.RequesterNick, request.Message)
	var msg struct {
		RequestId     int64
		Message       string
		RequesterUin  int64
		RequesterNick string
	}
	msg.RequestId = request.RequestId
	msg.RequesterUin = request.RequesterUin
	msg.RequesterNick = request.RequesterNick
	msg.Message = request.Message

	if "auto" == config.GlobalConfig.GetString("event.addFriend") {
		request.Accept()
		return
	} else {
		NewFriendRequest = make(map[int64]interface{})
		NewFriendRequest[msg.RequestId] = request

		buf, _ := json.Marshal(msg)
		url := config.GlobalConfig.GetString("event.url")
		_, ok := utils.SendJsonPost(string(buf), url, "POST")
		logs.Logger.Infof("新好友请求事件上报%v", ok)
		return
	}
}

func OnGroupInvited(client *client.QQClient, request *client.GroupInvitedRequest) {
	logs.Logger.Infof("收到新群[%v]邀请加入", request.GroupName)
	var msg struct {
		RequestId     int64
		GroupName       string
		GroupCode  int64
		InvitorNick string
		InvitorUin int64
	}
	msg.RequestId = request.RequestId
	msg.GroupName = request.GroupName
	msg.GroupCode = request.GroupCode
	msg.InvitorNick = request.InvitorNick
	msg.InvitorUin = request.InvitorUin

	if "auto" == config.GlobalConfig.GetString("event.groupInvited") {
		request.Accept()
		return
	} else {
		NewGroupInvited = make(map[int64]interface{})
		NewGroupInvited[msg.RequestId] = request

		buf, _ := json.Marshal(msg)
		url := config.GlobalConfig.GetString("event.url")
		_, ok := utils.SendJsonPost(string(buf), url, "POST")
		logs.Logger.Infof("新群邀请事件上报%v", ok)
		return
	}
}

func OnGroupMuted(client *client.QQClient, event *client.GroupMuteEvent) {
	//event.
}

func OnGroupMemberJoined(client *client.QQClient, event *client.MemberJoinGroupEvent) {
	logs.Logger.Infof("群[%v]有新成员%v加入了", event.Group.Name, event.Member.Nickname)

	var msg GroupMemberJoined
	var typ EventType

	msg.Nickname = event.Member.Nickname
	msg.CardName = event.Member.CardName
	msg.SpecialTitle = event.Member.SpecialTitle
	msg.Uin = event.Member.Uin
	msg.JoinTime = event.Member.JoinTime
	msg.LastSpeakTime = event.Member.LastSpeakTime
	msg.SpecialTitleExpireTime = event.Member.SpecialTitleExpireTime
	msg.Gender = event.Member.Gender
	msg.Level = event.Member.Level
	typ.Type = "GroupMemberJoined"
	typ.Msg = msg

	group := client.FindGroup(event.Group.Code)
	if group == nil {
		return
	}
	buf, _ := json.Marshal(typ)
	var userInfo struct{
		Callback string `gorm:"column:callback"`
	}
	cache.DB.Table("tb_users").Where("`qq` = ?", group.OwnerUin).First(&userInfo)
	if userInfo.Callback != "" {
		utils.SendJsonPost(string(buf), userInfo.Callback, "POST")
	}
}

// 监听系统日志
func OnLog(client *client.QQClient, event *client.LogEvent) {
	fmt.Println(event.Message)
}
package api

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"tinyQQ/help"
	"tinyQQ/lib/bot"
	"tinyQQ/lib/client"
	"tinyQQ/lib/message"
	"tinyQQ/logs"
)

// 发送私聊消息
func SendPrivateMessage(c *gin.Context) {
	var msg struct {
		QQ      int64  `json:"qq" form:"qq"`
		Content string `json:"content" form:"content"`
		CQ      bool   `json:"cq" form:"cq"`
	}
	if err := c.ShouldBind(&msg); err != nil {
		logs.Logger.WithError(err).Warn("参数绑定错误！")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  -1,
			"error": err.Error(),
		})
		return
	}

	if msg.QQ == 0 || msg.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  -1,
			"error": "参数不正确！",
		})
		return
	}

	err := bot.Instance.ReloadFriendList()
	if err != nil {
		logs.Logger.WithError(err).Panic("刷新好友列表错误！")
	}

	if friend := bot.Instance.FindFriend(msg.QQ); friend == nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "请先添加好友！",
		})
		return
	}

	// 纯文本消息
	if !msg.CQ {
		bot.Instance.SendPrivateMessage(msg.QQ, message.NewSendingMessage().Append(message.NewText(msg.Content)))
	} else {
		im, err := help.ParseMessage(msg.Content, msg.QQ, 0)
		if err != nil {
			logs.Logger.WithError(err).Panic("解析CQ码错误！")
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  -1,
				"error": "解析CQ码错误！",
			})
			return
		}
		bot.Instance.SendPrivateMessage(msg.QQ, im)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
	})
}

// 发送群聊消息
func SendGroupMessage(c *gin.Context) {
	var msg struct {
		Group   int64  `json:"group" form:"group"`
		Content string `json:"content" form:"content"`
		CQ      bool   `json:"cq" form:"cq"`
	}
	if err := c.ShouldBind(&msg); err != nil {
		logs.Logger.WithError(err).Warn("参数绑定错误！")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  -1,
			"error": err.Error(),
		})
		return
	}
	// 这里执行时间过长
	//err := bot.Instance.ReloadGroupList()
	//if err != nil {
	//	logs.Logger.WithError(err).Panic("刷新群列表错误！")
	//}
	if friend := bot.Instance.FindGroup(msg.Group); friend == nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "请先添加群！",
		})
		return
	}

	// 纯文本消息
	if !msg.CQ {
		bot.Instance.SendGroupMessage(msg.Group, message.NewSendingMessage().Append(message.NewText(msg.Content)))
	} else {
		im, err := help.ParseMessage(msg.Content, 0, msg.Group)
		if err != nil {
			logs.Logger.WithError(err).Panic("解析CQ码错误！")
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  -1,
				"error": "解析CQ码错误！",
			})
			return
		}
		bot.Instance.SendGroupMessage(msg.Group, im)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
	})
}

// 获取好友列表
func GetFriendList(c *gin.Context) {
	list, err := bot.Instance.GetFriendList()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  -1,
			"error": err.Error(),
		})
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"code": 200,
			"list": list,
		},
	)
}

// 处理好友请求
func SetNewFriendRequest(c *gin.Context) {
	var msg struct {
		ID     int64 `json:"id" form:"id"`
		Accept bool  `json:"accept" form:"accept"`
	}
	if err := c.ShouldBind(&msg); err != nil {
		logs.Logger.WithError(err).Warn("参数绑定错误！")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  -1,
			"error": err.Error(),
		})
		return
	}
	request := NewFriendRequest[msg.ID].(*client.NewFriendRequest)
	if msg.Accept {
		request.Accept()
	} else {
		request.Reject()
	}
}

// 设置群禁言
func SetGroupMuteAll(c *gin.Context) {
	var param struct {
		GroupId  int64 `json:"group_id" form:"group_id"`
		Mute bool  `json:"mute" form:"mute"`
	}
	if err := c.ShouldBind(&param); err != nil {
		logs.Logger.WithError(err).Warn("参数绑定错误！")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  -1,
			"error": err.Error(),
		})
		return
	}
	group := bot.Instance.FindGroupByUin(param.GroupId)
	if group == nil {
		return
	}
	group.MuteAll(param.Mute)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
	})
}

// 设置群成员禁言指定时间
func SetGroupMute(c *gin.Context) {
	var param struct {
		GroupId  int64 `json:"group_id" form:"group_id"`
		GroupMember int64 `json:"group_member" form:"group_member"`
		MuteTime uint32  `json:"mute_time" form:"mute_time"`
	}
	if err := c.ShouldBind(&param); err != nil {
		logs.Logger.WithError(err).Warn("参数绑定错误！")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  -1,
			"error": err.Error(),
		})
		return
	}
	group := bot.Instance.FindGroupByUin(param.GroupId)
	if group == nil {
		return
	}
	member := group.FindMember(param.GroupMember)
	if member == nil {
		return
	}
	member.Mute(param.MuteTime)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
	})
}

// 设置群公告
func SetGroupMemo(c *gin.Context) {
	var param struct {
		GroupId  int64 `json:"group_id" form:"group_id"`
		Memo string  `json:"memo" form:"memo"`
	}
	if err := c.ShouldBind(&param); err != nil {
		logs.Logger.WithError(err).Warn("参数绑定错误！")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  -1,
			"error": err.Error(),
		})
		return
	}
	group := bot.Instance.FindGroupByUin(param.GroupId)
	if group == nil {
		return
	}
	group.UpdateMemo(param.Memo)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
	})
}

// 设置群名
func SetGroupName(c *gin.Context) {
	var param struct {
		GroupId  int64 `json:"group_id" form:"group_id"`
		GroupName string  `json:"group_name" form:"group_name"`
	}
	if err := c.ShouldBind(&param); err != nil {
		logs.Logger.WithError(err).Warn("参数绑定错误！")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  -1,
			"error": err.Error(),
		})
		return
	}
	group := bot.Instance.FindGroupByUin(param.GroupId)
	if group == nil {
		return
	}
	group.UpdateName(param.GroupName)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
	})
}

func Test(c *gin.Context) {
	Restart <- struct{}{}
}

var Restart = make(chan struct{}, 1)

func RRestart(Args []string) {
	cmd := &exec.Cmd{}
	if runtime.GOOS == "windows" {
		file, err := exec.LookPath(Args[0])
		if err != nil {
			log.Errorf("重启失败:%s", err.Error())
			return
		}
		path, err := filepath.Abs(file)
		if err != nil {
			log.Errorf("重启失败:%s", err.Error())
		}
		Args = append([]string{"/c", "start ", path, "faststart"}, Args[1:]...)
		cmd = &exec.Cmd{
			Path:   "cmd.exe",
			Args:   Args,
			Stderr: os.Stderr,
			Stdout: os.Stdout,
		}
	} else {
		Args = append(Args, "faststart")
		cmd = &exec.Cmd{
			Path:   Args[0],
			Args:   Args,
			Stderr: os.Stderr,
			Stdout: os.Stdout,
		}
	}
	cmd.Start()
}
package server

import (
	"github.com/gin-gonic/gin"
	"tinyQQ/api"
	"tinyQQ/logs"
)

func RunHttp() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	router(r)
	logs.Logger.Info("HTTP Server Start successfully！http://127.0.0.1:5701")
	r.Run(":5701")
}

func router(r *gin.Engine) {
	path := r.Group("")
	path.Any("send_private_msg", api.SendPrivateMessage)
	path.Any("send_group_msg", api.SendGroupMessage)
	path.Any("get_friend_list", api.GetFriendList)
	path.Any("set_new_friend_request", api.SetNewFriendRequest)
	path.Any("set_group_mute_all", api.SetGroupMuteAll)
	path.Any("set_group_mute", api.SetGroupMute)
	path.Any("set_group_memo", api.SetGroupMemo)
	path.Any("set_group_name", api.SetGroupName)
	//path.Any("test", api.Test)

}

package server

import (
	"tinyQQ/api"
	"tinyQQ/lib/bot"
	"tinyQQ/logs"
)

func Event()  {
	bot.Instance.OnPrivateMessage(api.OnPrivateMessage)
	bot.Instance.OnGroupMessage(api.OnGroupMessage)
	bot.Instance.OnNewFriendAdded(api.OnNewFriendAdded)
	bot.Instance.OnNewFriendRequest(api.OnNewFriendRequest)
	bot.Instance.OnGroupInvited(api.OnGroupInvited)
	bot.Instance.OnGroupMuted(api.OnGroupMuted)
	bot.Instance.OnGroupMemberJoined(api.OnGroupMemberJoined)
	//bot.Instance.OnLog(api.OnLog)
	logs.Logger.Info("Event Server Start successfully！")
}
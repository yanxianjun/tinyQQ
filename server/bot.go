package server

import (
	"os"
	"os/signal"
	"tinyQQ/lib/bot"
	"tinyQQ/logs"
)

func ServerBefore () chan os.Signal {
	logs.Logger.Info("System initialization, please wait ...")
	bot.Init() // 快速初始化
	bot.StartService() // 初始化 Modules
	bot.UseProtocol(bot.AndroidPhone) // 使用协议 不同协议可能会有部分功能无法使用 在登陆前切换协议
	bot.Login() // 登录
	bot.RefreshList() // 刷新好友列表，群列表
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	return ch
}

func ServerAfter (ch chan os.Signal) {
	<-ch
	bot.Stop()
}
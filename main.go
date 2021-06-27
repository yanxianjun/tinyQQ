package main

import (
	"os"
	"tinyQQ/api"
	"tinyQQ/config"
	"tinyQQ/lib/bot"
	"tinyQQ/lib/utils"
	"tinyQQ/logs"
	"tinyQQ/server"
)

func init() {
	utils.WriteLogToFS()
	config.Init()
}

func main() {
	//go func() {
	//	time.Sleep(10*time.Second)
	//	os.Exit(-2)
	//}()

	ch := server.ServerBefore ()

	//if err := cache.InitMysqlCon(); err != nil {
	//	panic("数据库连接失败！")
	//}

	//go server.RunHttp()

	//if "" != config.GlobalConfig.GetString("event.url") {
	//	go server.Event()
	//}

	r := api.Restart
	arg := os.Args
	select {
	case <-r:
		logs.Logger.Info("正在重启中...")
		api.RRestart(arg)
	case <-ch:
		bot.Stop()
	}

	//server.ServerAfter(ch)
}
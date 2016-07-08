package main

import (
	_ "digs/docs"
	_ "digs/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/astaxie/beego/logs"
	"github.com/afex/hystrix-go/hystrix"
	"net"
	"net/http"
	"digs/common"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	//Configure routes
	beego.InsertFilter("*", beego.BeforeRouter,cors.Allow(&cors.Options{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "GET", "OPTIONS"},
		AllowHeaders: []string{"Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
		AllowCredentials: true,
	}))

	//Configure logger
	log := logs.NewLogger(1)
	log.EnableFuncCallDepth(true)
	log.SetLogger("multifile", `{"filename":"test.log","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)

	//Configure hystrix/monitoring
	setCommandParameters()
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "8084"), hystrixStreamHandler)

	beego.Run()
}

func setCommandParameters() {
	singleCommandConfig := hystrix.CommandConfig{
		Timeout:                1000,
		MaxConcurrentRequests:  5,
	}
	batchCommandConfig := hystrix.CommandConfig{
		Timeout:                5000,
		MaxConcurrentRequests:  5,
	}
	singleCommandHighConcurrencyConfig := hystrix.CommandConfig{
		Timeout:                1000,
		MaxConcurrentRequests:  50,
	}
	extrnalUnthrottled := hystrix.CommandConfig{
		Timeout:                1000,
		MaxConcurrentRequests:  10000,
	}
	commandConfigMap := make(map[string]hystrix.CommandConfig)
	commandConfigMap[common.MessageWrite] = singleCommandConfig
	commandConfigMap[common.MessageGetAll] = batchCommandConfig
	commandConfigMap[common.Notification] = singleCommandConfig
	commandConfigMap[common.UserAccount] = singleCommandConfig
	commandConfigMap[common.UserAccountGetAll] = batchCommandConfig
	commandConfigMap[common.SessionWrite] = singleCommandConfig
	commandConfigMap[common.SessionGet] = singleCommandHighConcurrencyConfig
	commandConfigMap[common.SessionDel] = singleCommandConfig
	commandConfigMap[common.FeedAdd] = singleCommandHighConcurrencyConfig
	commandConfigMap[common.FeedGet] = singleCommandConfig
	commandConfigMap[common.LocationGet] = singleCommandHighConcurrencyConfig
	commandConfigMap[common.LocationUpdate] = singleCommandHighConcurrencyConfig
	commandConfigMap[common.LocationUserFind] = batchCommandConfig
	commandConfigMap[common.AndroidPush] = extrnalUnthrottled
	commandConfigMap[common.IOSPush] = extrnalUnthrottled

	hystrix.Configure(commandConfigMap)

}

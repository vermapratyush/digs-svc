package main

import (
	_ "digs/docs"
	_ "digs/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/afex/hystrix-go/hystrix"
	"net"
	"net/http"
	"digs/common"
	"digs/logger"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	//Configure routes
	beego.InsertFilter("*", beego.BeforeRouter,cors.Allow(&cors.Options{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
		AllowCredentials: true,
	}))

	logger.Initialize()

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
	externalUnthrottled := hystrix.CommandConfig{
		Timeout:                1000,
		MaxConcurrentRequests:  10000,
		ErrorPercentThreshold:  101,
	}
	awsS3ExternalUnthrottled := hystrix.CommandConfig{
		Timeout:                5000,
		MaxConcurrentRequests:  10000,
		ErrorPercentThreshold:  101,
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
	commandConfigMap[common.FeedDel] = singleCommandConfig
	commandConfigMap[common.LocationGet] = singleCommandHighConcurrencyConfig
	commandConfigMap[common.LocationUpdate] = singleCommandHighConcurrencyConfig
	commandConfigMap[common.LocationUserFind] = batchCommandConfig
	commandConfigMap[common.UserGroup] = singleCommandHighConcurrencyConfig
	commandConfigMap[common.UserGroupBatch] = batchCommandConfig
	commandConfigMap[common.AndroidPush] = externalUnthrottled
	commandConfigMap[common.IOSPush] = externalUnthrottled
	commandConfigMap[common.MeetupAPI] = externalUnthrottled
	commandConfigMap[common.FourSquareAPI] = externalUnthrottled
	commandConfigMap[common.BitlyAPI] = externalUnthrottled
	commandConfigMap[common.AmazonS3] = awsS3ExternalUnthrottled

	hystrix.Configure(commandConfigMap)

}

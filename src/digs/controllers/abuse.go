package controllers

import "github.com/astaxie/beego"

type AbuseController struct {
	HttpBaseController
}

func (this *AbuseController) Post() {
	beego.Critical("ReportAbuse|Post=", string(this.Ctx.Input.RequestBody))
	this.Serve204()
}
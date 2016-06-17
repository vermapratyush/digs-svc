package controllers

import (
)
import (
	"digs/domain"
	"github.com/astaxie/beego"
	"encoding/json"
)

type HttpBaseController struct {
	beego.Controller
}

type WSBaseController struct {

}

func (this *HttpBaseController) Super(request *domain.BaseRequest) *HttpBaseController {
	if this.Ctx.Input.Method() == "POST" {
		json.Unmarshal(this.Ctx.Input.RequestBody, request)
	}
	request.SessionID = this.Ctx.Input.Header("SID")
	request.UserAgent = this.Ctx.Input.UserAgent()
	return this
}


type WebSocketController struct {
	beego.Controller
}

func (this *HttpBaseController) Serve500(err error) {
	this.Data["json"] = domain.ErrorResponse{
		StatusCode:500,
		Message:err.Error(),
	}
	this.Ctx.Output.SetStatus(500)
	this.ServeJSON()
	return

}

func (this *HttpBaseController) Serve200(obj interface{}) {
	this.Data["json"] = obj
	this.Ctx.Output.SetStatus(200)
	this.ServeJSON()
}

func (this *HttpBaseController) Serve304() {
	this.Ctx.Output.SetStatus(304)
	this.ServeJSON()
}

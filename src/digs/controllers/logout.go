package controllers

import (
	"digs/domain"
	"github.com/astaxie/beego"
	"encoding/json"
	"digs/models"
)

type LogoutController struct {
	HttpBaseController
}

func (this *LogoutController) Post()  {
	var request domain.UserLogoutRequest
	beego.Info(string(this.Ctx.Input.RequestBody))
	this.Super(&request.BaseRequest)
	json.Unmarshal(this.Ctx.Input.RequestBody, &request)

	err := models.DeleteUserAuth(request.SessionID)
	if err != nil {
		beego.Info(err)
		this.Serve500(err)
		return
	}
	beego.Info("304")
	this.Serve204()
}


package controllers

import (
	"digs/domain"
	"github.com/astaxie/beego"
	"encoding/json"
	"digs/models"
	"gopkg.in/mgo.v2"
)

type LogoutController struct {
	HttpBaseController
}

func (this *LogoutController) Post()  {
	var request domain.UserLogoutRequest
	beego.Info(string(this.Ctx.Input.RequestBody))
	this.Super(&request.BaseRequest)
	json.Unmarshal(this.Ctx.Input.RequestBody, &request)

	userAuth, err := models.FindSession("sid", request.SessionID)
	beego.Info(request.SessionID)
	if err != nil {
		beego.Error(err)
		if err == mgo.ErrNotFound {
			this.InvalidSessionResponse()
			return
		}
		this.Serve500(err)
		return
	}

	err = models.DeleteNotificationId(userAuth.UID, request.NotificationId)

	if err != nil {
		beego.Error("NotificationDeleteFailed|notificationId=", request.NotificationId, "|err=", err)
	}

	err = models.DeleteUserAuth(userAuth.Id)
	if err != nil {
		beego.Info(err)
		this.Serve500(err)
		return
	}
	beego.Info("304")
	this.Serve204()
}


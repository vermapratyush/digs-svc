package controllers

import (
	"digs/domain"
	"github.com/astaxie/beego"
	"encoding/json"
	"digs/models"
	"digs/socket"
	"errors"
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
		this.Serve500(errors.New("Unable to find session"))
		return
	}
	go socket.LeaveNode(userAuth.UID)

	err = models.DeleteNotification(request.NotificationId)

	if err != nil {
		beego.Error("NotificationDeleteFailed|notificationId=", request.NotificationId)
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


package controllers

import (
	"digs/domain"
	"github.com/astaxie/beego"
	"encoding/json"
	"digs/models"
	"errors"
)

type NotificationController struct {
	HttpBaseController
}

func (this *NotificationController) Post() {
	var request domain.NotificationRequest
	beego.Info(string(this.Ctx.Input.RequestBody))
	this.Super(&request.BaseRequest)
	json.Unmarshal(this.Ctx.Input.RequestBody, &request)

	session, err := models.FindSession("sid", request.SessionID)
	if err != nil {
		this.Serve500(errors.New("Unable to find session"))
		return
	}
	err = models.AddNotificationId(session.UID, request.NotificationID, request.OSType)
	if err != nil {
		this.Serve500(errors.New("Unable to register device"))
	} else {
		this.Serve204()
	}

}

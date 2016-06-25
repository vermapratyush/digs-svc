package controllers

import (
	"digs/domain"
	"github.com/astaxie/beego"
	"encoding/json"
	"digs/models"
	"errors"
)

type SettingController struct {
	HttpBaseController
}

func (this *SettingController) Post() {
	var request domain.SettingRequest
	beego.Info("Login Request", string(this.Ctx.Input.RequestBody))
	this.Super(&request.BaseRequest)
	json.Unmarshal(this.Ctx.Input.RequestBody, &request)

	userAuth, err := models.FindSession("sid", request.SessionID)
	if err != nil {
		beego.Error("SessionNotFound|err=", err)
		this.Serve500(errors.New("Unable to find session"))
		return
	}
	err = models.UpdateUserAccount(userAuth.UID, request.PushNotification, request.Range)

	if err != nil {
		beego.Error("SettingUpdateFailed|", err)
		this.Serve500(errors.New("Update Failed."))
	} else {
		this.Serve204()
	}
}

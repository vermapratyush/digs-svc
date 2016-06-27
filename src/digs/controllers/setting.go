package controllers

import (
	"github.com/astaxie/beego"
	"encoding/json"
	"digs/models"
	"errors"
	"digs/socket"
)

type SettingController struct {
	HttpBaseController
}

func (this *SettingController) Post() {

	var dataObject map[string]interface{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &dataObject)
	beego.Info("UpdateSetting|postData=", dataObject)

	sid, _ := dataObject["sessionId"].(string)
	userAuth, err := models.FindSession("sid", sid)
	if err != nil {
		beego.Error("SessionNotFound|err=", err)
		this.Serve500(errors.New("Unable to find session"))
		return
	}
	err = models.UpdateUserAccount(userAuth.UID, &dataObject)
	enablePush, _ := dataObject["enableNotification"].(bool)
	socket.LookUp[userAuth.UID].PushNotificationEnabled = enablePush

	if err != nil {
		beego.Error("SettingUpdateFailed|", err)
		this.Serve500(errors.New("Update Failed."))
	} else {
		this.Serve204()
	}
}

func (this *SettingController) Get()  {
	sid := this.GetString("sid")
	userAuth, err := models.FindSession("sid", sid)
	beego.Info("GetSetting|sid=", sid)
	if err != nil {
		beego.Error(err)
		this.Serve500(errors.New("Unable to get session"))
		return
	}
	userAccount, _ := models.GetUserAccount("uid", userAuth.UID)
	if userAccount.Settings == nil {
		this.Serve304()
		return
	}
	beego.Info("Setting=", userAccount.Settings)
	this.Serve200(&userAccount.Settings)

}

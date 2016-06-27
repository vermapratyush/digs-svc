package controllers

import (
	"github.com/astaxie/beego"
	"encoding/json"
	"digs/models"
	"errors"
	"digs/socket"
	"digs/domain"
)

type SettingController struct {
	HttpBaseController
}

func (this *SettingController) Post() {

	var request domain.SettingRequest
	json.Unmarshal(this.Ctx.Input.RequestBody, &request)
	request.Range = request.Range * 1000
	beego.Info("UpdateSetting|postData=", request)

	userAuth, err := models.FindSession("sid", request.SessionID)
	if err != nil {
		beego.Error("SessionNotFound|err=", err)
		this.Serve500(errors.New("Unable to find session"))
		return
	}
	err = models.UpdateUserAccount(userAuth.UID, &request)

	go updateLookUpTable(request.PushNotification, userAuth.UID)

	if err != nil {
		beego.Error("SettingUpdateFailed|", err)
		this.Serve500(errors.New("Update Failed."))
	} else {
		this.Serve204()
	}
}

func updateLookUpTable(notification bool, uid string) {

	var peer = socket.LookUp[uid]
	peer.PushNotificationEnabled = notification
	socket.LookUpLock.Lock()
	socket.LookUp[uid] = peer
	socket.LookUpLock.Unlock()
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
	beego.Info("Setting=", userAccount.Settings)
	this.Serve200(&userAccount.Settings)

}

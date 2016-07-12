package controllers

import (
	"github.com/astaxie/beego"
	"encoding/json"
	"digs/models"
	"errors"
	"digs/socket"
	"digs/domain"
	"gopkg.in/mgo.v2"
)

type SettingController struct {
	HttpBaseController
}

func (this *SettingController) Post() {

	var request domain.SettingRequest
	json.Unmarshal(this.Ctx.Input.RequestBody, &request)
	beego.Info("UpdateSetting|postData=", request)

	userAuth, err := models.FindSession("sid", request.SessionID)
	if err != nil {
		beego.Error("SessionNotFound|err=", err)
		if err == mgo.ErrNotFound {
			this.InvalidSessionResponse()
			return
		}
		this.Serve500(err)
		return
	} else {
		userAccount, err := models.GetUserAccount("uid", userAuth.UID)
		if (err == nil && userAccount.Settings.Range != request.Range) {
			go updatePersonActivity(userAccount, userAccount.Settings.Range, request.Range)
		} else if err == nil && userAccount.Settings.PublicProfile != request.PublicProfile {
			go hideJoinPersonActivity(request, userAuth.UID)
		}
	}

	err = models.UpdateUserAccount(userAuth.UID, &request)

	if err != nil {
		beego.Error("SettingUpdateFailed|", err)
		this.Serve500(errors.New("Update Failed."))
	} else {
		this.Serve204()
	}
}

func updatePersonActivity(userAccount *models.UserAccount, oldRange, newRange float64)  {
	userLocation, err := models.GetUserLocation(userAccount.UID)
	err1 := models.UpdateMessageRange(userAccount.UID, newRange)
	if (err == nil && len(userLocation.Location.Coordinates) != 0 && err1 == nil) {
		if oldRange > newRange {
			uidList := models.GetLiveUIDForFeed(userLocation.Location.Coordinates[0], userLocation.Location.Coordinates[1], oldRange, newRange)
			beego.Info("SettingsChanged|InformPartialUser=", uidList)
			socket.MulticastPersonCustom("leave", userAccount, userLocation.Location, uidList)
		} else {
			uidList := models.GetLiveUIDForFeed(userLocation.Location.Coordinates[0], userLocation.Location.Coordinates[1], newRange, oldRange)
			beego.Info("SettingsChanged|InformPartialUser=", uidList)
			socket.MulticastPersonCustom("join", userAccount, userLocation.Location, uidList)
		}
	}
}

func hideJoinPersonActivity(settingRequest domain.SettingRequest, uid string) {

	if settingRequest.PublicProfile == false {
		socket.MulticastPerson(uid, "hide")
	} else {
		socket.MulticastPerson(uid, "join")
	}

}

func (this *SettingController) Get()  {
	sid := this.GetString("sessionId")
	userAuth, err := models.FindSession("sid", sid)
	beego.Info("GetSetting|sid=", sid)
	if err != nil {
		beego.Error(err)
		if err == mgo.ErrNotFound {
			this.InvalidSessionResponse()
			return
		}
		this.Serve500(err)
		return
	}
	userAccount, _ := models.GetUserAccount("uid", userAuth.UID)
	beego.Info("Setting=", userAccount.Settings)
	this.Serve200(&userAccount.Settings)

}

package controllers

import (
	"encoding/json"
	"digs/models"
	"errors"
	"digs/socket"
	"digs/domain"
	"gopkg.in/mgo.v2"
	"digs/logger"
)

type SettingController struct {
	HttpBaseController
}

func (this *SettingController) Post() {

	var request domain.SettingRequest
	json.Unmarshal(this.Ctx.Input.RequestBody, &request)
	logger.Debug("SETTING|UpdateSetting|postData=%v", request)

	userAuth, err := models.FindSession("sid", request.SessionID)
	if err != nil {
		if err == mgo.ErrNotFound {
			this.InvalidSessionResponse()
			return
		}
		logger.Error("SESSION|UnableToFindSession|SID=", request.SessionID, "|Err=%v", err)
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
		logger.Error("SESSION|SettingUpdateFailed|SID=", request.SessionID, "|SettingRequest=%v", request, "|Err=%v", err)
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
			logger.Debug("SettingsChanged|UID=", userAccount.UID, "|InformPartialUser=%v", uidList)
			socket.MulticastPersonCustom("leave", userAccount, userLocation.Location, uidList, "")
		} else {
			uidList := models.GetLiveUIDForFeed(userLocation.Location.Coordinates[0], userLocation.Location.Coordinates[1], newRange, oldRange)
			logger.Debug("SettingsChanged|UID=", userAccount.UID, "|InformPartialUser=%v", uidList)
			socket.MulticastPersonCustom("join", userAccount, userLocation.Location, uidList, "")
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
	logger.Debug("GetSetting|SID=", sid, "|UID=", userAuth.UID)
	if err != nil {

		if err == mgo.ErrNotFound {
			this.InvalidSessionResponse()
			return
		}
		logger.Error("SESSION|UnableToFindSession|SID=", sid, "|Err=%v", err)
		this.Serve500(err)
		return
	}
	userAccount, _ := models.GetUserAccount("uid", userAuth.UID)
	logger.Debug("RESPONSE|GetSetting|SID=", sid, "|UID=", userAuth.UID, "|Response=", userAccount.Settings)
	this.Serve200(&userAccount.Settings)

}

package controllers

import (
	"digs/domain"
	"encoding/json"
	"digs/models"
	"gopkg.in/mgo.v2"
	"digs/logger"
)

type LogoutController struct {
	HttpBaseController
}

func (this *LogoutController) Post()  {
	var request domain.UserLogoutRequest
	logger.Debug("LOGOUT|Request=", string(this.Ctx.Input.RequestBody))
	this.Super(&request.BaseRequest)
	json.Unmarshal(this.Ctx.Input.RequestBody, &request)

	userAuth, err := models.FindSession("sid", request.SessionID)
	logger.Debug(request.SessionID)
	if err != nil {
		if err == mgo.ErrNotFound {
			this.InvalidSessionResponse()
			return
		}
		logger.Error("LOGOUT|LogoutFailed=", err)
		this.Serve500(err)
		return
	}

	err = models.DeleteNotificationId(userAuth.UID, request.NotificationId)

	if err != nil {
		logger.Error("NotificationDeleteFailed|SID=", userAuth, "|UID=", userAuth.UID, "|NotificationId=", request.NotificationId, "|err=", err)
	}

	err = models.DeleteUserAuth(userAuth.Id)
	if err != nil {
		logger.Error("UserAuthDeleteFailed|SID=", userAuth, "|UID=", userAuth.UID, "|NotificationId=", request.NotificationId, "|err=", err)
		this.Serve500(err)
		return
	}
	this.Serve204()
}


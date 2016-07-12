package controllers

import (
	"digs/domain"
	"encoding/json"
	"digs/models"
	"errors"
	"gopkg.in/mgo.v2"
	"digs/logger"
)

type NotificationController struct {
	HttpBaseController
}

func (this *NotificationController) Post() {
	var request domain.NotificationRequest
	logger.Debug("NOTIFICATION|Request=", string(this.Ctx.Input.RequestBody))
	this.Super(&request.BaseRequest)
	json.Unmarshal(this.Ctx.Input.RequestBody, &request)

	session, err := models.FindSession("sid", request.SessionID)
	if err != nil {
		if err == mgo.ErrNotFound {
			this.InvalidSessionResponse()
			return
		}
		logger.Error("NOTIFICATION|SID=", session.SID, "|UID=", session.UID, "|Err=", err)
		this.Serve500(err)
		return
	}
	err = models.AddNotificationId(session.UID, request.NotificationID, request.OSType)
	if err != nil {
		logger.Error("NOTIFICATION|NotificationAddFialed|SID=", session.SID, "|UID=", session.UID, "|Err=", err)
		this.Serve500(errors.New("Unable to register device"))
	} else {
		if request.NotificationID != request.OldNotificationID && request.OldNotificationID != "" {
			models.DeleteNotificationId(session.UID, request.OldNotificationID)
		}
		this.Serve204()
	}

}

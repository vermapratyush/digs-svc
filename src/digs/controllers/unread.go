package controllers

import (
	"digs/domain"
	"encoding/json"
	"digs/models"
	"gopkg.in/mgo.v2"
	"digs/logger"
)

type UnreadController struct {
	HttpBaseController
}

func (this *UnreadController) Post() {
	var request domain.UnreadRequest
	logger.Debug("UNREAD|Request=", string(this.Ctx.Input.RequestBody))
	this.Super(&request.BaseRequest)
	json.Unmarshal(this.Ctx.Input.RequestBody, &request)

	session, err := models.FindSession("sid", request.SessionID)
	if err != nil {
		if err == mgo.ErrNotFound {
			this.InvalidSessionResponse()
			return
		}
		logger.Error("UNREAD|SID=", session.SID, "|UID=", session.UID, "|Err=%v", err)
		this.Serve500(err)
		return
	}

	err = models.UpdateUnreadPointer(request.GID, session.UID, request.MID)

	if err != nil {
		this.Serve500(err)
		return
	}
	this.Serve204()
}

package controllers

import (
	"github.com/astaxie/beego"
	"encoding/json"
	"digs/domain"
	"digs/models"
	"gopkg.in/mgo.v2"
)

type AbuseController struct {
	HttpBaseController
}

func (this *AbuseController) Post() {
	var request domain.AbuseRequest

	beego.Info("REQUEST|AbuseRequest|", string(this.Ctx.Input.RequestBody))
	this.Super(&request.BaseRequest)
	json.Unmarshal(this.Ctx.Input.RequestBody, &request)
	session, err := models.FindSession("sid", request.SessionID)
	if err != nil {
		if err == mgo.ErrNotFound {
			this.InvalidSessionResponse()
			return
		}
		this.Serve500(err)
		return
	}

	if request.UID != "" {
		models.AddToBlockedContent(session.UID, "blockedUsers", request.UID)
		models.RemoveUserFromFeed(session.UID, request.UID)
	}
	if request.MID != "" {
		models.AddToBlockedContent(session.UID, "blockedMessages", request.MID)
		models.RemoveMessage(session.UID, request.MID)
	}

	this.Serve204()
}
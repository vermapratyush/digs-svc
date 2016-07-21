package controllers

import (
	"digs/domain"
	"encoding/json"
	"digs/models"
	"gopkg.in/mgo.v2"
	"digs/logger"
	"digs/common"
)

type MessagePinController struct {
	HttpBaseController
}

func (this *MessagePinController) Put() {
	var request domain.MessagePinAddDeleteRequest
	logger.Debug("AddMessagePin|Request=", string(this.Ctx.Input.RequestBody))
	this.Super(&request.BaseRequest)
	json.Unmarshal(this.Ctx.Input.RequestBody, &request)

	session, err := models.FindSession("sid", request.SessionID)
	if err != nil {
		if err == mgo.ErrNotFound {
			this.InvalidSessionResponse()
			return
		}
		logger.Error("AddMessagePin|SID=", session.SID, "|UID=", session.UID, "|Err=%v", err)
		this.Serve500(err)
		return
	}

	err = models.AddPinnedMessage(session.UID, request.MID)

	if err != nil {
		this.Serve500(err)
		return
	}
	this.Serve204()
}

func (this *MessagePinController) Delete() {
	var request domain.MessagePinAddDeleteRequest
	logger.Debug("DeleteMessagePin|Request=", string(this.Ctx.Input.RequestBody))
	this.Super(&request.BaseRequest)
	json.Unmarshal(this.Ctx.Input.RequestBody, &request)

	session, err := models.FindSession("sid", request.SessionID)
	if err != nil {
		if err == mgo.ErrNotFound {
			this.InvalidSessionResponse()
			return
		}
		logger.Error("DeleteMessagePin|SID=", session.SID, "|UID=", session.UID, "|Err=%v", err)
		this.Serve500(err)
		return
	}

	err = models.RemovePinnedMessage(session.UID, request.MID)

	if err != nil {
		this.Serve500(err)
		return
	}
	this.Serve204()
}

func (this *MessagePinController) Get() {

	sessionId := this.GetString("sessionId")
	userId := this.GetString("userId")

	logger.Debug("GetMessagePin|Request|SID=", sessionId, "|UID=", userId)

	session, err := models.FindSession("sid", sessionId)
	if err != nil {
		if err == mgo.ErrNotFound {
			this.InvalidSessionResponse()
			return
		}
		logger.Error("UNREAD|SID=", session.SID, "|UID=", session.UID, "|Err=%v", err)
		this.Serve500(err)
		return
	}

	userAccount, err := models.GetUserAccount("uid", userId)

	if err != nil {
		this.Serve500(err)
		return
	}

	messages, err := models.GetResolvedMessages(userAccount.PinnedMessages)
	if err != nil {
		this.Serve500(err)
		return
	}
	response := make([]domain.MessageGetResponse, len(messages))
	for idx, message := range (messages) {
		response[idx] = domain.MessageGetResponse{
			UID: message.From,
			MID: message.MID,
			Verified:message.UserAccount.Verified,
			From: common.GetName(message.UserAccount.FirstName, message.UserAccount.LastName),
			About: message.UserAccount.About,
			Message: message.Content,
			Timestamp: message.CreationTime.Unix() * int64(1000),
			ProfilePicture: message.UserAccount.ProfilePicture,
		}
	}
	this.Serve200(response)
}

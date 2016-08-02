package controllers

import (
	"digs/models"
	"digs/domain"
	"digs/common"
	"gopkg.in/mgo.v2"
	"digs/logger"
	"encoding/json"
	"fmt"
	"digs/mapper"
)

type GroupController struct {
	HttpBaseController
}

func (this *GroupController) Get() {
	gid := this.GetString("groupId")
	from, _ := this.GetInt64("from", 0)

	sid := this.GetString("sessionId")
	userAuth, err := models.FindSession("sid", sid)
	if err != nil {
		if err == mgo.ErrNotFound {
			this.InvalidSessionResponse()
			return
		}
		this.Serve500(err)
		return
	}
	messages, _ := models.GetMessageFromGroup(gid, from, common.MessageBatchSize)
	logger.Debug("GROUPGetResponse|Sid=", sid, "|UID=", userAuth.UID, "|GID=", gid, "|ResponseSize=", len(messages))

	response := composeResponse(gid, messages)
	this.Serve200(response)
}

func (this *GroupController) Post() {
	var request domain.GroupCreateRequest
	logger.Debug("REQUEST|GroupCreateRequest|", string(this.Ctx.Input.RequestBody))
	this.Super(&request.BaseRequest)
	json.Unmarshal(this.Ctx.Input.RequestBody, &request)

	sid := request.SessionID
	userAuth, err := models.FindSession("sid", sid)
	if err != nil {
		if err == mgo.ErrNotFound {
			this.InvalidSessionResponse()
			return
		}
		this.Serve500(err)
		return
	}

	userGroup := models.UserGroup{}
	messages := []models.UserGroupMessageResolved{}
	response := domain.CreateGroupResponse{}

	if this.Ctx.Input.Param(":version") == "v1" {
		request.IsGroup = false
	}
	//One-to-one chat
	if !request.IsGroup {
		otherUID := request.UIDS[0]
		if otherUID == userAuth.UID {
			otherUID = request.UIDS[1]
		}
		otherUserAccount, _ := models.GetUserAccount("uid", otherUID)
		userAccount, _ := models.GetUserAccount("uid", userAuth.UID)
		intersect := models.GetOneToOneCommonId(userAccount, otherUserAccount)
		if intersect == "" {
			userGroup, err = CreateOneToOneGroupChat("One-To-One-Group", fmt.Sprintf("Betweem %s and %s", request.UIDS[0], request.UIDS[1]), request.UIDS)
			if err != nil {
				this.Serve500(err)
				return
			}
		} else {
			userGroup, _ = models.GetGroupAccount(intersect)
			messages, _ = models.GetMessageFromGroup(userGroup.GID, 0, common.MessageBatchSize)
		}
	} else {
		userGroup, err = CreateGroupChat(request.GroupName, request.GroupAbout, request.GroupPicture, request.UIDS)
		if err != nil {
			this.Serve500(err)
			return
		}
	}

	response.GID = userGroup.GID
	response.GroupName = userGroup.GroupName
	response.GroupAbout = userGroup.GroupAbout
	response.GroupPicture = userGroup.GroupPicture
	response.Messages = composeResponse(userGroup.GID, messages)

	logger.Debug("GROUPCreateResponse|Sid=", sid, "|UID=", userAuth.UID, "GID=", response.GID, "|ResponseSize=", len(response.Messages))
	this.Serve200(response)
}

func (this *GroupController) GetDetails() {
	sid := this.GetString("sessionId")
	_, err := models.FindSession("sid", sid)
	if err != nil {
		if err == mgo.ErrNotFound {
			this.InvalidSessionResponse()
			return
		}
		this.Serve500(err)
		return
	}
	gid := this.Ctx.Input.Param(":groupId")
	userGroup, _ := models.GetGroupAccount(gid)

	userAccounts, _ := models.GetAllUserAccountIn(userGroup.UIDS)
	this.Serve200(&domain.GroupDetail {
		GID: gid,
		Users: mapper.MapUserAccountToPersonResponse(userAccounts),
		GroupName: userGroup.GroupName,
		GroupAbout: userGroup.GroupAbout,
		GroupPicture: userGroup.GroupPicture,
	})
}

func (this *GroupController) JoinGroup() {
	sid := this.GetString("sessionId")
	userAccount, err := models.FindSession("sid", sid)
	if err != nil {
		if err == mgo.ErrNotFound {
			this.InvalidSessionResponse()
			return
		}
		this.Serve500(err)
		return
	}
	gid := this.Ctx.Input.Param(":groupId")
	err = AddUserToGroup(userAccount.UID, gid)
	if err != nil {
		this.Serve500(err)
		return
	}
	this.Serve204()
}

func (this *GroupController) LeaveGroup() {
	sid := this.GetString("sessionId")
	userAccount, err := models.FindSession("sid", sid)
	if err != nil {
		if err == mgo.ErrNotFound {
			this.InvalidSessionResponse()
			return
		}
		this.Serve500(err)
		return
	}
	gid := this.Ctx.Input.Param(":groupId")
	err = RemoveUserFromGroup(userAccount.UID, gid)
	if err != nil {
		this.Serve500(err)
		return
	}
	this.Serve204()
}

func composeResponse(gid string, messages []models.UserGroupMessageResolved) []domain.MessageGetResponse {
	responseMessage := make([]domain.MessageGetResponse, len(messages))
	for idx, message := range (messages) {
		responseMessage[idx] = domain.MessageGetResponse{
			UID:message.UID,
			MID: message.MID,
			GID: gid,
			Verified:message.UserAccount.Verified,
			From: common.GetName(message.UserAccount.FirstName, message.UserAccount.LastName),
			About: message.UserAccount.About,
			Message: message.Content,
			Timestamp: message.CreationTime.UnixNano() / int64(1000000),
			ProfilePicture: message.UserAccount.ProfilePicture,
		}
	}
	return responseMessage
}
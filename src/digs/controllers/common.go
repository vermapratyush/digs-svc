package controllers

import (
	"digs/domain"
	"github.com/astaxie/beego"
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"digs/logger"
	"digs/models"
)

type HttpBaseController struct {
	beego.Controller
}

type WSBaseController struct {
	beego.Controller
	ws *websocket.Conn
}

func (this *HttpBaseController) Super(request *domain.BaseRequest) *HttpBaseController {
	request.HeaderSessionID = this.Ctx.Input.Header("SID")
	request.HeaderUserAgent = this.Ctx.Input.UserAgent()
	return this
}

func (this *WSBaseController) Prepare() {

	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		logger.Critical("Not a websocket handshake|Err=%v", err)
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		logger.Critical("Cannot setup WebSocket connection:%v", err)
		return
	}
	this.ws = ws
}

func (this *HttpBaseController) Serve500(err error) {
	this.Data["json"] = domain.GenericResponse{
		StatusCode:500,
		Message:err.Error(),
	}
	this.Ctx.Output.SetStatus(500)
	this.ServeJSON()
	return

}
func (this *HttpBaseController) ServeUnsupportedMedia() {
	this.Data["json"] = domain.GenericResponse{
		StatusCode:415,
	}
	this.Ctx.Output.SetStatus(415)
	this.ServeJSON()
	return

}

func (this *HttpBaseController) InvalidSessionResponse() {
	this.Data["json"] = &domain.GenericResponse{
		StatusCode:401,
		MessageCode:5000,
		Message:"Invalid Session",
	}
	this.Ctx.Output.SetStatus(401)
	this.ServeJSON()
}


func (this *HttpBaseController) Serve200(obj interface{}) {
	this.Data["json"] = obj
	this.Ctx.Output.SetStatus(200)
	this.ServeJSON()
}

func (this *HttpBaseController) Serve304() {
	this.Ctx.Output.SetStatus(304)
	this.ServeJSON()
}

func (this *HttpBaseController) Serve204() {
	this.Ctx.Output.SetStatus(204)
	this.ServeJSON()
}

func (this *WSBaseController) Respond(obj interface{})  {
	data, err := json.Marshal(obj)
	if err != nil {
		logger.Critical("Unable to reply back to the message sender Err=%v", err, "|obj=%v", obj)
		this.ws.WriteMessage(websocket.TextMessage, []byte("Unable to respond"))
		return
	}
	err = this.ws.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		logger.Critical("Error writing to websocket|Obj=%v", obj)
	}
}

func CreateGroupChat(groupName, groupAbout, groupPicture string, members []string) (models.UserGroup, error) {
	userGroup, err := models.CreateGroup(groupName, groupAbout, groupPicture, members)

	for _, uid := range (members) {
		_ = models.AddUserToGroupChat(uid, userGroup.GID)
	}

	return userGroup, err
}

func CreateOneToOneGroupChat(groupName, groupAbout string, members []string) (models.UserGroup, error) {
	userGroup, err := models.CreateGroup(groupName, groupAbout, "", members)

	for _, uid := range (members) {
		_ = models.AddUserToOneToOneGroupChat(uid, userGroup.GID)
	}

	return userGroup, err
}

func AddUserToGroup(uid, gid string) error {
	err := models.AddUserToGroupChat(uid, gid)
	if err != nil {
		return err
	}
	err = models.AddUserIdToGroup(uid, gid)
	if err != nil {
		return err
	}
	return nil
}

func RemoveUserFromGroup(uid, gid string) error {
	err := models.RemoveUserFromGroupChat(uid, gid)
	if err != nil {
		return err
	}
	err = models.RemoveUserIdFromGroup(uid, gid)
	if err != nil {
		return err
	}
	return nil
}
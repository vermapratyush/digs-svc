package controllers

import (
	"digs/domain"
	"github.com/astaxie/beego"
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
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
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}
	this.ws = ws
}

func (this *HttpBaseController) Serve500(err error) {
	this.Data["json"] = domain.ErrorResponse{
		StatusCode:500,
		Message:err.Error(),
	}
	this.Ctx.Output.SetStatus(500)
	this.ServeJSON()
	return

}

func (this *HttpBaseController) Serve200(obj interface{}) {
	this.Data["json"] = obj
	this.Ctx.Output.SetStatus(200)
	this.ServeJSON()
}

func (this *HttpBaseController) Serve204() {
	this.Ctx.Output.SetStatus(204)
	this.ServeJSON()
}

func (this *WSBaseController) Respond(obj interface{})  {
	data, err := json.Marshal(obj)
	if err != nil {
		beego.Critical("Unable to repoly back to the message sender Err=%s", err)
		this.ws.WriteMessage(websocket.TextMessage, []byte("Unable to respond"))
		return
	}
	err = this.ws.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		beego.Critical("Error writing to websocket")
	}
}
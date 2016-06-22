package socket

import (
	"github.com/gorilla/websocket"
	"digs/domain"
	"digs/models"
	"encoding/json"
	"github.com/astaxie/beego"
)

type Peer struct {
	Conn *websocket.Conn
	UID string
}

const (
	//MessageToServer
	UpdateLocation = "1:"
	SendMessage    = "2:"
	GetMessage     = "3:"
	Exit           = "4:"
	//MessageToClient
	Message        = "5:"
)

var LookUp = make(map[string]Peer)

func AddNode(uid string, ws *websocket.Conn) {
	beego.Info("NodeAdded|UID=", uid)
	LookUp[uid] = Peer{
		Conn:ws,
		UID:uid,
	}
}

func LeaveNode(uid string) {
	beego.Info("NodeLeft|UID=", uid)
	LookUp[uid].Conn.Close()
	delete(LookUp, uid)
}

func MulticastMessage(userAccount *models.UserAccount, msg *domain.MessageSendRequest) {
	defer DeadSocketWrite()
	uids := models.GetLiveUIDForFeed(msg.Location.Longitude, msg.Location.Latitude, msg.Reach)
	beego.Info("TotalUsers|Size=", len(uids))
	for idx := 0; idx < len(uids); idx++ {
		ws, present := LookUp[uids[idx]]
		if present == false || uids[idx] == userAccount.UID {
			continue
		}
		beego.Info("SendingMessage|From=", userAccount.UID, "|To=", uids[idx])
		response, err := json.Marshal(domain.MessageGetResponse{
			From:userAccount.FirstName + userAccount.LastName,
			UID:userAccount.UID,
			Message: msg.Body,
			Timestamp: msg.Timestamp,

		})
		if err != nil {
			beego.Critical("MessageSendFailed|Error=", uids[idx])
			continue
		}
		ws.Conn.WriteMessage(websocket.TextMessage, response)
	}
}

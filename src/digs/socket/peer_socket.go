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
	NotificationId []string
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
	notifications, err := models.GetNotificationIds(uid)
	nid := make(map[string]struct{})
	if err != nil {
		beego.Info("No notifications for user ", uid)
	} else {
		for _, notification := range (*notifications) {
			nid[notification.NotificationId]  = struct{}{}
		}
	}
	nidArray := make([]string, len(nid))
	idx := 0
	for k, _ := range (nid) {
		nidArray[idx] = k
		idx = idx + 1
	}

	beego.Info("NodeAdded|UID=", uid)
	LookUp[uid] = Peer{
		Conn:ws,
		UID:uid,
		NotificationId:nidArray,
	}
}

func LeaveNode(uid string) {
	beego.Info("NodeLeft|UID=", uid)
	_, ok := LookUp[uid]
	if ok && LookUp[uid].Conn != nil {
		LookUp[uid].Conn.Close()
	}
	delete(LookUp, uid)
}

func MulticastMessage(userAccount *models.UserAccount, msg *domain.MessageSendRequest) {
	defer DeadSocketWrite(userAccount)
	uids := models.GetLiveUIDForFeed(msg.Location.Longitude, msg.Location.Latitude, msg.Reach)
	beego.Info("TotalUsers|Size=", len(uids))
	for idx := 0; idx < len(uids); idx++ {
		ws, present := LookUp[uids[idx]]
		if present == false || uids[idx] == userAccount.UID {
			continue
		}
		beego.Info("SendingMessage|From=", userAccount.UID, "|To=", uids[idx])
		response, _ := json.Marshal(domain.MessageGetResponse{
			From:userAccount.FirstName + userAccount.LastName,
			UID:userAccount.UID,
			Message: msg.Body,
			Timestamp: msg.Timestamp,
			ProfilePicture:userAccount.ProfilePicture,
		})
		err := ws.Conn.WriteMessage(websocket.TextMessage, response)
		if err != nil {
			beego.Critical("MessageSendFailed|Error=", uids[idx])
			LeaveNode(uids[idx])
			continue
		}
	}
}

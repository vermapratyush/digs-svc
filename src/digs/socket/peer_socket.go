package socket

import (
	"github.com/gorilla/websocket"
	"digs/domain"
	"digs/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/NaySoftware/go-fcm"
	"digs/common"
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
		beego.Info("from", userAccount.UID, "to=", uids[idx])
		if uids[idx] == userAccount.UID {
			continue
		} else if (present == false) {
			sendPushMessage(userAccount, uids[idx], msg)
		} else {
			sendWSMessage(userAccount, uids[idx], msg, ws.Conn)
		}
	}


}

func sendWSMessage(userAccount *models.UserAccount, toUID string, msg *domain.MessageSendRequest, ws *websocket.Conn) {
	beego.Info("SendingWSMessage|From=", userAccount.UID, "|To=", toUID)
	response, _ := json.Marshal(domain.MessageGetResponse{
		From:userAccount.FirstName + userAccount.LastName,
		UID:userAccount.UID,
		Message: msg.Body,
		Timestamp: msg.Timestamp,
		ProfilePicture:userAccount.ProfilePicture,
	})
	err := ws.WriteMessage(websocket.TextMessage, response)
	if err != nil {
		beego.Critical("MessageSendFailed|Error=", toUID)
		LeaveNode(toUID)
		sendPushMessage(userAccount, toUID, msg)
	}
}

func sendPushMessage(userAccount *models.UserAccount, toUID string, msg *domain.MessageSendRequest) {
	beego.Info("SendingPushMessage|From=", userAccount.UID, "|To=", toUID)
	fcm := fcm.NewFcmClient(common.PushNotification_API_KEY)

	data := map[string]string{
		"title": userAccount.FirstName,
		"message": msg.Body,
		"image": "twitter",
	}

	notifications, err := models.GetNotificationIds(toUID)

	nid := make(map[string]struct{})
	if err != nil {
		beego.Info("No push notifications for user ", toUID)
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
	fcm.NewFcmRegIdsMsg(nidArray, data)
	status, err := fcm.Send(1)
	if err == nil {
		beego.Info(status)
	} else {
		beego.Error(err)
	}

}
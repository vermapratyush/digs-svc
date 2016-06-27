package socket

import (
	"github.com/gorilla/websocket"
	"digs/domain"
	"digs/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/NaySoftware/go-fcm"
	"digs/common"
	"sync"
	"runtime/debug"
)

type Peer struct {
	Conn *websocket.Conn
	UID string
	PushNotificationEnabled bool
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
var LookUpLock sync.RWMutex

func AddNode(uid string, ws *websocket.Conn) {

	beego.Info("NodeAdded|UID=", uid)

	LookUpLock.Lock()
	LookUp[uid] = Peer{
		Conn:ws,
		UID:uid,
	}
	LookUpLock.Unlock()
	MulticastPerson(uid, "join")
}

func LeaveNode(uid string) {
	beego.Info("NodeLeft|UID=", uid)
	_, present := LookUp[uid]
	if present && LookUp[uid].Conn != nil {
		MulticastPerson(uid, "leave")

		defer DeadSocketWrite(LookUp[uid])
		LookUp[uid].Conn.Close()
	} else if present {
		beego.Info("WSAlredyClosed|uid=", uid)
	}

	LookUpLock.Lock()
	delete(LookUp, uid)
	LookUpLock.Unlock()
}

func MulticastPerson(uid string, activity string) {
	userLocation, err := models.GetUserLocation(uid)
	if err != nil {
		beego.Error("UserLocationNotFound|err=", err)
		return
	}
	userAccount, _ := models.GetUserAccount("uid", uid)
	messageRange, _ := userAccount.Settings["messageRange"].(int64)
	isPublic, _ := userAccount.Settings["publicProfile"].(bool)

	//HACK: Find a better solution
	enableNotification, _ := userAccount.Settings["enableNotification"].(bool)
	var peer = LookUp[uid]
	peer.PushNotificationEnabled = enableNotification
	LookUpLock.Lock()
	LookUp[uid] = peer
	LookUpLock.Unlock()


	uids := models.GetLiveUIDForFeed(userLocation.Location.Coordinates[0], userLocation.Location.Coordinates[1], messageRange)
	for _, toUID := range(uids) {
		peer, present := LookUp[toUID]
		if uid == toUID || present == false {
			continue
		} else if (isPublic) {
			beego.Info("Stacktrace", string(debug.Stack()))
			beego.Info("Person=", uid, " activity=", activity, " to=", toUID)
			response, _ := json.Marshal(&domain.PersonResponse{
				Name: common.GetName(userAccount.FirstName, userAccount.LastName),
				UID: uid,
				Activity: activity,
				About: userAccount.About,
				ProfilePicture: userAccount.ProfilePicture,
			})
			err := sendWSMessage(peer, uid, response)
			if err != nil {
				beego.Error("MessageSendFailed|err=", err)
			}
		}
	}
}

func MulticastMessage(userAccount *models.UserAccount, msg *domain.MessageSendRequest) {
	reach, _ := userAccount.Settings["messageReach"].(int64)
	beego.Info("Searching people in radius of ", reach, " from ", msg.Location)
	uids := models.GetLiveUIDForFeed(msg.Location.Longitude, msg.Location.Latitude, reach)
	beego.Info("TotalUsers|Size=", len(uids))
	for idx := 0; idx < len(uids); idx++ {
		peer, present := LookUp[uids[idx]]
		beego.Info("from", userAccount.UID, "to=", uids[idx])
		if uids[idx] == userAccount.UID {
			continue
		} else if (present == false && peer.PushNotificationEnabled) {
			sendPushMessage(userAccount, uids[idx], msg)

		} else {
			response, _ := json.Marshal(domain.MessageGetResponse{
				From:common.GetName(userAccount.FirstName, userAccount.LastName),
				UID:userAccount.UID,
				Message: msg.Body,
				Timestamp: msg.Timestamp,
				ProfilePicture:userAccount.ProfilePicture,
			})
			err := sendWSMessage(peer, userAccount.UID, response)
			if err != nil && peer.PushNotificationEnabled {
				sendPushMessage(userAccount, uids[idx], msg)
			}
		}
	}


}

func sendWSMessage(toPeer Peer, fromUID string, data []byte) error {
	defer DeadSocketWrite(toPeer)
	beego.Info("SendingWSMessage|From=", fromUID, "|To=", toPeer.UID)

	err := toPeer.Conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		beego.Critical("MessageSendFailed|Error=", toPeer.UID)
		LeaveNode(fromUID)
		return err
	}
	return nil
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
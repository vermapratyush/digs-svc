package socket

import (
	"github.com/gorilla/websocket"
	"digs/domain"
	"digs/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"digs/common"
	"sync"
)


const (
	//MessageToServer
	UpdateLocation = "1:"
	SendMessage    = "2:"
	GetMessage     = "3:"
	Exit           = "4:"
	//MessageToClient
	Message        = "5:"
)

func AddNode(uid string, ws *websocket.Conn) {
	if uid == "" {
		return
	}
	beego.Info("NodeAdded|UID=", uid)

	SetLookUp(uid, Peer{
		Conn:ws,
		UID:uid,
		wsLock:&sync.Mutex{},
	})
	MulticastPerson(uid, "join")
}

func LeaveNode(uid string) {
	if uid == "" {
		return
	}
	beego.Info("NodeLeft|UID=", uid)
	ws, present := GetLookUp(uid)
	RemoveLookUp(uid)

	if present && ws.Conn != nil {
		MulticastPerson(uid, "leave")

		defer DeadSocketWrite(uid)
		ws.Conn.Close()
	} else if present {
		beego.Info("WSAlredyClosed|uid=", uid)
	}


}

func MulticastPerson(uid string, activity string) {
	userLocation, err := models.GetUserLocation(uid)
	if err != nil {
		beego.Error("UserLocationNotFound|err=", err)
		return
	}
	userAccount, _ := models.GetUserAccount("uid", uid)
	uidList := models.GetLiveUIDForFeed(userLocation.Location.Coordinates[0], userLocation.Location.Coordinates[1], userAccount.Settings.Range, -1)
	beego.Info(uidList)
	MulticastPersonCustom(activity, userAccount, userLocation.Location, uidList)

}

func MulticastPersonCustom(activity string, userAccount *models.UserAccount, userLocation models.Coordinate, uids []string)  {

	for _, toUID := range(uids) {
		peer, present := GetLookUp(toUID)
		beego.Info("From=", userAccount.UID, "|To=", toUID, "|present=", present)
		if toUID == "" || userAccount.UID == toUID || present == false {
			continue
		} else if (activity == "hide" || activity == "show" || userAccount.Settings.PublicProfile) {
			beego.Info("Person=", userAccount.UID, " activity=", activity, " to=", toUID)
			response, _ := json.Marshal(&domain.PersonResponse{
				Name: common.GetName(userAccount.FirstName, userAccount.LastName),
				UID: userAccount.UID,
				Activity: activity,
				About: userAccount.About,
				ProfilePicture: userAccount.ProfilePicture,
			})
			err := sendWSMessage(peer, userAccount.UID, response)
			if err != nil {
				beego.Error("MessageSendFailed|err=", err)
			}
		}
	}
}

func MulticastMessage(userAccount *models.UserAccount, msg *domain.MessageSendRequest) {

	beego.Info("Searching people in radius of ", userAccount.Settings.Range, " from ", msg.Location)
	uids := models.GetLiveUIDForFeed(msg.Location.Longitude, msg.Location.Latitude, userAccount.Settings.Range, -1)
	beego.Info("TotalUsers|Size=", len(uids))
	for _, toUID := range(uids) {
		peer, present := GetLookUp(toUID)
		if toUID == "" || toUID == userAccount.UID {
			continue
		}
		toUserAccount, _ := models.GetUserAccount("uid", toUID)
		models.AddToUserFeed(toUID, msg.MID)
		userLocation, err := models.GetUserLocation(toUID)
		if err != nil || toUserAccount.Settings.Range >= common.Distance(msg.Location.Latitude, msg.Location.Longitude, userLocation.Location.Coordinates[1], userLocation.Location.Coordinates[0]) {
			continue
		}

		if (!present && toUserAccount.Settings.PushNotification) {
			beego.Info("Push|from", userAccount.UID, "to=", toUID)
			sendPushMessage(userAccount, toUID, msg)

		} else if (present) {
			beego.Info("WS|from", userAccount.UID, "to=", toUID)
			beego.Info("Peer=", peer)
			response, _ := json.Marshal(domain.MessageGetResponse{
				From:common.GetName(userAccount.FirstName, userAccount.LastName),
				UID:userAccount.UID,
				MID:msg.MID,
				About: userAccount.About,
				Message: msg.Body,
				Timestamp: msg.Timestamp,
				ProfilePicture:userAccount.ProfilePicture,
			})
			err := sendWSMessage(peer, userAccount.UID, response)
			if err != nil && toUserAccount.Settings.PushNotification {
				sendPushMessage(userAccount, toUID, msg)
			}
		}
	}
}

func sendWSMessage(toPeer Peer, fromUID string, data []byte) error {
	defer DeadSocketWrite(toPeer.UID)
	beego.Info("SendingWSMessage|From=", fromUID, "|To=", toPeer.UID)

	err := SendData(toPeer.UID, data)
	if err != nil {
		beego.Critical("MessageSendFailed|Error=", toPeer.UID)
		LeaveNode(fromUID)
		return err
	}
	return nil
}

func sendPushMessage(userAccount *models.UserAccount, toUID string, msg *domain.MessageSendRequest) {
	beego.Info("SendingPushMessage|From=", userAccount.UID, "|To=", toUID)

	notifications, err := models.GetNotificationIds(toUID)
	if err != nil {
		beego.Error("NotificationIdFetch|err=", err)
		return
	}

	for _, notification := range(*notifications) {
		beego.Info(notification.OSType)
		if notification.OSType == "android" {
			androidPush(userAccount, notification.NotificationId, msg)
		} else {
			iosPush(userAccount, notification.NotificationId, msg)
		}
	}
}

func androidPush(userAccount *models.UserAccount, nid string, msg *domain.MessageSendRequest) {
	models.AndroidMessagePush(userAccount.UID, nid, msg.Body)

}

func iosPush(userAccount *models.UserAccount, nid string, msg *domain.MessageSendRequest)  {
	models.IOSMessagePush(userAccount.UID, nid, msg.Body)
}
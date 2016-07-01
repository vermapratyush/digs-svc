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
	"github.com/sideshow/apns2/certificate"
	apns "github.com/sideshow/apns2"
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
var LookUpLock sync.RWMutex

func AddNode(uid string, ws *websocket.Conn) {
	if uid == "" {
		return
	}
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
	uidList := models.GetLiveUIDForFeed(userLocation.Location.Coordinates[0], userLocation.Location.Coordinates[1], userAccount.Settings.Range, -1)
	MulticastPersonCustom(activity, userAccount, userLocation.Location, uidList)

}

func MulticastPersonCustom(activity string, userAccount *models.UserAccount, userLocation models.Coordinate, uids []string)  {

	for _, toUID := range(uids) {
		peer, present := LookUp[toUID]
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
		peer, present := LookUp[toUID]
		if toUID == "" || toUID == userAccount.UID {
			continue
		}
		pushEnabled, _ := models.GetUserAccount("uid", toUID)
		models.AddToUserFeed(toUID, msg.MID)
		if (!present && pushEnabled.Settings.PushNotification) {
			beego.Info("Push|from", userAccount.UID, "to=", toUID)
			sendPushMessage(userAccount, toUID, msg)

		} else if (present) {
			beego.Info("WS|from", userAccount.UID, "to=", toUID)
			beego.Info("Peer=", peer)
			response, _ := json.Marshal(domain.MessageGetResponse{
				From:common.GetName(userAccount.FirstName, userAccount.LastName),
				UID:userAccount.UID,
				MID:msg.MID,
				Message: msg.Body,
				Timestamp: msg.Timestamp,
				ProfilePicture:userAccount.ProfilePicture,
			})
			err := sendWSMessage(peer, userAccount.UID, response)
			if err != nil && pushEnabled.Settings.PushNotification {
				sendPushMessage(userAccount, toUID, msg)
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

	notifications, err := models.GetNotificationIds(toUID)
	if err != nil {
		beego.Error("NotificationIdFetch|err=", err)
		return
	}

	for _, notification := range(*notifications) {
		if notification.OSType == "android" {
			androidPush(userAccount, notification.NotificationId, msg)
		} else {
			iosPush(userAccount, notification.NotificationId, msg)
		}
	}
}

func androidPush(userAccount *models.UserAccount, nid string, msg *domain.MessageSendRequest) {
	beego.Info("AndroidPush|From=", userAccount.UID, "|To=", nid)
	fcm := fcm.NewFcmClient(common.PushNotification_API_KEY)

	data := map[string]string{
		"title": "powow",
		"message": msg.Body,
		"image": "twitter",
		"style": "inbox",
		"summaryText": "There are %n% notification",
	}

	fcm.NewFcmMsgTo(nid, data)
	status, err := fcm.Send(1)
	if err == nil {
		beego.Info(status)
	} else {
		beego.Error(err)
	}

}

func iosPush(userAccount *models.UserAccount, nid string, msg *domain.MessageSendRequest)  {
	beego.Info("IOSPush|From=", userAccount.UID, "|To=", nid)
	cert, pemErr := certificate.FromPemFile("socket/final.pem", "")
	if pemErr != nil {
		beego.Error("APNSCertError|err", pemErr)
		return
	}

	notification := &apns.Notification{}
	notification.DeviceToken = nid
	notification.Payload = []byte("{\"aps\":{\"alert\":\"" + msg.Body + "\"}}") // See Payload section below

	client := apns.NewClient(cert).Development()
	_, err := client.Push(notification)

	if err != nil {
		beego.Error("APNSPushError|err=", err)
		return
	}
}
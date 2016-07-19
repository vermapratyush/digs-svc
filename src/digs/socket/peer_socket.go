package socket

import (
	"github.com/gorilla/websocket"
	"digs/domain"
	"digs/models"
	"encoding/json"
	"digs/common"
	"sync"
	"digs/logger"
	"fmt"
)


const (
	//MessageToServer
	UpdateLocation  = "1:"
	SendMessage     = "2:"
	GroupMessage    = "3:"
	TypingMessage   = "4:"
	Exit            = "5:"
	//MessageToClient
	Message         = "6:"
)

func AddNode(uid string, ws *websocket.Conn) {
	if uid == "" {
		return
	}
	logger.Debug("PEER|NodeAdded|UID=", uid)

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
	logger.Debug("PEER|NodeLeft|UID=", uid)
	ws, present := GetLookUp(uid)
	RemoveLookUp(uid)

	if present && ws.Conn != nil {
		MulticastPerson(uid, "leave")

		defer DeadSocketWrite(uid)
		ws.Conn.Close()
	} else if present {
		logger.Debug("PEER|WSAlredyClosed|uid=", uid)
	}


}

func MulticastPerson(uid string, activity string) {
	userLocation, err := models.GetUserLocation(uid)
	if err != nil || len(userLocation.Location.Coordinates) == 0 {
		logger.Error("PEER|UserLocationNotFound|uid=",uid, "|err=%v", err)
		return
	}
	userAccount, _ := models.GetUserAccount("uid", uid)
	uidList := models.GetLiveUIDForFeed(userLocation.Location.Coordinates[0], userLocation.Location.Coordinates[1], userAccount.Settings.Range, -1)
	MulticastPersonCustom(activity, userAccount, userLocation.Location, uidList)

}

func MulticastPersonCustom(activity string, userAccount *models.UserAccount, userLocation models.Coordinate, uids []string)  {
	blockedMap := common.GetStringArrayAsMap(userAccount.BlockedUsers)

	for _, toUID := range(uids) {
		peer, present := GetLookUp(toUID)
		_, presentInBlock := blockedMap[toUID]

		if toUID == "" || userAccount.UID == toUID || present == false || presentInBlock {
			logger.Debug("PEER|NotSendingActivity|toUID=", toUID, "|FromUID=", userAccount.UID, "|Present=", present, "|Blocked=", presentInBlock)
			continue
		} else if (activity == "hide" || activity == "show" || userAccount.Settings.PublicProfile) {
			logger.Debug("PEER|NotSendingActivity|toUID=", toUID, "|FromUID=", userAccount.UID, "|Activity=", activity)
			response, _ := json.Marshal(&domain.PersonResponse{
				Name: common.GetName(userAccount.FirstName, userAccount.LastName),
				UID: userAccount.UID,
				Activity: activity,
				About: userAccount.About,
				ProfilePicture: userAccount.ProfilePicture,
			})
			err := sendWSMessage(peer, userAccount.UID, response)
			if err != nil {
				logger.Error("PEER|MessageSendFailed|ToUID=", toUID, "|From=", userAccount.UID, "|err=%v", err)
			}
		}
	}
}


func MulticastMessage(userAccount *models.UserAccount, msg *domain.MessageSendRequest) {

	uids := []string{}

	if msg.GID != "" {
		groupAccount, _ := models.GetGroupAccount(msg.GID)
		uids = groupAccount.UIDS
	} else {
		uids = models.GetLiveUIDForFeed(msg.Location.Longitude, msg.Location.Latitude, userAccount.Settings.Range, -1)
	}
	logger.Debug("TotalUsers|UID=", userAccount.UID, "|MID=", msg.MID, "|Location=%v", msg.Location, "|Size=", len(uids))
	sendingWS := make([]string, 0)
	sendingPush := make([]string, 0)

	for _, toUID := range(uids) {
		peer, present := GetLookUp(toUID)
		if toUID == "" || toUID == userAccount.UID {
			continue
		}
		toUserAccount, _ := models.GetUserAccount("uid", toUID)

		//Add to feed of the user or group accordingly
		AddToFeed(toUID, msg.GID, msg)

		//Blocking only works for group messages
		blocked := common.IsUserBlocked(toUserAccount.BlockedUsers, userAccount.UID)
		if blocked {
			logger.Debug("PEER|NotSendingMessage|toUID=", toUID, "|FromUID=", userAccount.UID, "|Blocked=", blocked)
			continue
		}

		if (!present && toUserAccount.Settings.PushNotification) {
			sendingPush = append(sendingPush, toUID)
			sendPushMessage(userAccount, toUID, msg)

		} else if (present) {
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
				sendingPush = append(sendingPush, toUID)
				sendPushMessage(userAccount, toUID, msg)
			} else if err == nil {
				sendingWS = append(sendingWS, toUID)
			}
		}
	}
	logger.Debug("WSMessage|len=", len(sendingWS), "|uid=", sendingWS)
	logger.Debug("PushMessage|len=", len(sendingPush), "|uid=", sendingPush)
}

func sendWSMessage(toPeer Peer, fromUID string, data []byte) error {
	defer DeadSocketWrite(toPeer.UID)

	err := SendData(toPeer.UID, data)
	if err != nil {
		logger.Critical("MessageSendFailed|ToUID=", toPeer.UID, "|From=", fromUID, "|Error=%v", err)
		LeaveNode(fromUID)
		return err
	}
	return nil
}

func sendPushMessage(userAccount *models.UserAccount, toUID string, msg *domain.MessageSendRequest) {

	notifications, err := models.GetNotificationIds(toUID)
	if err != nil {
		logger.Error("PEER|NotificationIdFetch|toUID=", toUID, "|err=%v", err)
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
	models.AndroidMessagePush(userAccount.UID, nid, fmt.Sprintf("%s: %s", userAccount.FirstName, msg.Body))

}

func iosPush(userAccount *models.UserAccount, nid string, msg *domain.MessageSendRequest)  {
	models.IOSMessagePush(userAccount.UID, nid, fmt.Sprintf("%s: %s", userAccount.FirstName, msg.Body))
}
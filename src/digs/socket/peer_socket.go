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
	"strings"
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
	MulticastPersonCustom(activity, userAccount, userLocation.Location, uidList, "")

}

func MulticastPersonCustom(activity string, userAccount *models.UserAccount, userLocation models.Coordinate, uids []string, gid string)  {
	blockedMap := common.GetStringArrayAsMap(userAccount.BlockedUsers)

	for _, toUID := range(uids) {
		peer, present := GetLookUp(toUID)
		_, presentInBlock := blockedMap[toUID]

		if toUID == "" || userAccount.UID == toUID || present == false || presentInBlock {
			logger.Debug("PEER|NotSendingActivity|toUID=", toUID, "|FromUID=", userAccount.UID, "|Present=", present, "|Blocked=", presentInBlock)
			continue
		} else if (activity == "hide" || activity == "show" || userAccount.Settings.PublicProfile) {
			logger.Debug("PEER|SendingActivity|toUID=", toUID, "|FromUID=", userAccount.UID, "|Activity=", activity)
			activeState := "active"
			if activity == "hide" || activity == "leave" {
				activeState = "inactive"
			}
			response, _ := json.Marshal(&domain.PersonResponse{
				Name: common.GetName(userAccount.FirstName, userAccount.LastName),
				UID: userAccount.UID,
				GID: gid,
				Activity: activity,
				About: userAccount.About,
				ActiveState:activeState,
				Verified:userAccount.Verified,
				ProfilePicture: userAccount.ProfilePicture,
				IsGroup: false,
			})
			err := sendWSMessage(peer, response)
			if err != nil {
				logger.Error("PEER|MessageSendFailed|ToUID=", toUID, "|From=", userAccount.UID, "|err=%v", err)
			}
		}
	}
}

func MulticastGroup(event *domain.PersonResponse, uids[] string) {
	for _, toUID := range(uids) {
		peer, present := GetLookUp(toUID)
		if present {
			logger.Debug("PEER|SendingActivity|toUID=", toUID, "|Activity=CreateGroup")
			messageByte, _:= json.Marshal(event)
			sendWSMessage(peer, messageByte)
		} else {
			sendPushPeople(toUID, event)
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
	logger.Debug("TotalUsers|UID=", userAccount.UID, "|MID=", msg.MID, "|Location=%v", msg.Location, "|Size=", uids)
	sendingWS := make([]string, 0)
	sendingPush := make([]string, 0)

	for _, toUID := range(uids) {
		peer, present := GetLookUp(toUID)
		if toUID == "" || toUID == userAccount.UID {
			continue
		}
		toUserAccount, _ := models.GetUserAccount("uid", toUID)

		//Add to feed of the user, group_feed has already been updated
		if msg.GID == "" {
			models.AddToUserFeed(toUID, msg.MID)
		}

		//Blocking only works for group messages
		blocked := common.IsUserBlocked(toUserAccount.BlockedUsers, userAccount.UID)
		if blocked {
			logger.Debug("PEER|NotSendingMessage|toUID=", toUID, "|FromUID=", userAccount.UID, "|Blocked=", blocked)
			continue
		}
		response := &domain.MessageGetResponse{
			From:common.GetName(userAccount.FirstName, userAccount.LastName),
			UID:userAccount.UID,
			GID:msg.GID,
			MID:msg.MID,
			Verified:userAccount.Verified,
			About: userAccount.About,
			Message: msg.Body,
			Timestamp: msg.Timestamp,
			ProfilePicture:userAccount.ProfilePicture,
		}
		responseString, _ := json.Marshal(response)

		if (!present && toUserAccount.Settings.PushNotification) {
			sendingPush = append(sendingPush, toUID)
			sendPushMessage(userAccount, toUID, response)

		} else if (present) {
			err := sendWSMessage(peer, responseString)
			if err != nil && toUserAccount.Settings.PushNotification {
				sendingPush = append(sendingPush, toUID)
				sendPushMessage(userAccount, toUID, response)
			} else if err == nil {
				sendingWS = append(sendingWS, toUID)
			}
		}
	}
	logger.Debug("WSMessage|len=", len(sendingWS), "|uid=", sendingWS)
	logger.Debug("PushMessage|len=", len(sendingPush), "|uid=", sendingPush)
}

func sendWSMessage(toPeer Peer, data []byte) error {
	defer DeadSocketWrite(toPeer.UID)

	err := SendData(toPeer.UID, data)
	if err != nil {
		logger.Critical("MessageSendFailed|ToUID=", toPeer.UID, "|Error=%v", err)
		return err
	}
	return nil
}

func sendPushMessage(userAccount *models.UserAccount, toUID string, response *domain.MessageGetResponse) {

	notifications, err := models.GetNotificationIds(toUID)
	if err != nil {
		logger.Error("PEER|NotificationIdFetch|toUID=", toUID, "|err=%v", err)
		return
	}

	for _, notification := range(*notifications) {
		if notification.OSType == "android" {
			androidMessagePush(userAccount, notification.NotificationId, response)
		} else {
			iosPush(userAccount, notification.NotificationId, response)
		}
	}
}

func sendPushPeople(uid string, response *domain.PersonResponse) {
	notifications, err := models.GetNotificationIds(uid)
	if err != nil {
		logger.Error("PEER|NotificationIdFetch|toUID=", uid, "|err=%v", err)
		return
	}
	message := fmt.Sprintf("You have been added to the group: %s", response.Name)
	for _, notification := range(*notifications) {
		if notification.OSType == "android" {
			androidPeoplePush(uid, notification.NotificationId, message)
		} else {
			iosPeoplePush(uid, notification.NotificationId, message)
		}
	}
}

func androidMessagePush(userAccount *models.UserAccount, nid string, msg *domain.MessageGetResponse) {
	additionalData, _ := json.Marshal(msg)
	if strings.Contains(msg.Message, "<img") {
		models.AndroidMessagePush(userAccount.UID, nid, fmt.Sprintf("%s has sent an image", userAccount.FirstName), string(additionalData))
	} else {
		models.AndroidMessagePush(userAccount.UID, nid, fmt.Sprintf("%s: %s", userAccount.FirstName, msg.Message), string(additionalData))
	}
}

func androidPeoplePush(uid, nid, message string) {
	models.AndroidMessagePush(uid, nid, message, "")
}

func iosPeoplePush(uid, nid, message string) {
	models.IOSMessagePush(uid, nid, message)
}

func iosPush(userAccount *models.UserAccount, nid string, msg *domain.MessageGetResponse)  {
	if strings.Contains(msg.Message, "<img") {
		models.IOSMessagePush(userAccount.UID, nid, fmt.Sprintf("%s has sent an image", userAccount.FirstName))
	} else {
		models.IOSMessagePush(userAccount.UID, nid, fmt.Sprintf("%s: %s", userAccount.FirstName, msg.Message))
	}

}
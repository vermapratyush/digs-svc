package controllers

import (
	"digs/domain"
	"strings"
	"digs/models"
	"encoding/json"
	"digs/socket"
	"digs/common"
	"digs/logger"
	"strconv"
)

type WSMessengerController struct {
	WSBaseController
}

func (this *WSMessengerController) Get() {
	sid := this.GetString("sessionId")

	logger.Debug("WSConnection|SID=", sid)

	userAuth, err := models.FindSession("sid", sid)
	if err != nil {
		this.Respond(&domain.GenericResponse{
			StatusCode:401,
			MessageCode:5000,
			Message:"Invalid Session",
		})
		return
	}
	logger.Debug("UserConnected|UID=%v", userAuth)
	this.Respond(&domain.GenericResponse{
		StatusCode: 200,
		Message: "Connection Established",
		MessageCode: 3000,
	})

	var location = domain.Coordinate{}
	if (this.GetString("longitude") != "" && this.GetString("latitude") != "") {
		latFloat, _ := strconv.ParseFloat(this.GetString("latitude"), 64)
		longFloat, _ := strconv.ParseFloat(this.GetString("longitude"), 64)
		models.AddUserNewLocation(longFloat, latFloat, userAuth.UID)
		location = domain.Coordinate{
			Longitude:longFloat,
			Latitude:latFloat,
		}
	}

	socket.AddNode(userAuth.UID, this.ws)
	defer socket.LeaveNode(userAuth.UID)


	for {
		_, request, err := this.ws.ReadMessage()

		if err != nil {
			logger.Error("MessageReadFailed|SID=", userAuth.SID, "|UID=", userAuth.UID, "|Err=%v", err)
			return
		}
		logger.Debug("REQUEST|Sid=", userAuth.UID, "UID=", userAuth.UID, "|WSRequest=", string(request))
		response, _ := serve(request, userAuth, &location)

		if (response != nil) {
			logger.Debug("RESPONSE|Sid=", userAuth.SID, "|UID=", userAuth.UID, "|Response=%v", response)
			this.Respond(response)
		}

	}
}

func serve(requestBody []byte, userAuth *models.UserAuth, location *domain.Coordinate) (interface{}, error) {
	message := string(requestBody)

	switch  {
	case strings.HasPrefix(message, socket.UpdateLocation):
		var newLocation = domain.Coordinate{}
		_ = json.Unmarshal(requestBody[len(socket.UpdateLocation):], &newLocation)
		updateLocation(location, &newLocation, userAuth)

		return nil, nil

	case strings.HasPrefix(message, socket.SendMessage):

		var msg = domain.MessageSendRequest{}
		_ = json.Unmarshal(requestBody[len(socket.SendMessage):], &msg)
		updateLocation(location, &msg.Location, userAuth)
		if msg.Location.Longitude == 0.0 || msg.Location.Latitude == 0.0 {
			msg.Location = *location
		}
		err := handleMessage(userAuth.UID, &msg)
		if err != nil {
			logger.Critical("MessageRecv|NotHandled|SID=", userAuth.SID, "|UID=", userAuth.UID, "|Err=%v", err)
			return serveMessageRecvErr(msg), nil
		} else {
			return serveMessageRecvOk(msg), nil
		}
	case strings.HasPrefix(message, socket.GroupMessage):

		var msg = domain.MessageSendRequest{}
		_ = json.Unmarshal(requestBody[len(socket.SendMessage):], &msg)

		err := handleMessage(userAuth.UID, &msg)
		if err != nil {
			logger.Critical("MessageRecv|NotHandled|SID=", userAuth.SID, "|UID=", userAuth.UID, "|Err=%v", err)
			return serveMessageRecvErr(msg), nil
		} else {
			return serveMessageRecvOk(msg), nil
		}
	case strings.HasPrefix(message, socket.TypingMessage):

		var msg = domain.MessageTypingRequest{}
		_ = json.Unmarshal(requestBody[len(socket.TypingMessage):], &msg)

		handleTyping(userAuth.UID, &msg)
		return nil, nil
	default:
		logger.Critical("UnknownCommand|NotHandled|SID=", userAuth.SID, "|UID=", userAuth.UID)
		return nil, nil
	}

}

func handleTyping(uid string, msg *domain.MessageTypingRequest) {
	if msg.GID == "" {
		if msg.IsTyping {
			go socket.MulticastPerson(uid, "typing")
		} else {
			go socket.MulticastPerson(uid, "nottyping")
		}
	} else {
		userAccount, _ := models.GetUserAccount("uid", uid)
		userLocation, _ := models.GetUserLocation(uid)
		userGroup, _ := models.GetGroupAccount(msg.GID)
		if msg.IsTyping {
			go socket.MulticastPersonCustom("typing", userAccount, userLocation.Location, userGroup.UIDS, userGroup.GID)
		} else {
			go socket.MulticastPersonCustom("nottyping", userAccount, userLocation.Location, userGroup.UIDS, userGroup.GID)
		}
	}
}

func updateLocation(oldLocation, newLocation *domain.Coordinate, userAuth *models.UserAuth) {
	if newLocation.Latitude == 0.0 || newLocation.Longitude == 0.0 {
		if oldLocation.Longitude != 0.0 && oldLocation.Latitude != 0.0 {
			newLocation = oldLocation
		}
		return
	}

	if oldLocation == nil || (oldLocation != nil && common.Distance(oldLocation, newLocation) > 5000) {
		if oldLocation == nil {
			oldLocation = &domain.Coordinate{}
		}
		oldLocation.Longitude = newLocation.Longitude
		oldLocation.Latitude = newLocation.Latitude
		go models.AddUserNewLocation(newLocation.Longitude, newLocation.Latitude, userAuth.UID)
	}
	oldLocation = newLocation
}

func handleMessage(uid string, msg *domain.MessageSendRequest) (error) {

	_, err := models.CreateMessage(uid, msg.MID, msg.Location.Longitude, msg.Location.Latitude, msg.Body)

	if err != nil {
		return err
	}
	socket.AddToFeed(uid, msg.GID, msg)
	userAccount, err := models.GetUserAccount("uid", uid)
	if err != nil {
		return err
	}

	go socket.MulticastMessage(userAccount, msg)
	return nil
}


func serveMessageRecvErr(msg domain.MessageSendRequest) *domain.MessageReceivedResponse {
	return &domain.MessageReceivedResponse{
		StatusCode:500,
		RequestId:msg.MID,
		GID: msg.GID,
	}
}

func serveMessageRecvOk(msg domain.MessageSendRequest) *domain.MessageReceivedResponse {
	return &domain.MessageReceivedResponse{
		StatusCode:200,
		RequestId:msg.MID,
		GID: msg.GID,
	}
}
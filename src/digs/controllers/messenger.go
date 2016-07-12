package controllers

import (
	"digs/domain"
	"strings"
	"digs/models"
	"encoding/json"
	"digs/socket"
	"digs/common"
	"digs/logger"
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

	socket.AddNode(userAuth.UID, this.ws)
	defer socket.LeaveNode(userAuth.UID)
	var location = domain.Coordinate{}

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

		err := handleMessage(userAuth.UID, &msg)
		if err != nil {
			logger.Critical("MessageRecv|NotHandled|SID=", userAuth.SID, "|UID=", userAuth.UID, "|Err=%v", err)
			return &domain.MessageReceivedResponse{
				StatusCode:500,
				RequestId:msg.MID,
			}, nil
		} else {
			return &domain.MessageReceivedResponse{
				StatusCode:200,
				RequestId:msg.MID,
			}, nil
		}
	default:
		logger.Critical("UnknownCommand|NotHandled|SID=", userAuth.SID, "|UID=", userAuth.UID)
		return nil, nil
	}

}

func updateLocation(oldLocation, newLocation *domain.Coordinate, userAuth *models.UserAuth) {

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
	models.AddToUserFeed(uid, msg.MID)

	userAccount, err := models.GetUserAccount("uid", uid)
	if err != nil {
		return err
	}

	go socket.MulticastMessage(userAccount, msg)
	return nil
}


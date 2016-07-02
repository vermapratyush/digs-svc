package controllers

import (
	"digs/domain"
	"strings"
	"digs/models"
	"encoding/json"
	"digs/socket"
	"github.com/astaxie/beego"
)

type WSMessengerController struct {
	WSBaseController
}

func (this *WSMessengerController) Get() {
	sid := this.GetString("sessionId")

	beego.Info("WSConnection|SID=", sid)

	userAuth, err := models.FindSession("sid", sid)
	if err != nil {
		this.Respond(&domain.GenericResponse{
			StatusCode:422,
			MessageCode:5000,
			Message:"Invalid Session",
		})
		return
	}
	beego.Info("UserConnected|UID=", userAuth)
	this.Respond(&domain.GenericResponse{
		StatusCode: 200,
		Message: "Connection Established",
		MessageCode: 3000,
	})

	socket.AddNode(userAuth.UID, this.ws)
	defer socket.LeaveNode(userAuth.UID)

	for {
		_, request, err := this.ws.ReadMessage()

		if err != nil {
			beego.Info("Err", err.Error())
		}
		if err != nil {
			beego.Critical("NodeConnectionLost|Error", err)
			return
		}
		response, _ := serve(request, userAuth)
		if (response != nil) {
			beego.Info("From sid=", userAuth, "Response", response)
			this.Respond(response)
		}

	}
}

func serve(requestBody []byte, userAuth *models.UserAuth) (interface{}, error) {

	var location = domain.Coordinate{}
	message := string(requestBody)

	switch  {
	case strings.HasPrefix(message, socket.UpdateLocation):
		var newLocation = domain.Coordinate{}
		_ = json.Unmarshal(requestBody[len(socket.UpdateLocation):], &newLocation)
		updateLocation(&location, &newLocation, userAuth)

		return nil, nil

	case strings.HasPrefix(message, socket.SendMessage):
		beego.Info("From sid=", userAuth, "Request", string(requestBody))

		var msg = domain.MessageSendRequest{}
		_ = json.Unmarshal(requestBody[len(socket.SendMessage):], &msg)
		updateLocation(&location, &msg.Location, userAuth)
		beego.Info("SendMessage|Message=", msg)

		err := handleMessage(userAuth.UID, &msg)
		if err != nil {
			beego.Critical("Unable to handle message %s", err)
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
		beego.Warn("UnknownCommand|Request=", string(requestBody))
		return nil, nil
	}

}

func updateLocation(oldLocation, newLocation *domain.Coordinate, userAuth *models.UserAuth) {

	if oldLocation == nil || oldLocation.Longitude != newLocation.Longitude || oldLocation.Latitude != newLocation.Latitude {
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


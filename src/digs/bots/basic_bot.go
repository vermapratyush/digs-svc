package bots

import (
	"digs/models"
	"github.com/satori/go.uuid"
	"digs/logger"
)

//type Model interface{}
type Integration func(bot Bot, userAccount models.UserAccount, coordinate models.Coordinate) (interface{}, error)
type MessageGeneration func(models.UserAccount, models.Coordinate, []string, interface{}) string

type Bot struct {
	BotName string
	FromUser string
	Radius float64
	Users []models.UserAccount
	//BotModel Model
	BotIntegration Integration
	BotMessageGeneration MessageGeneration
	IBot
}

type IBot interface {
	//Initialize(botName, fromUser string, radius float64)
	GetUserProfile(uid string) *models.UserAccount
	GetUserLocation(uid string) *models.UserLocation
	NearByPeople(longitude, latitude, maxDistance float64) []string
	AddToUserFeed(string, *models.Message)
	CreateBotMessage(uid string, body string) *models.Message
	SendPushForMessage(uid string, message *models.Message)


	//Custom override the following
	//BotIntegration(models.UserAccount, models.Coordinate) (interface{}, error)
	//GenerateMessage(toUser models.UserAccount, toLocation models.Coordinate, nearBy []string, botData ...interface{}) string

	//Sending Strategy
	//SpiderWebStrategy() map[string]bool
}


func SpiderWebStrategy(botInfo Bot)  map[string]bool {
	processed := make(map[string]bool)
	for _, user := range(botInfo.Users)  {
		_, present := processed[user.UID]
		if !present {
			userLocation := botInfo.GetUserLocation(user.UID)
			if len(userLocation.Location.Coordinates) == 0 {
				continue
			}

			nearBy := botInfo.NearByPeople(userLocation.Location.Coordinates[0], userLocation.Location.Coordinates[1], botInfo.Radius)
			meetupList, err := botInfo.BotIntegration(botInfo, user, userLocation.Location)
			if (err == nil) {
				for _, near := range(nearBy) {
					_, present := processed[near]
					if present {
						continue
					}

					otherUser := botInfo.GetUserProfile(near)
					otherUserLocation := botInfo.GetUserLocation(otherUser.UID)
					if len(userLocation.Location.Coordinates) == 0 {
						continue
					}
					body := botInfo.BotMessageGeneration(*otherUser, otherUserLocation.Location, nearBy, meetupList)
					message := botInfo.CreateBotMessage(otherUser.UID, body)

					botInfo.AddToUserFeed(otherUser.UID, message)
					botInfo.SendPushForMessage(otherUser.UID, message)
					processed[otherUser.UID] = true
					logger.Debug("SponsorNotification|User=", user.UID, "|Message=", message)
				}
			} else {
				logger.Debug("NoNotification|User=", user.UID)
			}
			processed[user.UID] = true
		}
	}
	return processed
}

func (this Bot) GetUserProfile(uid string) *models.UserAccount {
	userAccount, _ := models.GetUserAccount("uid", uid)
	return userAccount
}

func (this Bot) GetUserLocation(uid string) *models.UserLocation {
	location, _ := models.GetUserLocation(uid)
	return &location
}

func (this Bot) NearByPeople(longitude, latitude, maxDistance float64) []string {
	uids := models.GetLiveUIDForFeed(longitude, latitude, maxDistance, -1)
	return uids
}

func (this Bot) SendPushForMessage(uid string, message *models.Message) {
	devices, _ := models.GetNotificationIds(uid)
	for _, device := range(*devices) {
		if device.OSType == "android" {
			models.AndroidMessagePush(uid, device.NotificationId, message.Content, "", "individual", "", this.BotName)
		} else {
			models.IOSMessagePush(uid, device.NotificationId, message.Content)
		}
	}
}

func (this Bot) CreateBotMessage(uid string, body string) *models.Message {
	uuid := this.BotName + uuid.NewV4().String()
	message, _ := models.CreateMessage(this.FromUser, uuid, 0.0, 0.0, body)
	return message
}

func (this Bot) AddToUserFeed(uid string, message *models.Message) {
	models.AddToUserFeed(uid, message.MID)
}
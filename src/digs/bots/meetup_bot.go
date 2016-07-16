package bots

import (
	"digs/models"
	"github.com/satori/go.uuid"
	"digs/controllers"
	"fmt"
	"sort"
	"strings"
	"digs/common"
	"digs/logger"
)

var (
	BotName = "MeetupBot"
	FromUser = "PowowInfo"
	DefaultRadius = 30000.0
)
type MeetupBotcontroller struct {
	controllers.HttpBaseController
}

func (this *MeetupBotcontroller) Get() {
	userId := this.GetString("userId")
	users := make([]models.UserAccount, 0)
	if userId != "" {
		userAccount, _ := models.GetUserAccount("uid", userId)
		users = append(users, *userAccount)
	} else {
		users, _ = models.GetAllUserAccount()
	}

	processed := make(map[string]struct{})
	for _, user := range(users)  {
		_, present := processed[user.UID]
		if !present {
			userLocation := GetUserLocation(user.UID)
			if len(userLocation.Location.Coordinates) == 0 {
				continue
			}

			nearBy := NearByPeople(userLocation.Location.Coordinates[0], userLocation.Location.Coordinates[1], DefaultRadius)
			meetupList := GetTopMeetup(userLocation.Location)
			if (len(meetupList) > 0) {
				for _, near := range(nearBy) {
					_, present := processed[near]
					if present {
						continue
					}

					otherUser := GetUserProfile(near)
					otherUserLocation := GetUserLocation(otherUser.UID)
					if len(userLocation.Location.Coordinates) == 0 {
						continue
					}
					message := GenerateMessage(*otherUser, otherUserLocation.Location, nearBy, meetupList)

					AddToUserFeed(otherUser.UID, message)
					SendPushForMessage(otherUser.UID, message)
					processed[otherUser.UID] = struct {}{}
					logger.Debug("SponsorNotification|User=", user.UID, "|Message=", message)
				}
			} else {
				logger.Debug("NoNotification|User=", user.UID)
			}
			processed[user.UID] = struct {}{}
		}
	}
}

func GetUserProfile(uid string) *models.UserAccount {
	userAccount, _ := models.GetUserAccount("uid", uid)
	return userAccount
}

func GetUserLocation(uid string) *models.UserLocation {
	location, _ := models.GetUserLocation(uid)
	return &location
}

func NearByPeople(longitude, latitude, maxDistance float64) []string {
	uids := models.GetLiveUIDForFeed(longitude, latitude, maxDistance, -1)
	return uids
}

func GenerateMessage(user models.UserAccount, location models.Coordinate, nearBy []string, meetupList []models.MeetupEvent) *models.Message {

	body := ""
	if len (nearBy) == 1 {
		body = fmt.Sprintf("Hey %s, We found few meetups for tomorrow which you might like to try out with people nearby. Why don't you give it a shot!.<br/>", user.FirstName)
	} else {
		body = fmt.Sprintf("Hey %s, We found few meetups for tomorrow which you might like to try out with people nearby. Why don't you club together with %d other powow users in your area and go to one of the events.<br/>", user.FirstName, len(nearBy))
	}

	for idx, meetup := range (meetupList) {
		body = fmt.Sprintf("%s%d. <a href=\"%s\">%s</a><br/>", body, (idx + 1),  meetup.EventUrl, meetup.Title)
	}
	message := CreateBotMessage(user.UID, body)
	return message
}

func GetTopMeetup(location models.Coordinate) []models.MeetupEvent {

	events := models.GetMeetup(location.Coordinates[0], location.Coordinates[1], DefaultRadius)
	sort.Sort(models.ByYesCount{events.Results})
	eventsWithinRange := make([]models.MeetupEvent, 0)
	for _, event := range(events.Results) {
		distFromMeetup := common.DistanceLong(location.Coordinates[0], location.Coordinates[1],
			event.EventLatLong.Longitude, event.EventLatLong.Latitude)
		if (distFromMeetup < DefaultRadius) {
			eventsWithinRange = append(eventsWithinRange, event)
		}
	}
	finalList := make([]models.MeetupEvent, 3)
	if (len(eventsWithinRange) > 3) {
		finalList = eventsWithinRange[0:3]
	} else {
		finalList = eventsWithinRange
	}
	for idx, event := range(finalList) {
		short := models.ShortenUrl(event.EventUrl)
		short = strings.Replace(short, "http://", "http://www.", 1)
		finalList[idx].EventUrl = short
	}

	return finalList
}

func SendPushForMessage(uid string, message *models.Message) {
	devices, _ := models.GetNotificationIds(uid)
	for _, device := range(*devices) {
		if device.OSType == "android" {
			models.AndroidMessagePush(uid, device.NotificationId, message.Content)
		} else {
			models.IOSMessagePush(uid, device.NotificationId, message.Content)
		}
	}
}

func CreateBotMessage(uid string, body string) *models.Message {
	uuid := BotName + uuid.NewV4().String()
	message, _ := models.CreateMessage(FromUser, uuid, 0.0, 0.0, body)
	return message
}

func AddToUserFeed(uid string, message *models.Message) {
	models.AddToUserFeed(uid, message.MID)
}
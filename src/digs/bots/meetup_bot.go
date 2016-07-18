package bots

import (
	"digs/models"
	"digs/controllers"
	"fmt"
	"sort"
	"strings"
	"digs/common"
	"errors"
)

type MeetupBotController struct {
	controllers.HttpBaseController
}

type MeetupBot Bot

func (this *MeetupBotController) Post() {

	userId := this.GetString("userId")
	users := make([]models.UserAccount, 0)
	if userId != "" {
		userAccount, _ := models.GetUserAccount("uid", userId)
		users = append(users, *userAccount)
	} else {
		users, _ = models.GetAllUserAccount()
	}

	botInfo := MeetupBot{
		BotName: "MeetupBot",
		FromUser: "PowowInfo",
		Users:users,
		Radius:10000.0,
		BotIntegration:meetupBotIntegration,
		BotMessageGeneration:meetupGenerateMessage,
	}

	processed := SpiderWebStrategy(Bot(botInfo))

	this.Serve200(processed)
}


func meetupGenerateMessage(toUser models.UserAccount, toLocation models.Coordinate, nearBy []string, meetupList interface{}) string {
	body := ""
	if len (nearBy) == 1 {
		body = fmt.Sprintf("Hey %s, We found few meetups for tomorrow which you might like to try out with people nearby. Why don't you give it a shot!.<br/>", toUser.FirstName)
	} else {
		body = fmt.Sprintf("Hey %s, We found few meetups for tomorrow which you might like to try out with people nearby. Why don't you club together with %d other powow users in your area and go to one of the events.<br/>", toUser.FirstName, len(nearBy))
	}

	meetupTypecast := meetupList.([]MeetupEvent)
	for idx, meetup := range (meetupTypecast) {
		body = fmt.Sprintf("%s%d. <a href=\"%s\">%s</a><br/>", body, (idx + 1),  meetup.EventUrl, meetup.Title)
	}
	return body
}

func meetupBotIntegration(bot Bot, toUser models.UserAccount, location models.Coordinate) (interface{}, error) {

	events := GetMeetup(location.Coordinates[0], location.Coordinates[1], bot.Radius)
	sort.Sort(ByYesCount{events.Results})
	eventsWithinRange := make([]MeetupEvent, 0)
	for _, event := range(events.Results) {
		distFromMeetup := common.DistanceLong(location.Coordinates[0], location.Coordinates[1],
			event.EventLatLong.Longitude, event.EventLatLong.Latitude)
		if (distFromMeetup < bot.Radius) {
			eventsWithinRange = append(eventsWithinRange, event)
		}
	}
	finalList := make([]MeetupEvent, 3)
	if (len(eventsWithinRange) > 3) {
		finalList = eventsWithinRange[0:3]
	} else {
		finalList = eventsWithinRange
	}
	for idx, event := range(finalList) {
		short := ShortenUrl(event.EventUrl)
		short = strings.Replace(short, "http://", "http://www.", 1)
		finalList[idx].EventUrl = short
	}

	if (len(finalList) == 0) {
		return nil, errors.New("Not enough meetups")
	}
	return finalList, nil
}


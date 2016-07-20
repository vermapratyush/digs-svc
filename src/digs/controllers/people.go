package controllers

import (
	"digs/models"
	"errors"
	"digs/domain"
	"digs/common"
	"digs/socket"
	"gopkg.in/mgo.v2"
	"digs/logger"
	"strconv"
)

type PeopleController struct {
	HttpBaseController
}

func (this *PeopleController) Get() {
	sid := this.GetString("sessionId")
	longitude, longErr := this.GetFloat("longitude")
	latitude, latErr := this.GetFloat("latitude")

	userAuth, err := models.FindSession("sid", sid)

	if longErr != nil || latErr != nil {
		logger.Error("PEOPLE|LatLongFormatError|SID=", userAuth.SID, "|UID=", userAuth.UID, "|Lat=", latitude, "|Long=", longitude, "|Err=%v", err)
		this.Serve500(errors.New("Location cordinate not provided in proper format"))
		return
	} else {
		latFloat, _ := strconv.ParseFloat(this.GetString("latitude"), 64)
		longFloat, _ := strconv.ParseFloat(this.GetString("longitude"), 64)
		models.AddUserNewLocation(longFloat, latFloat, userAuth.UID)
	}

	if err != nil {
		if err == mgo.ErrNotFound {
			this.InvalidSessionResponse()
			return
		}
		this.Serve500(err)
		logger.Error("PEOPLE|SessionRetrieveError|SID=", userAuth.SID, "|UID=", userAuth.UID, "|Lat=", latitude, "|Long=", longitude, "|Err=%v", err)
		return
	}
	userAccount, err := models.GetUserAccount("uid", userAuth.UID)
	if err != nil {
		logger.Error("PEOPLE|UserNotFound|SID=", userAuth.SID, "|UID=", userAuth.UID, "|Lat=", latitude, "|Long=", longitude, "|Err=%v", err)
		this.Serve500(errors.New("User not found"))
		return
	}

	uidList := models.GetLiveUIDForFeed(longitude, latitude, userAccount.Settings.Range, -1)
	users, err := models.GetAllUserAccountIn(uidList)
	if err != nil {
		logger.Error("PEOPLE|GetUserAccountFailed|SID=", userAuth.SID, "|UID=", userAuth.UID, "|Lat=", latitude, "|Long=", longitude, "|Err=%v", err)
		return
	}

	blockedMap := common.GetStringArrayAsMap(userAccount.BlockedUsers)

	//TODO: Find a better solution, too make realloc
	people := make([]domain.PersonResponse, 0, len(uidList))
	for idx := 0; idx < len(users); idx = idx + 1 {
		user := users[idx]
		_, present := socket.GetLookUp(user.UID)
		_, presentInBlock := blockedMap[user.UID]

		if user.UID == userAccount.UID || presentInBlock  {
			continue
		}

		activeState := "active"
		if !present {
			activeState = "inactive"
		}
		people = append(people, domain.PersonResponse{
			Name: common.GetName(user.FirstName, user.LastName),
			UID: user.UID,
			About: user.About,
			Activity: "join",
			ActiveState: activeState,
			ProfilePicture: user.ProfilePicture,
		})
	}

	people = addPeopleWhoCommunicatedOneOnOne(userAuth.UID, people)

	//addAlwaysActiveBot(people)

	logger.Debug("PEOPLE|SID=", userAuth.SID, "|UID=", userAuth.UID, "|FeedSize=", len(people))

	this.Serve200(people)
}

func addPeopleWhoCommunicatedOneOnOne(uid string, people []domain.PersonResponse) []domain.PersonResponse {
	oneOneOne, _ := models.GetGroupsUserIsMemberOf(uid)

	for _, user := range (oneOneOne) {
		addUser := true
		for _, person := range(people) {
			if user.UserAccount.UID == person.UID {
				addUser = false
			}
		}
		if addUser {
			logger.Debug("adding")
			people = append(people, domain.PersonResponse {
				Name: common.GetName(user.UserAccount.FirstName, user.UserAccount.LastName),
				UID: user.UserAccount.UID,
				About: user.UserAccount.About,
				Activity: "join",
				ActiveState: "out_of_range",
				ProfilePicture: user.UserAccount.ProfilePicture,
			})
		}
	}
	return people
}

func addAlwaysActiveBot(people []domain.PersonResponse) {
	containsSitesh := false
	containsPratyush := false
	for _, person := range(people) {
		if person.UID == "10210146992256811" {
			containsSitesh = true
			person.ActiveState = "active"
		} else if person.UID == "10154168592560450" {
			containsPratyush = true
			person.ActiveState = "active"
		}
	}
	user, _ := models.GetUserAccount("uid", "10210146992256811")
	if !containsSitesh {
		people = append(people, domain.PersonResponse{
			Name: common.GetName(user.FirstName, user.LastName),
			UID: user.UID,
			About: user.About,
			Activity: "join",
			ActiveState: "active",
			ProfilePicture: user.ProfilePicture,
		})
	}
	user, _ = models.GetUserAccount("uid", "10154168592560450")
	if !containsPratyush {
		people = append(people, domain.PersonResponse{
			Name: common.GetName(user.FirstName, user.LastName),
			UID: user.UID,
			About: user.About,
			Activity: "join",
			ActiveState: "active",
			ProfilePicture: user.ProfilePicture,
		})
	}
}




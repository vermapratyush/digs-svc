package controllers

import (
	"errors"
	"digs/domain"
	"digs/common"
	"digs/socket"
	"gopkg.in/mgo.v2"
	"digs/logger"
	"strconv"
	"github.com/deckarep/golang-set"
	"fmt"
	"digs/models"
	"github.com/astaxie/beego"
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
	blockedGroup := common.GetStringArrayAsMap(userAccount.BlockedGroups)

	//TODO: Find a better solution, too make realloc
	people := make([]*domain.PersonResponse, 0, len(uidList))
	for idx := 0; idx < len(users); idx = idx + 1 {
		user := users[idx]
		_, present := socket.GetLookUp(user.UID)
		_, presentInBlock := blockedMap[user.UID]

		if user.UID == userAccount.UID || presentInBlock || !user.Settings.PublicProfile  {
			continue
		}

		activeState := "active"
		if !present {
			activeState = "inactive"
		}
		people = append(people, &domain.PersonResponse{
			Name: common.GetName(user.FirstName, user.LastName),
			UID: user.UID,
			About: user.About,
			Activity: "join",
			ActiveState: activeState,
			Verified:user.Verified,
			ProfilePicture: user.ProfilePicture,
			IsGroup: false,
		})
	}

	//TODO: HACK: GET RID
	version := this.Ctx.Input.Param(":version")
	if version == "v1" {

		people = addPeopleWhoCommunicatedOneOnOneHack(userAccount.UID, people[0:], blockedMap)
		people = addUnreadCount(userAuth.UID, people[0:])
	} else {
		people = addPeopleWhoCommunicatedOneOnOne(userAccount, people[0:], blockedMap)
		people = addJoinedGroups(userAccount, people[0:])
		people = addGroupsNearBy(userAccount, &domain.Coordinate{Longitude:longitude, Latitude:latitude}, blockedGroup, people[0:])
		people = addFourSquareGroups(userAccount, &domain.Coordinate{Longitude:longitude, Latitude:latitude}, blockedGroup, people[0:])
	}

	//addAlwaysActiveBot(people)

	logger.Debug("PEOPLE|SID=", userAuth.SID, "|UID=", userAuth.UID, "|FeedSize=", len(people))

	this.Serve200(people)
}

func addFourSquareGroups(userAccount *models.UserAccount, coordinate *domain.Coordinate, blockedGroup map[string]struct{}, people []*domain.PersonResponse) []*domain.PersonResponse {

	fourSquare := models.SearchFourSquareVenue(coordinate.Longitude, coordinate.Latitude, 150.0)

	for _, venue := range(fourSquare.Response.Venues) {
		member := false
		for _, person := range(people) {
			if person.GID == "foursquare-" + venue.Id {
				member = true
			}
		}
		icon := "https://ss3.4sqi.net/img/categories_v2/none_bg_64.png"
		if len (venue.Categories) > 0 {
			icon = venue.Categories[0].CategoryIcon.Prefix + "bg_64" + venue.Categories[0].CategoryIcon.Suffix
		}
		beego.Info(icon)
		if _, present := blockedGroup[venue.Id]; !member && !present {
			people = append(people, &domain.PersonResponse{
				Name: venue.Name,
				GID: fmt.Sprintf("foursquare-%s", venue.Id),
				About: "Location imported from FourSquare",
				ActiveState: "nearby_group",
				UnreadCount: 0,
				MemberCount: int(venue.VenueStats.UsersCount),
				IsGroup: true,
				ProfilePicture: icon,
			})
		}
	}
	return people
}

func addJoinedGroups(userAccount *models.UserAccount, people []*domain.PersonResponse) []*domain.PersonResponse {
	userGroups := models.GetUserGroups(userAccount.MultiGroupId())
	for _, group := range(userGroups) {
		unread := common.IndexOf(group.MIDS, group.MessageRead[userAccount.UID])
		if unread < 0 {
			unread = len(group.MIDS);
		}
		people = append(people, &domain.PersonResponse{
			Name: group.GroupName,
			GID: group.GID,
			About: group.GroupAbout,
			ActiveState: "joined_group",
			UnreadCount: int64(unread),
			MemberCount: len(group.UIDS),
			IsGroup: true,
			ProfilePicture: group.GroupPicture,
		})
	}

	return people
}

func addGroupsNearBy(userAccount *models.UserAccount, coordinate *domain.Coordinate, blockedGroup map[string]struct{}, people []*domain.PersonResponse) []*domain.PersonResponse {
	nearByPeople := models.GetLiveUIDForFeed(coordinate.Longitude, coordinate.Latitude, userAccount.Settings.Range, -1)
	userAccounts, _ := models.GetAllUserAccountIn(nearByPeople)
	groupIds := mapset.NewSet()
	for _, userAccount := range (userAccounts) {
		for _, gid := range (userAccount.MultiGroupId()) {
			_, blocked := blockedGroup[gid]
			if !blocked {
				groupIds.Add(gid)
			}
		}
	}
	//TODO: HACK: Add AlwaysShow Groups
	groupIds.Add("5f839ce5-5894-4de0-8b78-8149a5febdd5")
	
	for _, group := range(people) {
		if group.ActiveState == "joined_group" {
			groupIds.Remove(group.GID)
		}
	}
	groupIdString := make([]string, groupIds.Cardinality())
	idx := 0
	for gid := range(groupIds.Iter()) {
		groupIdString[idx] = gid.(string)
		idx++
	}
	userGroups := models.GetUserGroups(groupIdString)
	for _, group := range(userGroups) {
		if len(group.MIDS) == 0 {
			continue
		}
		people = append(people, &domain.PersonResponse{
			Name: group.GroupName,
			GID: group.GID,
			About: group.GroupAbout,
			ActiveState: "nearby_group",
			UnreadCount: 0,
			MemberCount: len(group.UIDS),
			IsGroup: true,
			ProfilePicture: group.GroupPicture,
		})
	}
	return people
}

//TODO: HACK: GET RID
func addUnreadCount(uid string, people []*domain.PersonResponse) []*domain.PersonResponse {

	userGroups, _ := models.GetUnreadMessageCountOneOnOne(uid)
	userGroupMap := map[string]models.UserGroup{}
	for _, userGroup := range (userGroups) {
		otherUserInGroup := userGroup.UIDS[0]
		if otherUserInGroup == uid {
			otherUserInGroup = userGroup.UIDS[1]
		}
		userGroupMap[otherUserInGroup] = userGroup
	}

	for _, person := range (people) {
		idx := common.IndexOf(userGroupMap[person.UID].MIDS, userGroupMap[person.UID].MessageRead[uid])
		if idx < 0 {
			idx = 0
		}
		person.UnreadCount = int64(idx)
	}

	return people
}

func addPeopleWhoCommunicatedOneOnOneHack(uid string, people []*domain.PersonResponse, blockedMap map[string]struct{}) []*domain.PersonResponse {
	oneOneOne, _ := models.GetGroupsUserIsMemberOf(uid)

	for _, user := range (oneOneOne) {
		_, presentInBlock := blockedMap[user.UserAccount.UID]
		if len(user.MessageIds) == 0 || presentInBlock {
			continue
		}
		addUser := true
		for _, person := range(people) {
			if user.UserAccount.UID == person.UID {
				addUser = false
			}
		}
		if addUser {
			people = append(people, &domain.PersonResponse {
				Name: common.GetName(user.UserAccount.FirstName, user.UserAccount.LastName),
				UID: user.UserAccount.UID,
				About: user.UserAccount.About,
				Activity: "join",
				ActiveState: "out_of_range",
				Verified:user.UserAccount.Verified,
				ProfilePicture: user.UserAccount.ProfilePicture,
			})
		}
	}
	return people
}

func addPeopleWhoCommunicatedOneOnOne(userAccount *models.UserAccount, people []*domain.PersonResponse, blockedMap map[string]struct{}) []*domain.PersonResponse {

	userGroups := models.GetUserGroups(userAccount.OneToOneGroupId())

	userIdsForOneOnOne := make([]string, 0)
	unreadCountPerGroup := make(map[string]int64, len(userGroups))
	for _, group := range(userGroups) {
		if len(group.MIDS) > 0 {
			unread := common.IndexOf(group.MIDS, group.MessageRead[userAccount.UID])
			if unread < 0 {
				unread = len(group.MIDS);
			}
			if group.UIDS[0] != userAccount.UID {
				userIdsForOneOnOne = append(userIdsForOneOnOne, group.UIDS[0])
				unreadCountPerGroup[group.UIDS[0]] = int64(unread)
			} else {
				userIdsForOneOnOne = append(userIdsForOneOnOne, group.UIDS[1])
				unreadCountPerGroup[group.UIDS[1]] = int64(unread)
			}

		}
	}

	if len(userIdsForOneOnOne) > 0 {
		oneOnOne, _ := models.GetAllUserAccountIn(userIdsForOneOnOne)
		for _, user := range (oneOnOne) {

			_, presentInBlock := blockedMap[user.UID]
			if presentInBlock || !user.Settings.PublicProfile {
				continue
			}
			addUser := true
			for _, person := range(people) {
				if user.UID == person.UID {
					person.UnreadCount = unreadCountPerGroup[person.UID]
					addUser = false
				}
			}
			if addUser {
				people = append(people, &domain.PersonResponse {
					Name: common.GetName(user.FirstName, user.LastName),
					UID: user.UID,
					About: user.About,
					Activity: "join",
					ActiveState: "out_of_range",
					Verified:user.Verified,
					ProfilePicture: user.ProfilePicture,
					IsGroup: false,
					UnreadCount: unreadCountPerGroup[user.UID],
				})
			}
		}

	}

	return people
}
package models

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"time"
	"gopkg.in/mgo.v2/bson"
)

var DefaultDatabase = "heroku_qnx0661v"
var Session, _ = mgo.Dial(fmt.Sprintf("mongodb://localhost:27017/%s", DefaultDatabase))

//Create index MID, From
//Do Execute db.messages.ensureIndex({location:"2dsphere"})
type Coordinate struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

type Message struct {
	MID          string `bson:"mid" json:"mid"`
	From         string `bson:"from" json:"from"`
	Location     Coordinate `bson:"location" json:"location"`
	Content      string `bson:"content" json:"content"`
	CreationTime time.Time `bson:"creationTime" json:"creationTime"`
}

//Create Index UID, email
//CreationTime Asc
type UserAccount struct {
	UID             string `bson:"uid" json:"uid"`
	FirstName       string `bson:"firstName" json:"firstName"`
	LastName        string `bson:"lastName" json:"lastName"`
	Email           string `bson:"email" json:"email"`
	About           string `bson:"about" json:"about"`
	ProfilePicture  string `bson:"profilePicture" json:"profilePicture"`
	FBID            string `json:"fbid" bson:"fbid"`
	Locale          string `json:"locale" bson:"locale"`
	CreationTime    time.Time `bson:"creationTime" json:"creationTime"`
	FBVerified      bool `bson:"fbVerified" json:"fbVerified"`
	Settings        Setting `json:"settings" bson:"settings"`
	BlockedUsers    []string `json:"blockedUsers" bson:"blockedUsers"`
	BlockedMessages []string `json:"blockedMessages" bson:"blockedMessages"`
}

type Setting struct {
	Range            float64 `json:"messageRange" bson:"messageRange"`
	PublicProfile    bool `json:"publicProfile" bson:"publicProfile"`
	PushNotification bool `json:"enableNotification" bson:"enableNotification"`
}

//Order by creation time asc
//Index uid, sid, accessToken
type UserAuth struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UID          string `bson:"uid" json:"uid"`
	SID          string `bson:"sid" json:"sid"`
	AccessToken  string `bson:"accessToken" json:"accessToken"`
	CreationTime time.Time `bson:"creationTime" json:"creationTime"`
}

type UserLocation struct {
	UID          string `json:"uid" bson:"uid"`
	Location     Coordinate `json:"location" bson:"location"`
	MessageRange float64 `json:"messageRange" bson:"messageRange"`
	CreationTime time.Time `bson:"creationTime" json:"creationTime"`
}

type Notification struct {
	UID            string `json:"uid" bson:"uid"`
	NotificationId string `json:"notificationId" bson:"notificationId"`
	CreationTime   time.Time `bson:"creationTime" json:"creationTime"`
	OSType         string `bson:"os" json:"os"`
}

type MessageHistory struct {
	MID []string `json:"mid" bson:"mid"`
	UID string `json:"uid" bson:"uid"`
}

type UserGroup struct {
	GID        string `json:"gid" bson:"gid"`
	GroupName  string `json:"groupName" bson:"groupName"`
	GroupAbout string `json:"groupAbout" bson:"groupAbout"`
	UIDS       []string `json:"uids" bson:"uids"`
	MIDS       []string `json:"mids" bson:"mids"`
}

type UserGroupMessageResolved struct {
	MID          string `json:"mid" bson:"mid"`
	UID          string `json:"uid" bson:"uid"`
	Content      string `json:"content" bson:"content"`
	UserAccount  UserAccount `json:"userAccount" bson:"userAccount"`
	CreationTime time.Time `json:"creationTime" bson:"creationTime"`
}
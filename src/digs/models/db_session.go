package models

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"time"
	"gopkg.in/mgo.v2/bson"
)

var DefaultDatabase = "heroku_qnx0661v"
var Session, _ = mgo.Dial(fmt.Sprintf("mongodb://node-js:node-js@ds015194.mlab.com:15194/%s", DefaultDatabase))

//Create index MID, From
//Do Execute db.messages.ensureIndex({location:"2dsphere"})


type Coordinate struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates"`
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
	UID string `bson:"uid" json:"uid"`
	FirstName string `bson:"firstName" json:"firstName"`
	LastName string `bson:"lastName" json:"lastName"`
	Email string `bson:"email" json:"email"`
	About string `bson:"about" json:"about"`
	ProfilePicture string `bson:"profilePicture" json:"profilePicture"`
	FBID string `json:"fbid" bson:"fbid"`
	Locale string `json:"locale" bson:"locale"`
	CreationTime time.Time `bson:"creationTime" json:"creationTime"`
	FBVerified bool `bson:"fbVerified" json:"fbVerified"`
}

//Order by creation time asc
//Index uid, sid, accessToken
type UserAuth struct {
	Id bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UID string `bson:"uid" json:"uid"`
	SID string `bson:"sid" json:"sid"`
	AccessToken string `bson:"accessToken" json:"accessToken"`
	CreationTime time.Time `bson:"creationTime" json:"creationTime"`
}


type UserLocation struct {
	UID string `json:"uid" bson:"uid"`
	Username string `json:"username" bson:"username"`
	Location Coordinate `json:"location" bson:"location"`
	CreationTime time.Time `bson:"creationTime" json:"creationTime"`
}

type Notification struct {
	UID string `json:"uid" bson:"uid"`
	NotificationId string `json:"notificationId" bson:"notificationId"`
	CreationTime time.Time `bson:"creationTime" json:"creationTime"`
	OSType string `bson:"osType" json:"osType"`
}
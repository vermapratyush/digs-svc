package models

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"time"
)

var DefaultDatabase = "heroku_qnx0661v"
var Session, _ = mgo.Dial(fmt.Sprintf("mongodb://127.0.0.1:27017/%s", DefaultDatabase))

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
	CreationTime time.Time `bson:"creationTime" json:"creationTime"`
	ProfilePicture string `bson:"profilePicture" json:"profilePicture"`
	FBVerified string `bson:"fbVerified" json:"fbVerified"`
	FBID string `json:"fbid" bson:"fbid"`
	Locale string `json:"locale" bson:"locale"`
}

//Order by creation time asc
//Index uid, sid, accessToken
type UserAuth struct {
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
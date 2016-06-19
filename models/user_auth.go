package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

//Order by creation time desc
//Index uid, sid, accessToken
type UserAuth struct {
	UID string `bson:"uid" json:"uid"`
	SID string `bson:"sid" json:"sid"`
	AccessToken string `bson:"accessToken" json:"accessToken"`
	CreationTime time.Time `bson:"creationTime" json:"creationTime"`
}

func AddUserAuth(uid string, accessToken string, sid string) error {
	conn := Session.Clone()
	defer conn.Close()

	c := conn.DB("heroku_qnx0661v").C("auth")
	err := c.Insert(&UserAuth{
		UID: uid,
		SID: sid,
		AccessToken: accessToken,
		CreationTime:time.Now(),
	})
	return err
}

func FindSessionFromUID(uid string) string {
	conn := Session.Clone()
	defer conn.Close()

	c := conn.DB("heroku_qnx0661v").C("auth")
	res := UserAuth{}
	_ = c.Find(bson.M{"uid": uid}).Select(bson.M{"sid":"1"}).Sort("-creationTime").One(&res);
	return res.SID
}
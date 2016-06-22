package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

func AddUserAuth(uid string, accessToken string, sid string) error {
	conn := Session.Clone()
	defer conn.Close()

	c := conn.DB(DefaultDatabase).C("auth")
	err := c.Insert(&UserAuth{
		UID: uid,
		SID: sid,
		AccessToken: accessToken,
		CreationTime:time.Now(),
	})
	return err
}

func FindSession(fieldName, fieldValue string) (*UserAuth) {
	conn := Session.Clone()
	defer conn.Close()

	c := conn.DB(DefaultDatabase).C("auth")
	res := UserAuth{}
	_ = c.Find(bson.M{fieldName: fieldValue}).Sort("-creationTime").One(&res);
	return &res
}
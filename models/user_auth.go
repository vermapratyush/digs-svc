package models

import (
	"time"
)

type UserAuth struct {
	UID string `bson:"uid" json:"uid"`
	AccessToken string `bson:"accessToken" json:"accessToken"`
	CreationTime time.Time `bson:"creationTime" json:"creationTime"`
}

func AddUserAuth(uid string, accessToken string) error {
	conn := Session.Clone()
	defer conn.Close()

	c := conn.DB("heroku_qnx0661v").C("auth")
	err := c.Insert(&UserAuth{
		UID: uid,
		AccessToken: accessToken,
		CreationTime:time.Now(),
	})
	return err
}

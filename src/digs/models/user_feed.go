package models

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/astaxie/beego"
)

func AddToUserFeed(uid string, mid string) error {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("user_feed")
	defer conn.Close()

	query := bson.M{"uid": uid}
	update := bson.M{"$push": bson.M{"mid": mid }}

	change, err := c.Upsert(query, update)
	if err != nil {
		beego.Info(change)
	}
	return err
}

func GetUserFeed(uid string) (*MessageHistory, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("user_feed")
	defer conn.Close()

	var history MessageHistory
	err := c.Find(bson.M{"uid": uid}).One(&history)

	return &history, err
}

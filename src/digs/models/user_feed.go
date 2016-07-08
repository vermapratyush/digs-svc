package models

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/astaxie/beego"
	"github.com/afex/hystrix-go/hystrix"
	"digs/common"
)

func AddToUserFeed(uid string, mid string) error {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("user_feed")
	defer conn.Close()

	err := hystrix.Do(common.FeedAdd, func() error {
		query := bson.M{"uid": uid}
		update := bson.M{"$push": bson.M{"mid": mid }}

		change, err := c.Upsert(query, update)
		if err != nil {
			beego.Info(change)
		}
		return err
	}, nil)

	return err
}

func GetUserFeed(uid string) (*MessageHistory, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("user_feed")
	defer conn.Close()

	var history MessageHistory

	err := hystrix.Do(common.FeedGet, func() error {
		err := c.Find(bson.M{"uid": uid}).One(&history)
		return err
	}, nil)


	return &history, err
}

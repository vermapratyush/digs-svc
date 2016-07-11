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

func RemoveMessage(uid, mid string) error {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("user_feed")
	defer conn.Close()

	err := hystrix.Do(common.FeedDel, func() error {
		query := bson.M{"uid": uid}
		update := bson.M{
			"$pull": bson.M{
				"mid": mid,
			},
		}
		err := c.Update(query, update)
		if err != nil {
			beego.Error("Abuse|RemoveMessageFailed|err=", err)
		}
		return err
	}, nil)


	return err
}

func RemoveUserFromFeed(uid, abusiveUID string) error {
	conn := Session.Clone()
	messageConn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("user_feed")
	messageC := messageConn.DB(DefaultDatabase).C("messages")
	defer conn.Close()
	defer messageConn.Close()

	err := hystrix.Do(common.FeedDel, func() error {

		blockableMessages := []Message{}
		err := messageC.Find(bson.M{
			"from": abusiveUID,
		}).All(&blockableMessages)

		if err != nil {
			beego.Error("Abuse|RemoveMessageFailedBeforeFetch|err=", err)
			return err
		}

		mid := make([]string, len(blockableMessages))
		for idx, m := range(blockableMessages) {
			mid[idx] = m.MID
		}
		query := bson.M{"uid": uid}
		update := bson.M{
			"$pull": bson.M{
				"mid": bson.M{
					"$in": mid,
				},
			},
		}

		err = c.Update(query, update)
		if err != nil {
			beego.Error("Abuse|RemoveMessageFailedAfterFetch|err=", err)
		}
		return err


	}, nil)


	return err
}
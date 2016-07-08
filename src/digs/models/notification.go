package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"github.com/afex/hystrix-go/hystrix"
	"digs/common"
	"github.com/astaxie/beego"
)


func DeleteNotificationId(uid string, nid string) error {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("notifications")
	defer conn.Close()

	err := hystrix.Do(common.Notification, func() error {
		err := c.Remove(bson.M{"uid": uid, "notificationId": nid})
		if err != nil {
			beego.Error(err)
		}
		return err
	}, nil)

	return err
}

func AddNotificationId(uid string, nid string, os string) error {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("notifications")
	defer conn.Close()

	key := bson.M{"uid": uid, "notificationId": nid}
	value := bson.M{
		"uid": uid,
		"notificationId": nid,
		"os": os,
		"creationTime": time.Now(),
	}

	err := hystrix.Do(common.Notification, func() error {
		_, err := c.Upsert(key, value)
		if err != nil {
			beego.Error(err)
		}
		return err
	}, nil)

	return err
}

func GetNotificationIds(uid string) (*[]Notification, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("notifications")
	defer conn.Close()

	notifications := []Notification{}

	err := hystrix.Do(common.Notification, func() error {
		err := c.Find(bson.M{
			"uid": uid,
		}).All(&notifications)
		if err != nil {
			beego.Error(err)
		}
		return err
	}, nil)


	return &notifications, err

}
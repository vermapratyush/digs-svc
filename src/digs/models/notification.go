package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"github.com/afex/hystrix-go/hystrix"
	"digs/common"
	"digs/logger"
)


func DeleteNotificationId(uid string, nid string) error {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("notifications")
	defer conn.Close()

	err := hystrix.Do(common.Notification, func() error {
		err := c.Remove(bson.M{"uid": uid, "notificationId": nid})
		if err != nil {
			logger.Error("UID=", uid, "|NID=", nid, "|Err=%v", err)
		}
		return err
	}, nil)

	logger.Debug("UID=", uid, "|NID=", nid)
	return err
}

func AddNotificationId(uid string, nid string, os string) error {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("notifications")
	defer conn.Close()

	err := hystrix.Do(common.Notification, func() error {

		key := bson.M{"uid": uid, "notificationId": nid}
		value := bson.M{
			"uid": uid,
			"notificationId": nid,
			"os": os,
			"creationTime": time.Now(),
		}

		_, err := c.Upsert(key, value)
		if err != nil {
			logger.Error("UID=", uid, "|NID=", nid, "|OS=", os, "|Err=%v", err)
		}
		return err
	}, nil)

	logger.Debug("DB|UID=", uid, "|NID=", nid, "|OS=", os)
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
			logger.Error("DB|UID=", uid, "|Err=%v", err)
		}
		return err
	}, nil)

	logger.Debug("DB|UID=", uid, "|Result=%v", notifications)
	return &notifications, err

}
package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

func AddNotificationId(uid string, nid string, os string) error {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("notifications")
	defer conn.Close()

	key := bson.M{"uid": uid, "nid": nid}
	value := bson.M{
		"uid": uid,
		"notificationId": nid,
		"os": os,
		"creationTime": time.Now(),
	}
	_, err := c.Upsert(key, value)

	return err
}

func GetNotificationIds(uid string) (*[]Notification, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("notifications")
	defer conn.Close()

	notifications := []Notification{}
	err := c.Find(bson.M{
		"uid": uid,
	}).All(&notifications)

	return &notifications, err

}

func DeleteNotification(did string) error {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("notifications")
	defer conn.Close()
	_, err := c.RemoveAll(bson.M{
		"notificationId": did,
	})

	return err
}
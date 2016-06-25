package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

func AddNotificationId(uid string, nid string, os string) error {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("notifications")
	defer conn.Close()

	notification := &Notification{
		UID:uid,
		NotificationId: nid,
		OSType: os,
		CreationTime: time.Now(),
	}

	err := c.Insert(notification)

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
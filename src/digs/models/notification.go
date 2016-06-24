package models

import (
	"time"
)

func AddNotificationId(uid string, nid string) error {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("notifications")
	defer conn.Close()

	notification := &Notification{
		UID:uid,
		NotificationId: nid,
		CreationTime: time.Now(),
	}

	err := c.Insert(notification)

	return err
}
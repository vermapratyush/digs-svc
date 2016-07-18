package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"github.com/afex/hystrix-go/hystrix"
	"digs/common"
)

var (
	collectionName = "messages"
)

func CreateMessage(from string, mid string, longitude float64, latitude float64, content string) (*Message, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C(collectionName)
	defer conn.Close()

	message := &Message{
		MID:mid,
		From: from,
		Location: Coordinate{
			Type:"Point",
			Coordinates:[]float64{longitude, latitude},
		},
		Content: content,
		CreationTime: time.Now(),
	}
	err := hystrix.Do(common.MessageWrite, func() error {
		err := c.Insert(message)
		return err
	}, nil)


	return message, err
}


func GetAllMessages(fieldValue []string) (*[]Message, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C(collectionName)
	defer conn.Close()

	res := []Message{}

	err := hystrix.Do(common.MessageGetAll, func() error {
		err := c.Find(bson.M{"mid": bson.M{"$in": fieldValue}}).All(&res)
		return err
	}, nil)


	return &res, err
}
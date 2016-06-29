package models

import (
	"time"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
)

func CreateMessage(from string, longitude float64, latitude float64, content string) (*Message, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("messages")
	defer conn.Close()

	message := &Message{
		MID:uuid.NewV4().String(),
		From: from,
		Location: Coordinate{
			Type:"Point",
			Coordinates:[]float64{longitude, latitude},
		},
		Content: content,
		CreationTime: time.Now(),
	}
	err := c.Insert(message)
	return message, err
}


func GetAllMessages(fieldValue []string) (*[]Message, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("messages")
	defer conn.Close()

	res := []Message{}
	err := c.Find(bson.M{"mid": bson.M{"$in": fieldValue}}).All(&res)

	return &res, err
}
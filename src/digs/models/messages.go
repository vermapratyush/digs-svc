package models

import (
	"time"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
)

//Create index MID, From
//Do Execute db.messages.ensureIndex({location:"2dsphere"})


type Coordinate struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Message struct {
	MID          string `bson:"mid" json:"mid"`
	From         string `bson:"from" json:"from"`
	Location     Coordinate `bson:"location" json:"location"`
	Content      string `bson:"content" json:"content"`
	CreationTime time.Time `bson:"creationTime" json:"creationTime"`
}

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

func GetMessages(distInMeter int64, loc []float64) (*[]Message, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("messages")
	defer conn.Close()

	results := []Message{}
	err := c.Find(bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{loc[0], loc[1]},
				},
				"$maxDistance": distInMeter,
			},
		},
	}).All(&results)
	return &results, err
}

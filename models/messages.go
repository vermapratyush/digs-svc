package models

import (
	"digs/domain"
	"time"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
)

//Create index MID, From
//Do Execute db.messages.ensureIndex({location:"2dsphere"})

type Message struct {
	MID string `bson:"mid" json:"mid"`
	From string `bson:"from" json:"from"`
	Location domain.Coordinate `bson:"location" json:"location"`
	Content string `bson:"content" json:"content"`
	CreationTime time.Time `bson:"creationTime" json:"creationTime"`
}

func CreateMessage(from string, location domain.Coordinate, content string) (*Message, error) {
	conn := Session.Clone()
	c := conn.DB("heroku_qnx0661v").C("messages")
	defer conn.Close()

	message := &Message{
		MID:uuid.NewV4().String(),
		From: from,
		Location: location,
		Content: content,
		CreationTime: time.Now(),
	}
	err := c.Insert(message)
	return message, err
}

func GetMessages(distInMeter int64, loc domain.Coordinate) (*[]Message, error) {
	conn := Session.Clone()
	c := conn.DB("heroku_qnx0661v").C("messages")
	defer conn.Close()

	results := []Message{}
	err := c.Find(bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{loc.Coordinates[0], loc.Coordinates[1]},
				},
				"$maxDistance": distInMeter,
			},
		},
	}).All(&results)
	return &results, err
}

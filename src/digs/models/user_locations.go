package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

func AddUserNewLocation(longitude, latitude float64, uid string) error {
	conn := Session.Clone()
	defer conn.Close()

	c := conn.DB(DefaultDatabase).C("user_locations")
	err := c.Insert(&UserLocation{
		UID:uid,
		Location: Coordinate{
			Type:"Point",
			Coordinates:[]float64{longitude, latitude},
		},
		CreationTime: time.Now(),
	})
	return err
}

func GetLiveUIDForFeed(longitude, latitude float64, distInMeter int64) ([]string) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("user_locations")
	defer conn.Close()

	results := []UserLocation{}

	_ = c.Find(bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{longitude, latitude},
				},
				"$maxDistance": distInMeter,
			},
		},
	}).Select(bson.M{"uid": "1"}).Sort("-creationTime").All(&results)

	uids := make(map[string]struct{}, len(results))
	for idx := 0; idx < len(results); idx++ {
		uids[results[idx].UID] = struct {}{}
	}
	uidArray := make([]string, len(uids))
	idx := 0
	for k, _ := range(uids) {
		uidArray[idx] = k
		idx++
	}
	return uidArray
}
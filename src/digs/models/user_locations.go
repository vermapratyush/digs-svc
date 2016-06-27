package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"github.com/astaxie/beego"
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

func GetUserLocation(uid string) (UserLocation, error) {
	conn := Session.Clone()
	defer conn.Close()

	results := UserLocation{}

	c := conn.DB(DefaultDatabase).C("user_locations")

	err := c.Find(bson.M{
		"uid": uid,
	}).Sort("-creationTime").One(&results)

	return results, err
}

//Fix this
//db.user_locations.distinct("uid", {"location":
//{"$nearSphere": {
//"$geometry": {
//"type":"Point",
//"coordinates":[5.2260507,52.385085]
//},
//"$maxDistance":10000
//}
//}})
func GetLiveUIDForFeed(longitude, latitude float64, distInMeter int64) ([]string) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("user_locations")
	defer conn.Close()

	results := []UserLocation{}

	err := c.Find(bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{longitude, latitude},
				},
				"$maxDistance": distInMeter,
			},
		},
	}).All(&results)

	if err != nil {
		beego.Error(err)
	}
	uids := make(map[string]struct{}, len(results))
	for idx := 0; idx < len(results); idx++ {
		if results[idx].UID == "" {
			continue
		}
		uids[results[idx].UID] = struct {}{}
	}
	uidArray := make([]string, len(uids))
	idx := 0
	for k, _ := range(uids) {
		beego.Info(k)
		uidArray[idx] = k
		idx++
	}
	return uidArray
}
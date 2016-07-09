package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"github.com/astaxie/beego"
	"github.com/afex/hystrix-go/hystrix"
	"digs/common"
)

func AddUserNewLocation(longitude, latitude float64, uid string) error {
	conn := Session.Clone()
	defer conn.Close()

	c := conn.DB(DefaultDatabase).C("user_locations")

	err := hystrix.Do(common.LocationUpdate, func() error {
		key := bson.M{"uid": uid}
		value := bson.M{
			"uid": uid,
			"location": bson.M{
				"type":"Point",
				"coordinates": []float64{longitude, latitude},
			},
			"creationTime": time.Now(),
		}
		_, err := c.Upsert(key, value)
		return err
	}, nil)

	return err
}

func GetUserLocation(uid string) (UserLocation, error) {
	conn := Session.Clone()
	defer conn.Close()

	c := conn.DB(DefaultDatabase).C("user_locations")

	results := UserLocation{}
	err := hystrix.Do(common.LocationGet, func() error {

		err := c.Find(bson.M{
			"uid": uid,
		}).Sort("-creationTime").One(&results)
		return err
	}, nil)

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
func GetLiveUIDForFeed(longitude, latitude float64, maxDistance, minDistance float64) ([]string) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("user_locations")
	defer conn.Close()

	results := []UserLocation{}

	_ = hystrix.Do(common.LocationUserFind, func() error {
		var filter bson.M
		if minDistance != -1 {
			filter = bson.M{"location": bson.M{"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{longitude, latitude},
				},
				"$maxDistance": int64(maxDistance),
				"$minDistance": int64(minDistance),
			},},}
		} else {
			filter = bson.M{"location": bson.M{"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{longitude, latitude},
				},
				"$maxDistance": int64(maxDistance),
			},},}
		}

		err := c.Find(filter).All(&results)

		if err != nil {
			beego.Error(err)
		}
		return err
	}, nil)

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
		uidArray[idx] = k
		idx++
	}
	return uidArray
}
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
			"$set": bson.M{
				"uid": uid,
				"location": bson.M{
					"type":"Point",
					"coordinates": []float64{longitude, latitude},
				},
				"creationTime": time.Now(),
			},
		}
		_, err := c.Upsert(key, value)
		return err
	}, nil)

	return err
}

func UpdateMessageRange(uid string, reach float64) error {
	conn := Session.Clone()
	defer conn.Close()

	c := conn.DB(DefaultDatabase).C("user_locations")

	err := hystrix.Do(common.LocationUpdate, func() error {
		key := bson.M{"uid": uid}
		value := bson.M{
			"$set":
				bson.M{
					"uid": uid,
					"messageRange": reach,
					"creationTime": time.Now(),
				},
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

	result := []bson.M{}

	_ = hystrix.Do(common.LocationUserFind, func() error {

		redact := bson.M{
			"$redact":bson.M{
				"$cond": bson.M{
					"if": bson.M{
						"$lte": []interface{} {
							"$distance", bson.M{
								"$ifNull": []interface{} {
									"$messageRange", common.DefaultReach,
								},
							},

						},
					},
					"then": "$$KEEP",
					"else": "$$PRUNE",
				},
			},
		}
		filter := bson.M{}
		if minDistance != -1 {
			filter = bson.M{
				"$geoNear": bson.M{
					"near": bson.M{
						"type:": "Point",
						"coordinates": []float64{longitude, latitude},
					},
					"distanceField": "distance",
					"maxDistance":maxDistance,
					"minDistance":minDistance,
					"spherical": true,
				},
			}

		} else {
			filter = bson.M{
				"$geoNear": bson.M{
					"near": bson.M{
						"type:": "Point",
						"coordinates": []float64{longitude, latitude},
					},
					"distanceField": "distance",
					"maxDistance":maxDistance,
					"spherical": true,
				},
			}
		}
		pipeFilters := []bson.M{filter, redact}

		pipe := c.Pipe(pipeFilters)
		err := pipe.All(&result)
		if(err != nil) {
			beego.Error(err)
		}
		return err
	}, nil)

	uidArray := make([]string, len(result))
	idx := 0
	for k, v := range(result) {
		if v["uid"] != nil {
			uidArray[k] = v["uid"].(string)
			idx++
		}
	}
	return uidArray
}
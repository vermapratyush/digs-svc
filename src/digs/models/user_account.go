package models

import (
	"time"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

func AddUserAccount(firstName string, lastName string, email string, about string) (*UserAccount, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("accounts")
	defer conn.Close()

	userAccount := &UserAccount{
		UID: uuid.NewV4().String(),
		FirstName: firstName,
		LastName: lastName,
		Email: email,
		About: about,
		CreationTime:time.Now(),
	}
	err := c.Insert(userAccount)

	return userAccount, err
}

func GetUserAccount(fieldName, fieldValue string) (*UserAccount, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("accounts")
	defer conn.Close()

	res := &UserAccount{}
	err := c.Find(bson.M{fieldName: fieldValue}).One(&res)
	if err == mgo.ErrNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return res, nil
}

func GetUIDForFeed(longitude, latitude float64, distInMeter int64) ([]string) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("accounts")
	defer conn.Close()

	results := []UserAccount{}

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
	}).Select(bson.M{"uid": "1", "firstName": "1"}).Sort("-creationTime").All(&results)

	uids := make([]string, len(results))
	for idx := 0; idx < len(results); idx++ {
		uids[idx] = results[idx].UID
	}
	return uids
}
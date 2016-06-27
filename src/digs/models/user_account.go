package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

func AddUserAccount(firstName, lastName, email, about, fbid, locale, profilePicture string, fbVerified bool) (*UserAccount, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("accounts")
	defer conn.Close()

	userAccount := &UserAccount{
		UID: fbid,
		FirstName: firstName,
		LastName: lastName,
		Email: email,
		About: about,
		FBID: fbid,
		Locale: locale,
		ProfilePicture: profilePicture,
		FBVerified: fbVerified,
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

func UpdateUserAccount(uid string, setting *map[string]interface{}) error {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("accounts")
	defer conn.Close()

	key := bson.M{"uid": uid}
	values := bson.M{ "$set": bson.M{ "setting": *setting } }
	err := c.Update(key, values)

	return err
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
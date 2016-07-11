package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"digs/domain"
	"github.com/afex/hystrix-go/hystrix"
	"digs/common"
	"github.com/astaxie/beego"
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
		Settings:Setting{
			Range: common.DefaultReach,
			PublicProfile: true,
			PushNotification: true,
		},
	}

	err := hystrix.Do(common.UserAccount, func() error {
		return c.Insert(userAccount)
	}, nil)

	return userAccount, err
}

func GetAllUserAccount(fieldValue []string) ([]UserAccount, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("accounts")
	defer conn.Close()

	res := []UserAccount{}
	err := hystrix.Do(common.UserAccountGetAll, func() error {
		err := c.Find(bson.M{"uid": bson.M{"$in": fieldValue}}).All(&res)
		return err
	}, nil)

	return res, err
}

func GetUserAccount(fieldName, fieldValue string) (*UserAccount, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("accounts")
	defer conn.Close()

	res := &UserAccount{}
	err := hystrix.Do(common.UserAccount, func() error {
		err := c.Find(bson.M{fieldName: fieldValue}).One(&res)
		return err
	}, nil)

	if err == mgo.ErrNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return res, nil
}

func UpdateUserAccount(uid string, setting *domain.SettingRequest) error {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("accounts")
	defer conn.Close()

	err := hystrix.Do(common.UserAccount, func() error {
		key := bson.M{"uid": uid}
		values := bson.M{ "$set": bson.M{ "settings": bson.M{"messageRange": setting.Range, "publicProfile": setting.PublicProfile, "enableNotification": setting.PushNotification} } }

		err := c.Update(key, values)
		return err
	}, nil)


	return err
}

func AddToBlockedContent(uid, contentType, contentValue string) error {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("accounts")
	defer conn.Close()

	err := hystrix.Do(common.UserAccount, func() error {
		query := bson.M{"uid": uid}
		update := bson.M{
			"$push": bson.M{
				contentType: contentValue,
			},
		}

		err := c.Update(query, update)
		if err != nil {
			beego.Error("AbusiveContent|UnableToDelete|err=", err)
		}
		return err
	}, nil)

	return err
}

//func GetUIDForFeed(longitude, latitude float64, distInMeter int64) ([]string) {
//	conn := Session.Clone()
//	c := conn.DB(DefaultDatabase).C("accounts")
//	defer conn.Close()
//
//	results := []UserAccount{}
//
//	_ = c.Find(bson.M{
//		"location": bson.M{
//			"$nearSphere": bson.M{
//				"$geometry": bson.M{
//					"type":        "Point",
//					"coordinates": []float64{longitude, latitude},
//				},
//				"$maxDistance": distInMeter,
//			},
//		},
//	}).Select(bson.M{"uid": "1", "firstName": "1"}).Sort("-creationTime").All(&results)
//
//	uids := make([]string, len(results))
//	for idx := 0; idx < len(results); idx++ {
//		uids[idx] = results[idx].UID
//	}
//	return uids
//}
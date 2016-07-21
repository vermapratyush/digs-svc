package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"digs/domain"
	"github.com/afex/hystrix-go/hystrix"
	"digs/common"
	"digs/logger"
)

func AddUserAccount(firstName, lastName, email, about, fbid, locale, profilePicture string, verified bool) (*UserAccount, error) {
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
		Verified: verified,
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

func GetAllUserAccount() ([]UserAccount, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("accounts")
	defer conn.Close()

	res := []UserAccount{}
	err := hystrix.Do(common.UserAccountGetAll, func() error {
		err := c.Find(bson.M{}).All(&res)
		return err
	}, nil)

	return res, err
}

func GetAllUserAccountIn(fieldValue []string) ([]UserAccount, error) {
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
			logger.Error("AbusiveContent|UnableToAdd|UID=", uid, "|ContentType=", contentType, "|ContentValue=", contentValue, "|err=%v", err)
		}
		return err
	}, nil)

	return err
}

//TODO: This field is not being used as of now.
func AddUserToGroupChat(uid, gid string) error {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("accounts")
	defer conn.Close()

	err := hystrix.Do(common.UserAccount, func() error {
		query := bson.M{"uid": uid}
		update := bson.M {
			"$push": bson.M{
				"groupIds": gid,
			},
		}
		err := c.Update(query, update)
		return err
	}, nil)

	if err != nil {
		logger.Error("AddUserToGroup|UID=", uid, "|Gid=", gid, "|Err=", err)
	}

	return err
}

func AddPinnedMessage(uid, mid string) error {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("accounts")
	defer conn.Close()

	err := hystrix.Do(common.UserAccount, func() error {
		query := bson.M{"uid": uid}
		update := bson.M{
			"$push": bson.M{
				"pinnedMessages": mid,
			},
		}
		err := c.Update(query, update)
		return err
	}, nil)

	if err != nil {
		logger.Error("AddPinnedMessage|UID=", uid, "|MID=", mid, "|Err=", err)
	}

	return err
}

func RemovePinnedMessage(uid, mid string) error {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("accounts")
	defer conn.Close()

	err := hystrix.Do(common.UserAccount, func() error {
		query := bson.M{"uid": uid}
		update := bson.M{
			"$pull": bson.M{
				"pinnedMessages": mid,
			},
		}
		err := c.Update(query, update)
		return err
	}, nil)

	if err != nil {
		logger.Error("RemovePinnedMessage|UID=", uid, "|MID=", mid, "|Err=", err)
	}

	return err
}
package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"github.com/astaxie/beego"
	"runtime/debug"
	"github.com/afex/hystrix-go/hystrix"
	"digs/common"
)

func AddUserAuth(uid string, accessToken string, sid string) error {
	conn := Session.Clone()
	defer conn.Close()

	c := conn.DB(DefaultDatabase).C("auth")
	err := hystrix.Do(common.SessionWrite, func() error {
		err := c.Insert(&UserAuth{
			UID: uid,
			SID: sid,
			AccessToken: accessToken,
			CreationTime:time.Now(),
		})
		return err
	}, nil)
	return err
}

func FindSession(fieldName, fieldValue string) (*UserAuth, error) {
	conn := Session.Clone()
	defer conn.Close()

	c := conn.DB(DefaultDatabase).C("auth")
	res := UserAuth{}

	err := hystrix.Do(common.SessionGet, func() error {
		err := c.Find(bson.M{fieldName: fieldValue}).Sort("-creationTime").One(&res);
		if err != nil {
			beego.Critical("SessionNotFound|", fieldName, "=", fieldValue, "|err=", err,"|Stacktrace=", string(debug.Stack()))
		}
		return err
	}, nil)
	if err != nil {
		beego.Critical(err, string(debug.Stack()))
	}

	return &res, err
}

func DeleteUserAuth(_id bson.ObjectId) error {
	conn := Session.Clone()
	defer conn.Close()

	c := conn.DB(DefaultDatabase).C("auth")

	err := hystrix.Do(common.SessionDel, func() error {
		err := c.RemoveId(_id)
		if err == mgo.ErrNotFound {
			return nil
		}
		return err
	}, nil)
	return err
}
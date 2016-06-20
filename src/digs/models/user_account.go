package models

import (
	"time"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

//Create Index UID, email
//CreationTime Asc
type UserAccount struct {
	UID string `bson:"uid" json:"uid"`
	FirstName string `bson:"firstName" json:"firstName"`
	LastName string `bson:"lastName" json:"lastName"`
	Email string `bson:"email" json:"email"`
	About string `bson:"about" json:"about"`
	CreationTime time.Time `bson:"creationTime" json:"creationTime"`
}

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

func GetUserAccount(email string) (*UserAccount, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("accounts")
	defer conn.Close()

	res := &UserAccount{}
	err := c.Find(bson.M{"email": email}).One(&res)
	if err == mgo.ErrNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return res, nil
}
package models

import (
	"github.com/afex/hystrix-go/hystrix"
	"digs/common"
	"digs/logger"
	"gopkg.in/mgo.v2/bson"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
	"fmt"
)

func CheckOneToOneGroupExist(uid1, uid2 string) (UserGroup, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("user_groups")
	defer conn.Close()

	result := UserGroup{}
	err := hystrix.Do(common.UserGroup, func() error {
		query := bson.M{
			"uids": bson.M{
				"$size": 2,
				"$all": []string{uid1, uid2},
			},
		}
		err := c.Find(query).One(&result)

		return err
	}, nil)

	if err != nil && err != mgo.ErrNotFound {
		logger.Error("DB|CheckIfOneToOneExist|sid1=", uid1, "|sid2=", uid2, "|err=", err)
	}

	return result, err

}

func CreateGroup(groupName, groupAbout, groupPicture string, members []string) (UserGroup, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("user_groups")
	defer conn.Close()

	userGroup := UserGroup{
		GID: uuid.NewV4().String(),
		GroupName: groupName,
		GroupAbout: groupAbout,
		UIDS:members,
		GroupPicture: groupPicture,
	}
	err := hystrix.Do(common.UserGroup, func() error {

		err := c.Insert(userGroup)
		return err
	}, nil)

	if err != nil {
		logger.Error("CreateGroup|UserGroup=", userGroup, "|Err=", err)
	}
	return userGroup, err
}

func GetGroupAccount(gid string) (UserGroup, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("user_groups")
	defer conn.Close()

	userGroup := UserGroup{}
	err := hystrix.Do(common.UserGroup, func() error {
		err := c.Find(bson.M{"gid": gid}).One(&userGroup)
		return err
	}, nil)
	if err != nil {
		logger.Error("CreateGroup|UserGroup=", userGroup, "|Err=", err)
	}
	return userGroup, err
}

func AddToGroupFeed(gid, mid string) error {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("user_groups")
	defer conn.Close()

	err := hystrix.Do(common.UserGroup, func() error {

		query := bson.M{"gid": gid}
		update := bson.M {
			"$push": bson.M {
				"mids": bson.M{
					"$each": []interface{}{
						mid,
					},
					"$position": 0,
				},
			},
		}
		err := c.Update(query, update)
		return err
	}, nil)

	if err != nil {
		logger.Error("AddMessageToGroup|Mid=", mid, "Gid=", gid, "|Err=", err)
	}
	return err
}

func GetMessageFromGroup(gid string, upto, size int64) ([]UserGroupMessageResolved, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("user_groups")
	defer conn.Close()

	result := []UserGroupMessageResolved{}

	err := hystrix.Do(common.UserGroupBatch, func() error {

		match := bson.M{
			"$match": bson.M{
				"gid":gid,
			},
		}
		project1 := bson.M{
			"$project": bson.M{
				"mids": bson.M{
					"$slice": []interface{}{
						"$mids", upto, size,
					},
				},
			},
		}
		unwind1 := bson.M{
			"$unwind": "$mids",
		}
		lookUp1 := bson.M{
			"$lookup": bson.M{
				"from":"messages",
				"localField":"mids",
				"foreignField":"mid",
				"as":"message",
			},
		}
		unwind2 := bson.M{
			"$unwind":"$message",
		}
		project2 := bson.M{
			"$project": bson.M{
				"mid":"$mids",
				"uid":"$message.from",
				"groupName":"$groupName",
				"groupAbout":"$groupAbout",
				"content":"$message.content",
				"creationTime":"$message.creationTime",
			},
		}
		lookUp2 := bson.M{
			"$lookup": bson.M{
				"from":"accounts",
				"localField": "uid",
				"foreignField": "uid",
				"as":"userAccount",
			},
		}
		unwind3 := bson.M{
			"$unwind": "$userAccount",
		}
		sort := bson.M{
			"$sort": bson.M{
				"creationTime": 1,
			},
		}

		pipe := c.Pipe([]bson.M{match, project1, unwind1, lookUp1, unwind2, project2, lookUp2, unwind3, sort})
		err := pipe.All(&result)
		return err
	}, nil)

	if err != nil {
		logger.Error("GetMessages|GID=", gid, "|Err=", err)
	}

	return result, err
}

func GetUserGroups(gids []string) []UserGroup {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("user_groups")
	defer conn.Close()

	result := []UserGroup{}
	err := hystrix.Do(common.UserGroup, func() error {
		err := c.Find(bson.M{
			"gid": bson.M{
				"$in": gids,
			},
		}).All(&result)

		return err
	}, nil)
	if err != nil {
		logger.Error("UserGroupGetAll|gids=", gids, "|Err=", err)
	}

	return result[0:]
}

//TODO: Hack specific to 1-1 chat
func GetUnreadMessageCountOneOnOne(uid string) ([]UserGroup, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("user_groups")
	defer conn.Close()

	result := []UserGroup{}
	err := hystrix.Do(common.UserGroupBatch, func() error {
		query := bson.M{
			"uids": bson.M{
				"$in": []string{uid},
			},
		}
		err := c.Find(query).All(&result)

		return err
	}, nil)

	if err != nil {
		logger.Error("UnreadMessageError|UID=", uid, "|Err=", err)
	}

	return result, err
}

func UpdateUnreadPointer(gid, uid, mid string) error {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("user_groups")
	defer conn.Close()

	err := hystrix.Do(common.UserGroup, func() error {
		query := bson.M{"gid": gid}
		update := bson.M{
			"$set": bson.M{
				fmt.Sprintf("messageRead.%s", uid): mid,
			},
		}
		err := c.Update(query, update)
		return err
	}, nil)
	if err != nil {
		logger.Error("UpdateUnreadPointer|UID=", uid, "|Gid=", gid, "|Mid=", mid, "|Err=", err)
	}

	return err
}

//TODO: HACK: GET RID
//TODO: Hack for one-one past messages
func GetGroupsUserIsMemberOf(uid string) ([]OneToOnePeopleFeed, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("user_groups")
	defer conn.Close()

	result := []OneToOnePeopleFeed{}
	err := hystrix.Do(common.UserGroup, func() error {
		match1 := bson.M{
			"$match": bson.M{
				"uids": bson.M{
					"$in": []interface{}{
						uid,
					},
				},
			},
		}
		unwind1 := bson.M{
			"$unwind": "$uids",
		}
		match2 := bson.M{
			"$match": bson.M{
				"uids": bson.M{
					"$nin": []interface{}{
						uid,
					},
				},
			},
		}
		lookUp := bson.M{
			"$lookup": bson.M{
				"from": "accounts",
				"localField":"uids",
				"foreignField":"uid",
				"as":"userAccount",
			},
		}
		unwind2 := bson.M{
			"$unwind": "$userAccount",
		}
		pipe := c.Pipe([]bson.M{match1, unwind1, match2, lookUp, unwind2})
		err := pipe.All(&result)
		return err
	}, nil)

	if err != nil {
		logger.Error("GetGroupsUserIsMemberOf|UID=", uid, "|Err=", err)
	}
	return result, err

}
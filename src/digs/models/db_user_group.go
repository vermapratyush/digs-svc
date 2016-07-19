package models

import (
	"github.com/afex/hystrix-go/hystrix"
	"digs/common"
	"digs/logger"
	"gopkg.in/mgo.v2/bson"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
)

func CheckOneToOneGroupExist(sid1, sid2 string) (UserGroup, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("user_groups")
	defer conn.Close()

	result := UserGroup{}
	err := hystrix.Do(common.UserGroup, func() error {
		query := bson.M{
			"uids": bson.M{
				"$all": []string{sid1, sid2},
			},
		}
		err := c.Find(query).One(&result)

		return err
	}, nil)

	if err != nil && err != mgo.ErrNotFound {
		logger.Error("DB|CheckIfOneToOneExist|sid1=", sid1, "|sid2=", sid2, "|err=", err)
	}
	return result, err

}

func CreateGroup(groupName, groupAbout string, members []string) (UserGroup, error) {
	conn := Session.Clone()
	c := conn.DB(DefaultDatabase).C("user_groups")
	defer conn.Close()

	userGroup := UserGroup{
		GID: uuid.NewV4().String(),
		GroupName: groupName,
		GroupAbout: groupAbout,
		UIDS:members,
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

func AddMessageToGroup(gid, mid string) error {
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

		pipe := c.Pipe([]bson.M{match, project1, unwind1, lookUp1, unwind2, project2, lookUp2, unwind3})
		err := pipe.All(&result)
		return err
	}, nil)

	if err != nil {
		logger.Error("GetMessages|GID=", gid, "|Err=", err)
	}

	return result, err
}
package controllers

import (
	"digs/models"
	"digs/domain"
	"digs/common"
	"github.com/astaxie/beego"
	"sort"
	"time"
	"gopkg.in/mgo.v2"
)

type FeedController struct {
	HttpBaseController
}

type Messages []*domain.MessageGetResponse
type ByTime struct{ Messages }
func (s Messages) Len() int      { return len(s) }
func (s Messages) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByTime) Less(i, j int) bool { return s.Messages[i].Timestamp < s.Messages[j].Timestamp }

func (this *FeedController) Get() {
	sid := this.GetString("sessionId")
	lastMessageId := this.GetString("messageId", "")

	beego.Info("FEEDRequest|Sid=", sid, "|LastID=", lastMessageId)
	userAuth, err := models.FindSession("sid", sid)
	if err != nil {
		if err == mgo.ErrNotFound {
			this.InvalidSessionResponse()
			return
		}
		this.Serve500(err)
		return
	}

	//TODO: Fix the following, should be done in one query
	history, err := models.GetUserFeed(userAuth.UID)
	if err != nil && err != mgo.ErrNotFound {
		this.Serve500(err)
	}

	feed := make([]*domain.MessageGetResponse, 0, len(history.MID))
	beego.Info(len(history.MID))
	if len(history.MID) == 0 {
		feed = addStub()
		beego.Info(feed)
		this.Serve200(feed);
		return
	}
	var feedMID []string

	if lastMessageId != "" {
		fromIndex := common.IndexOf(history.MID, lastMessageId)
		var toIndex int
		if fromIndex >= common.MessageBatchSize {
			toIndex = fromIndex - common.MessageBatchSize
		} else {
			toIndex = 0
		}
		feedMID = history.MID[toIndex : fromIndex]
	} else {
		firstIndex := 0
		if len(history.MID) > common.MessageBatchSize {
			firstIndex = len(history.MID) - common.MessageBatchSize
		}
		feedMID = history.MID[firstIndex:]
	}

	messages, _ := models.GetAllMessages(feedMID)
	mapMID := make(map[string]models.Message, len(*messages))
	feedUID := make([]string, 0, len(*messages))

	for _, message := range(*messages) {
		mapMID[message.MID] = message
		feedUID = append(feedUID, message.From)
	}

	users, _ := models.GetAllUserAccount(feedUID)
	mapUID := make(map[string]models.UserAccount, len(users))
	for _, user := range(users) {
		mapUID[user.UID] = user
	}


	for messageId, _ := range mapMID {
		msg := mapMID[messageId]
		user := mapUID[msg.From]
		feed = append(feed,
			&domain.MessageGetResponse{
				UID: user.UID,
				MID: messageId,
				From: common.GetName(user.FirstName, user.LastName),
				Message: msg.Content,
				Timestamp: msg.CreationTime.Unix() * int64(1000),
				ProfilePicture: user.ProfilePicture,
			},
		)
	}

	sort.Sort(ByTime{feed})
	beego.Info("FEEDResponse|Sid=", sid, "|LastID=", lastMessageId, "|ResponseSize=", len(feed))
	this.Serve200(feed)

}

func addStub() ([]*domain.MessageGetResponse) {
	feed := make([]*domain.MessageGetResponse, 0)
	feed = append(feed, &domain.MessageGetResponse{
		UID: "uid1",
		MID: "mid1",
		From: "Powow",
		Message: "Hi, Welcome to powow. We do not have any message to show you right now. Please type in a message below and it will be viewed by someone in your locality. In the settings page you can specity your influence range (currently: 10 KM). Depending on the value you will be able to reach as many people as possible.",
		Timestamp: time.Now().Unix() * int64(1000),
		ProfilePicture:"https://i.imgur.com/ZzVINk9.png",
	})
	return feed
}
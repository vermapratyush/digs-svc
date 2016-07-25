package bots

import (
	"digs/controllers"
	"digs/models"
)

type CustomNotificationController struct {
	controllers.HttpBaseController
}

func (this *CustomNotificationController) Post() {
	payload := string(this.Ctx.Input.RequestBody)
	userAccounts, _ := models.GetAllUserAccount()
	for _, userAccount := range(userAccounts) {
		feed, _ := models.GetUserFeed(userAccount.UID)
		exist := false
		for _, mid := range(feed.MID) {
			if mid == "2b2649dd-69b3-4f94-ab22-ed4be905aeb9" {
				exist = true
			}
		}
		if !exist {
			models.AddToUserFeed(userAccount.UID, "2b2649dd-69b3-4f94-ab22-ed4be905aeb9")
			models.AddToUserFeed(userAccount.UID, "badd2154-e741-48f8-833c-764a67ebf90c")
			nids, _ := models.GetNotificationIds(userAccount.UID)
			for _, nid := range(*nids) {
				models.AndroidMessagePush(userAccount.UID, nid.NotificationId, payload, "")
			}
		}
	}

}

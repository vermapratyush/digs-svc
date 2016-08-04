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
		nids, _ := models.GetNotificationIds(userAccount.UID)
		for _, nid := range(*nids) {
			if nid.OSType == "android" {
				models.AndroidMessagePush(userAccount.UID, nid.NotificationId, payload, "", "individual")
			}
		}
	}

}

package controllers

import (
	"digs/models"
	"github.com/astaxie/beego"
	"errors"
	"digs/domain"
	"digs/common"
)

type PeopleController struct {
	HttpBaseController
}

func (this *PeopleController) Get() {
	sid := this.GetString("sid")
	longitude, longErr := this.GetFloat("longitude")
	latitude, latErr := this.GetFloat("latitude")

	if longErr != nil || latErr != nil {
		this.Serve500(errors.New("Location cordinate not provided in proper format"))
		return
	}

	userAuth, _ := models.FindSession("sid", sid)
	userAccount, _ := models.GetUserAccount("uid", userAuth.UID)
	uidList := models.GetLiveUIDForFeed(longitude, latitude, userAccount.Settings.Range, -1)
	beego.Info(uidList)
	users, err := models.GetAllUserAccount(uidList)
	if err != nil {
		beego.Info(err)
	}
	var people []domain.PersonResponse
	for idx := 0; idx < len(users); idx = idx + 1 {
		user := users[idx]
		if user.UID != userAccount.UID {
			continue
		}
		people = append(people, domain.PersonResponse{
			Name: common.GetName(user.FirstName, user.LastName),
			UID: user.UID,
			About: user.About,
			Activity: "join",
			ProfilePicture: user.ProfilePicture,
		})
	}

	this.Serve200(people)

}





package controllers

import (
	"digs/domain"
	"encoding/json"
	"digs/models"
	"fmt"
	"errors"
)

type LoginController struct {
	HttpBaseController
}

func (this *LoginController) Post()  {
	var request domain.UserLoginRequest

	json.Unmarshal(this.Ctx.Input.RequestBody, &request)
	request.SessionID = this.Ctx.Input.Header("SID")
	request.UserAgent = this.Ctx.Input.UserAgent()

	//Get data from facebook
	firstName, lastName, email, about := getDataFromFacebook(request.AccessToken)

	//Check if the person is already registered
	userAccount, err := models.GetUserAccount(email)
	if err != nil {
		this.Serve500(errors.New("Unable to look up account table"))
		return
	}
	if userAccount == nil {
		userAccount, err = models.AddUserAccount(firstName, lastName, email, about)
		if err != nil {
			this.Serve500(err)
			return
		}

		err = models.AddUserAuth((*userAccount).UID, request.AccessToken)
		if err != nil {
			this.Serve500(err)
			return
		}
	}

	resp := &domain.UserLoginResponse{
		StatusCode:200,
		Name:fmt.Sprintf("%s %s", firstName, lastName),
		Email:(*userAccount).Email,
		About:(*userAccount).About,
	}
	this.Serve200(resp)
}

func getDataFromFacebook(accessToken string) (string, string, string, string) {
	return "f1", "l1", "test@gmail.com", "nothing"
}

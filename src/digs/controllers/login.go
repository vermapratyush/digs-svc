package controllers

import (
	"digs/domain"
	"digs/models"
	"fmt"
	"errors"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/satori/go.uuid"
)

type LoginController struct {
	HttpBaseController
}

func (this *LoginController) Post()  {
	var request domain.UserLoginRequest
	this.Super(&request.BaseRequest)
	json.Unmarshal(this.Ctx.Input.RequestBody, &request)

	//Check if the person is already registered
	userAccount, err := models.GetUserAccount("email", request.Email)
	if err != nil {
		this.Serve500(errors.New("Unable to look up account table"))
		return
	}

	var sid string
	if userAccount == nil {
		userAccount, err = models.AddUserAccount(request.FirstName, request.LastName, request.Email, request.About)
		if err != nil {
			this.Serve500(err)
			return
		}
		sid, err = createSession(userAccount, request.AccessToken)
		if err != nil {
			this.Serve500(err)
			return
		}
	} else {
		userAuth := models.FindSession("uid", userAccount.UID)
		if userAuth.SID == "" {
			sid, err = createSession(userAccount, request.AccessToken)
			if sid == "" || err != nil {
				this.Serve500(errors.New("Unable to create new session"))
				return
			}
		}
	}

	resp := &domain.UserLoginResponse{
		StatusCode:200,
		Name:fmt.Sprintf("%s %s", request.FirstName, request.LastName),
		Email:request.Email,
		SessionId:sid,
		About:request.About,
	}
	this.Serve200(resp)
}

func createSession(userAccount *models.UserAccount, accessToken string) (string, error) {
	sid := uuid.NewV4().String()
	beego.Info("Session Created|SID=", sid, "|UID=", (*userAccount).UID, "|Email=", userAccount.Email)

	err := models.AddUserAuth((*userAccount).UID, accessToken, sid)
	return sid, err
}
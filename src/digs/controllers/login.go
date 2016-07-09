package controllers

import (
	"digs/domain"
	"digs/models"
	"errors"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/satori/go.uuid"
	"strings"
)

type LoginController struct {
	HttpBaseController
}

func (this *LoginController) Post()  {
	var request domain.UserLoginRequest
	beego.Info("REQUEST|LoginRequest|", string(this.Ctx.Input.RequestBody))
	this.Super(&request.BaseRequest)
	json.Unmarshal(this.Ctx.Input.RequestBody, &request)
	//Check if the person is already registered
	userAccount, err := models.GetUserAccount("uid", request.FBID)
	if err != nil {
		beego.Error("Unable to get user Account|Err=", err)
		this.Serve500(errors.New("Unable to look up account table"))
		return
	}

	var sid, uid string
	if userAccount == nil {
		request.ProfilePicture = strings.Replace(request.ProfilePicture, "http://", "https://", 1)
		userAccount, err = models.AddUserAccount(request.FirstName, request.LastName, request.Email, request.About, request.FBID, request.Locale, request.ProfilePicture, request.FBVerified)
		if err != nil {
			beego.Error("Unable to create user Account|Err=", err)
			this.Serve500(err)
			return
		}
	}
	uid = userAccount.UID
	sid, err = createSession(userAccount, request.AccessToken)
	if sid == "" || err != nil {
		beego.Critical("SessionCreationFailed|err=", err)
		this.Serve500(errors.New("Unable to create new session"))
		return
	}

	resp := &domain.UserLoginResponse{
		StatusCode:200,
		SessionId:sid,
		UserId:uid,
		Settings:domain.SettingResponse{
			Range:userAccount.Settings.Range,
			PublicProfile:userAccount.Settings.PublicProfile,
			PushNotification:userAccount.Settings.PushNotification,
		},
	}
	beego.Info("Login Response=", resp)
	this.Serve200(resp)
}

func createSession(userAccount *models.UserAccount, accessToken string) (string, error) {
	sid := uuid.NewV4().String()
	beego.Info("SessionCreated|SID=", sid, "|UID=", userAccount.UID, "|Email=", userAccount.Email)

	err := models.AddUserAuth((*userAccount).UID, accessToken, sid)
	return sid, err
}
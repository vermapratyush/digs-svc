package controllers

import (
	"digs/domain"
	"digs/models"
	"errors"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/satori/go.uuid"
	"strings"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/aws"
	"fmt"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
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

		go sendWelcomeMail(userAccount)

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

func sendWelcomeMail(userAccount *models.UserAccount) {
	svc := ses.New(awsSession.New(), &aws.Config{Region: aws.String("eu-west-1")})
	content := `<!doctype html>
<html>
<head>
<meta name="viewport" content="width=device-width">
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
<title>Really Simple HTML Email Template</title>
<style>
/* -------------------------------------
    GLOBAL
------------------------------------- */
* {
  font-family: "Helvetica Neue", "Helvetica", Helvetica, Arial, sans-serif;
  font-size: 100%;
  line-height: 1.6em;
  margin: 0;
  padding: 0;
}

img {
  max-width: 600px;
  width: auto;
}

body {
  -webkit-font-smoothing: antialiased;
  height: 100%;
  -webkit-text-size-adjust: none;
  width: 100% !important;
}


/* -------------------------------------
    ELEMENTS
------------------------------------- */
a {
  color: #348eda;
}

.btn-primary {
  Margin-bottom: 10px;
  width: auto !important;
}

.btn-primary td {
  background-color: #348eda;
  border-radius: 25px;
  font-family: "Helvetica Neue", Helvetica, Arial, "Lucida Grande", sans-serif;
  font-size: 14px;
  text-align: center;
  vertical-align: top;
}

.btn-primary td a {
  background-color: #348eda;
  border: solid 1px #348eda;
  border-radius: 25px;
  border-width: 10px 20px;
  display: inline-block;
  color: #ffffff;
  cursor: pointer;
  font-weight: bold;
  line-height: 2;
  text-decoration: none;
}

.last {
  margin-bottom: 0;
}

.first {
  margin-top: 0;
}

.padding {
  padding: 10px 0;
}


/* -------------------------------------
    BODY
------------------------------------- */
table.body-wrap {
  padding: 20px;
  width: 100%;
}

table.body-wrap .container {
  border: 1px solid #f0f0f0;
}


/* -------------------------------------
    FOOTER
------------------------------------- */
table.footer-wrap {
  clear: both !important;
  width: 100%;
}

.footer-wrap .container p {
  color: #666666;
  font-size: 12px;

}

table.footer-wrap a {
  color: #999999;
}


/* -------------------------------------
    TYPOGRAPHY
------------------------------------- */
h1,
h2,
h3 {
  color: #111111;
  font-family: "Helvetica Neue", Helvetica, Arial, "Lucida Grande", sans-serif;
  font-weight: 200;
  line-height: 1.2em;
  margin: 40px 0 10px;
}

h1 {
  font-size: 36px;
}
h2 {
  font-size: 28px;
}
h3 {
  font-size: 22px;
}

p,
ul,
ol {
  font-size: 14px;
  font-weight: normal;
  margin-bottom: 10px;
}

ul li,
ol li {
  margin-left: 5px;
  list-style-position: inside;
}

/* ---------------------------------------------------
    RESPONSIVENESS
------------------------------------------------------ */

/* Set a max-width, and make it display as block so it will automatically stretch to that width, but will also shrink down on a phone or something */
.container {
  clear: both !important;
  display: block !important;
  Margin: 0 auto !important;
  max-width: 600px !important;
}

/* Set the padding on the td rather than the div for Outlook compatibility */
.body-wrap .container {
  padding: 20px;
}

/* This should also be a block element, so that it will fill 100% of the .container */
.content {
  display: block;
  margin: 0 auto;
  max-width: 600px;
}

/* Let's make sure tables in the content area are 100% wide */
.content table {
  width: 100%;
}

</style>
</head>
<body bgcolor="#f6f6f6">

<!-- body -->
<table class="body-wrap" bgcolor="#f6f6f6">
  <tr>
    <td></td>
    <td class="container" bgcolor="#FFFFFF">

      <!-- content -->
      <div class="content">
      <table>
        <tr>
          <td>
          <center>
          <a href="http://powow.info"><img src="https://raw.githubusercontent.com/PowowInfo/powowinfo.github.io/master/img/icon.png"  target="_blank" style="width:90px; height:90px;margin:10px" /></a>
          </center>
            <h3>Hi ` + userAccount.FirstName + `,</h3>
            <p>First of all thank you for signing up on <a href="http://powow.info" target="_blank">Powow</a> and welcome to the Powow community. Powow is a your geo-location based community with no central organizer. You decide how wide or narrow your comfort radius for community is and you are implicitly part of it.</p>

            <h3>How is Powow different?</h3>
            <p>In contrast to previous generations which lived, worked and enjoyed at fixed locations current generation is mobile, dynamic and globally minded. But, moving is still too hard. Not any more! Powow allows you to immediately discover people, events and activities happening in your locality by simply getting on the app.</p>

			<h3>Few examples how you can make use of Powow:</h3>
			<ul>
			<li>You are at a gathering and want to share something with people there without asking everyone to jump hoops of asking emails, creating groups or any other thing.</li>
			<li>You just met someone and want to connect with them for future purposes but, don’t want to go through collecting and giving those antique dreaded visiting cards in 21st century.</li>
			<li>You want to organize an event in your society or hostel, and want to connect with interested people and plan further action plan.</li>
			</ul>
            <!-- button -->

            <!-- /button -->
            <p>If you would like to know more about what we are up to, follow us on <a href="https://medium.com/powowinfo/introducing-powow-693036ca8c0f">Medium.</a></p>
            <p>We wish you best and hope you do love using Powow. In case of any feedback please do get in touch with us, we would be happy to listen to you.</p>
            <div style="margin-top:40px"></div>
            <p style="font-style:italic">Thanks, have a lovely day.</p>
            <p style="font-style:italic">The Powow Team</p>
            <div>
            <p style="float: left; margin-right:10px"><a href="http://twitter.com/PowowInfo"  target="_blank"><img src="https://raw.githubusercontent.com/PowowInfo/powowinfo.github.io/master/img/twitter-icon.gif"></img></p>
            <p style="float: left; margin-right:10px"><a href="http://facebook.com/PowowInfo"  target="_blank"><img src="https://raw.githubusercontent.com/PowowInfo/powowinfo.github.io/master/img/facebook-icon.gif"></img></a></p>
            <p style="float: left; margin-right:10px;"><a href="mailto:hey@powow.info" target="_blank"><img  src="https://raw.githubusercontent.com/PowowInfo/powowinfo.github.io/master/img/mail-icon.png"></img></a></p>
            <p style="float: left; margin-right:10px;"><a href="https://medium.com/powowinfo/introducing-powow-693036ca8c0f" target="_blank"><img style="margin-top:-3px" src="https://raw.githubusercontent.com/PowowInfo/powowinfo.github.io/master/img/blog-icon.png"></img></a></p>
            <div style="clear: left;
"></div>
            </div>
          </td>
        </tr>
      </table>
      </div>
      <!-- /content -->

    </td>
    <td></td>
  </tr>
</table>
<!-- /body -->

<!-- footer -->
<table class="footer-wrap">
  <tr>
    <td></td>
    <td class="container">


    </td>
    <td></td>
  </tr>
</table>
<!-- /footer -->

</body>
</html>
`

	params := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(userAccount.Email),
			},
		},
		Message: &ses.Message{ // Required
			Body: &ses.Body{ // Required
				Html: &ses.Content{
					Data:    aws.String(content),
				},
				//Text: &ses.Content{
				//	Data:    aws.String("MessageData"),
				//},
			},
			Subject: &ses.Content{ // Required
				Data:    aws.String("Welcome to Powow"),
			},
		},
		Source: aws.String("hey@powow.info"),
		ReplyToAddresses: []*string{
			aws.String("hey@powow.info"),
		},
	}
	resp, err := svc.SendEmail(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

package models

import (
	"github.com/afex/hystrix-go/hystrix"
	"digs/common"
	"github.com/astaxie/beego"
	"github.com/sideshow/apns2/certificate"
	apns "github.com/sideshow/apns2"
	"github.com/NaySoftware/go-fcm"
)
var (
	APNSCert, pemErr = certificate.FromPemFile("models/final.pem", "")
)


func AndroidMessagePush(uid, nid, message string)  {
	_ = hystrix.Go(common.AndroidPush, func() error {
		beego.Info("AndroidPush|From=", uid, "|To=", nid)
		fcm := fcm.NewFcmClient(common.PushNotification_API_KEY)

		data := map[string]string{
			"title": "powow",
			"message": message,
			"image": "twitter",
			"style": "inbox",
			"summaryText": "There are %n% notification",
		}

		fcm.NewFcmMsgTo(nid, data)
		status, err := fcm.Send(1)
		if err == nil {
			beego.Info(status)
		} else {
			beego.Error(err)
		}
		return err
	}, nil)
}

func IOSMessagePush(uid, nid, message string) {
	_ = hystrix.Go(common.IOSPush, func() error {
		beego.Info("IOSPush|From=", uid, "|To=", nid)

		if pemErr != nil {
			beego.Error("APNSCertError|err", pemErr)
			return pemErr
		}

		notification := &apns.Notification{}
		notification.DeviceToken = nid
		notification.Topic = "info.powow.app"
		notification.Payload = []byte("{\"aps\":{\"alert\":\"" + message + "\"}}") // See Payload section below

		client := apns.NewClient(APNSCert).Development()
		resp, err := client.Push(notification)

		if err != nil {
			beego.Error("APNSPushError|err=", err)
		} else {
			beego.Info("APNSResponse=", resp)
		}
		return err
	}, nil)
}

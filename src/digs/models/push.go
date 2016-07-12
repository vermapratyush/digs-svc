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
	APNSCert, pemErr = certificate.FromPemFile("models/apn_production.pem", "")
)


func AndroidMessagePush(uid, nid, message string)  {
	_ = hystrix.Go(common.AndroidPush, func() error {
		fcm := fcm.NewFcmClient(common.PushNotification_API_KEY)

		data := map[string]string{
			"title": "powow",
			"message": message,
			"image": "twitter",
			"style": "inbox",
			"summaryText": "There are %n% notification",
		}

		fcm.NewFcmMsgTo(nid, data)
		_, err := fcm.Send(1)
		if err != nil {
			beego.Error(err)
		}
		return err
	}, nil)
}

func IOSMessagePush(uid, nid, message string) {
	_ = hystrix.Go(common.IOSPush, func() error {

		if pemErr != nil {
			beego.Error("APNSCertError|err", pemErr)
			return pemErr
		}

		notification := &apns.Notification{}
		notification.DeviceToken = nid
		notification.Topic = "info.powow.app"
		notification.Payload = []byte("{\"aps\":{\"alert\":\"" + message + "\"}}") // See Payload section below

		client := apns.NewClient(APNSCert).Production()
		_, err := client.Push(notification)

		if err != nil {
			beego.Error("APNSPushError|err=", err)
		}
		return err
	}, nil)
}

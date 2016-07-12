package models

import (
	"github.com/afex/hystrix-go/hystrix"
	"digs/common"
	"github.com/sideshow/apns2/certificate"
	apns "github.com/sideshow/apns2"
	"github.com/NaySoftware/go-fcm"
	"digs/logger"
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
		response, err := fcm.Send(1)
		if err != nil || response.StatusCode != 200 {
			logger.Error("PUSH|Android|UID=", uid, "|NID=", nid, "|OS=", "|Response=", response, "|Err=", err)
		}
		return err
	}, nil)
}

func IOSMessagePush(uid, nid, message string) {
	_ = hystrix.Go(common.IOSPush, func() error {

		if pemErr != nil {
			logger.Error("APNSCertError|err", pemErr)
			return pemErr
		}

		notification := &apns.Notification{}
		notification.DeviceToken = nid
		notification.Topic = "info.powow.app"
		notification.Payload = []byte("{\"aps\":{\"alert\":\"" + message + "\"}}") // See Payload section below

		client := apns.NewClient(APNSCert).Production()
		response, err := client.Push(notification)

		if err != nil || response.StatusCode != 200 {
			logger.Error("PUSH|ios|UID=", uid, "|NID=", nid, "|OS=", "|Response=", response, "|Err=", err)
		}
		return err
	}, nil)
}

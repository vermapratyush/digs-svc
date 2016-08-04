package models

import (
	"github.com/afex/hystrix-go/hystrix"
	"digs/common"
	"github.com/sideshow/apns2/certificate"
	apns "github.com/sideshow/apns2"
	"github.com/NaySoftware/go-fcm"
	"digs/logger"
	"encoding/json"
)
var (
	APNSCert, pemErr = certificate.FromPemFile("models/apn_production.pem", "")
)

type Payload struct {
	Aps *APS `json:"aps"`
}
type APS struct {
	Alert string `json:"alert"`
}

func AndroidMessagePush(uid, nid, message, additionalData string)  {
	_ = hystrix.Go(common.AndroidPush, func() error {
		fcm := fcm.NewFcmClient(common.PushNotification_API_KEY)

		data := map[string]string{
			"title": "powow",
			"message": message,
			"image": "https://raw.githubusercontent.com/PowowInfo/powowinfo.github.io/master/img/image_300.png",
			"style": "inbox",
			"additionalData": additionalData,
			"content-available": "1",
			"summaryText": "There are %n% notification",
		}

		fcm.SetContentAvailable(true)
		fcm.NewFcmMsgTo(nid, data)
		response, err := fcm.Send(1)
		if err != nil || response.StatusCode != 200 {
			logger.Error("PUSH|Android|UID=", uid, "|NID=", nid, "|OS=", "|Response=%v", response, "|Err=%v", err)
		}
		return err
	}, nil)
}

func IOSMessagePush(uid, nid, message string) {
	_ = hystrix.Go(common.IOSPush, func() error {

		if pemErr != nil {
			logger.Error("APNSCertError|err=%v", pemErr)
			return pemErr
		}

		notification := &apns.Notification{}
		notification.DeviceToken = nid
		notification.Topic = "info.powow.app"
		payload := &Payload{
			Aps:&APS{
				Alert:message,
			},
		}
		payloadByte, _ := json.Marshal(payload)
		notification.Payload = payloadByte
		client := apns.NewClient(APNSCert).Production()
		response, err := client.Push(notification)
		if err != nil || response.StatusCode != 200 {
			logger.Error("PUSH|ios|UID=", uid, "|NID=", nid, "|OS=", "|Response=%v", response, "|Err=%v", err)
		}
		return err
	}, nil)
}

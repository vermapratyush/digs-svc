package controllers

import (
	"github.com/satori/go.uuid"
	"digs/domain"
	"fmt"
	"github.com/astaxie/beego"
	"digs/models"
	"digs/logger"
	"strings"
	"gopkg.in/mgo.v2"
)

type MediaController struct {
	HttpBaseController
}

func (this *MediaController) Put()  {

	file, fileHeader, err := this.GetFile("picture")
	if err != nil {
		beego.Debug(err)
	}

	if !strings.HasPrefix(fileHeader.Header.Get("Content-Type"), "image") {
		this.ServeUnsupportedMedia()
		return
	}

	if err != nil {
		logger.Error("FileUploadFailed|Err=", err)
		this.Serve500(err)
		return
	}

	sid := this.GetString("sessionId")
	userAuth, err := models.FindSession("sid", sid)
	if err != nil {
		if err == mgo.ErrNotFound {
			this.InvalidSessionResponse()
			return
		}
		logger.Error("LOGOUT|LogoutFailed=%v", err)
		this.Serve500(err)
		return
	}

	blobUUID := uuid.NewV4().String()
	err = models.PutS3Object(file, blobUUID, fileHeader.Header.Get("Content-Type"), userAuth.UID)

	if err != nil {
		this.Serve500(err)
		return
	}
	resource := domain.MessagePutResponse{
		ResourceUrl:fmt.Sprintf("https://s3-eu-west-1.amazonaws.com/%s/%s", "powow-file-sharing", blobUUID),
	}
	this.Serve200(resource)
	return
}


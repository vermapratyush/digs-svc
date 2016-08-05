package controllers

import (
	"github.com/satori/go.uuid"
	"digs/domain"
	"fmt"
	"digs/models"
	"digs/logger"
	"gopkg.in/mgo.v2"
	"mime/multipart"
	"image/jpeg"
	"image/png"
	"os"
	"github.com/nfnt/resize"
	"image"
	"io/ioutil"
	"mime"
)

type MediaController struct {
	HttpBaseController
}

//Width, Height, Quality
var ResizeOptions = map[string][]uint{
	"original": []uint{},
	"thumbnail": []uint{50, 50, 100},
	"medium": []uint{300, 300, 75},
}

func (this *MediaController) Put()  {

	file, fileHeader, err := this.GetFile("picture")
	fileExt := fileHeader.Header.Get("Content-Type")

	if !isFileTypeSupported(fileExt) {
		this.ServeUnsupportedMedia()
		return
	}

	if err != nil {
		logger.Error("ImageUpload|Err=", err)
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
		logger.Error("ImageUpload|SessionError|err=", err)
		this.Serve500(err)
		return
	}

	blobUUID := uuid.NewV4().String()
	_, err = processAndUploadImages(file, fileExt, blobUUID, userAuth.UID)

	if err != nil {
		this.Serve500(err)
		return
	}
	resource := domain.MessagePutResponse{
		ResourceUrl:fmt.Sprintf("https://s3-eu-west-1.amazonaws.com/%s/%s-original", "powow-file-sharing", blobUUID),
	}
	this.Serve200(resource)
	return
}

func processAndUploadImages(file multipart.File, fileExt string, outputFile string, uid string) (*os.File, error) {

	originalFileName := fmt.Sprintf(fmt.Sprintf("/tmp/%s-original.%s", outputFile, fileExt))
	resizedFileName := fmt.Sprintf(fmt.Sprintf("/tmp/%s.%s", outputFile, fileExt))
	originalFile, err := os.Create(originalFileName)

	defer func() {
		os.Remove(originalFileName)
		os.Remove(resizedFileName)
		originalFile.Close()
	}()

	inputBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Critical("ImageUpload|FileReadError|Err=", err)
		return nil, err
	}
	_, err = originalFile.Write(inputBytes)
	if err != nil {
		logger.Critical("ImageUpload|FileWriteErorr|Err=", err)
		return nil, err
	}

	for fileQualityPrefix, imageQuality := range(ResizeOptions) {
		newBlobName := fmt.Sprintf("%s-%s", outputFile, fileQualityPrefix)

		originalImageCopy, _ := os.Open(originalFile.Name())
		originalImg, _, err := image.Decode(originalImageCopy)
		if err != nil {
			logger.Critical("ImageUpload|FileDecodeError|Err=", err)
			return nil, nil
		}
		resizeImageFile, _ := os.Create(resizedFileName)

		if len(imageQuality) != 3 || fileExt == "gif" {
			err = models.PutS3Object(originalFile, newBlobName, mime.TypeByExtension(fileExt), uid)
			if err != nil {
				logger.Critical("ImageUpload|S3UploadFailed|err=", err)
			}
			continue
		}

		resizedImage := resize.Thumbnail(imageQuality[0], imageQuality[1], originalImg, resize.Lanczos3)
		if err != nil {
			return nil, err
		}
		if fileExt == "jpeg" || fileExt == "jpg" {
			jpeg.Encode(resizeImageFile, resizedImage, &jpeg.Options{Quality: int(imageQuality[2])})
		}
		if fileExt == "png" {
			png.Encode(resizeImageFile, resizedImage)
		}

		err = models.PutS3Object(resizeImageFile, newBlobName, "image/webp", uid)
		if err != nil {
			logger.Critical("ImageUpload|S3UploadFailed|err=", err)
		}

		resizeImageFile.Close()
	}

	return nil, nil
}

func isFileTypeSupported(fileExt string) bool {
	if (fileExt == "jpeg" || fileExt == "jpg" || fileExt == "png" || fileExt == "gif") {
		return true
	}
	return false
}
package models

import (
	"github.com/afex/hystrix-go/hystrix"
	"digs/common"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"digs/logger"
	"mime/multipart"
)

func PutS3Object(data multipart.File, key, contentType, uid string) error {
	svc := s3.New(session.New(), aws.NewConfig().WithRegion("eu-west-1"))

	err := hystrix.Do(common.AmazonS3, func() error {
		params := s3.PutObjectInput{
			Bucket: aws.String(common.AmazonS3BucketName),
			Key: aws.String(key),
			Metadata: map[string]*string{
				"user-id": aws.String(uid),
			},
			ContentType: aws.String(contentType),
			Body:data,
		}

		_, err := svc.PutObject(&params)

		if err != nil {
			logger.Error("S3Put|Key=", key, "|Err=", err)
		}
		return err
	}, nil)

	return err
}

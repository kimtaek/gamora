package aws

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/caarlos0/env"
	"github.com/kimtaek/gamora/pkg/helper"
	"mime/multipart"
	"strings"
)

type Configure struct {
	Mode               string `env:"APP_MODE" envDefault:"debug"`
	AccessKey          string `env:"AWS_ACCESS_KEY" envDefault:"AKIAZI..."`
	SecretAccessKey    string `env:"AWS_SECRET_KEY" envDefault:"MT+8FpEPD..."`
	S3Region           string `env:"AWS_S3_REGION" envDefault:"ap-northeast-1"`
	S3Bucket           string `env:"AWS_S3_BUCKET" envDefault:"s3.....com"`
	UniversalTelephone string `env:"UNIVERSAL_TELEPHONE" envDefault:""`
}

var Config Configure

func Setup() {
	_ = env.Parse(&Config)
}

type manager struct {
	Location string
	Session  *session.Session
	Uploader *s3manager.Uploader
}

func NewAws() *manager {
	s := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials: credentials.NewStaticCredentialsFromCreds(credentials.Value{
				AccessKeyID:     Config.AccessKey,
				SecretAccessKey: Config.SecretAccessKey,
			}),
			Region: aws.String(Config.S3Region),
		},
	}))

	return &manager{
		Session:  s,
		Uploader: s3manager.NewUploader(s),
	}
}

func (a *manager) SendSMS(phoneNumber string, message string) error {
	if Config.UniversalTelephone == "82-1000000000" {
		return nil
	}

	if Config.UniversalTelephone != "" {
		phoneNumber = Config.UniversalTelephone
	}

	service := sns.New(a.Session)
	params := &sns.PublishInput{
		Message:     aws.String(message),
		PhoneNumber: aws.String(phoneNumber),
	}
	resp, err := service.Publish(params)

	if err != nil {
		helper.Error(err.Error())
		return err
	}

	helper.Info(resp)
	return nil
}

func (a *manager) Upload(file multipart.File, fileName string, extension string) (url string, version *string, err error) {
	if fileName == "" {
		return "", nil, errors.New("file name not found")
	}

	var contentType string
	switch strings.ToLower(extension) {
	case ".jpg":
		contentType = "image/jpeg"
	case ".jpeg":
		contentType = "image/jpeg"
	case ".gif":
		contentType = "image/gif"
	case ".png":
		contentType = "image/png"
	default:
		return "", nil, errors.New("not allow the file extension")
	}

	result, err := a.Uploader.Upload(&s3manager.UploadInput{
		ACL:         aws.String("public-read"),
		Body:        file,
		Bucket:      aws.String(Config.S3Bucket),
		ContentType: aws.String(contentType),
		Key:         aws.String(a.Location + "/" + fileName),
	})

	if err != nil {
		return "", nil, fmt.Errorf("failed to upload file, %v", err)
	}

	return result.Location, result.VersionID, nil
}

func (a *manager) Delete(location string, name string, version *string) error {
	service := s3.New(a.Session)
	_, _ = service.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(Config.S3Bucket), Key: aws.String(location + "/" + name), VersionId: version})
	return service.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(Config.S3Bucket),
		Key:    aws.String(location + "/" + name),
	})
}

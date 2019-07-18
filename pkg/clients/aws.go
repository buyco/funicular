package clients

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/buyco/funicular/internal/utils"
	"io"
	"sync"
)

func NewAWSConfig() *aws.Config {
	return aws.NewConfig()
}

func NewAWSSession(config *aws.Config) *session.Session {
	return session.Must(session.NewSession(config))
}

func NewS3Client(session *session.Session) *s3.S3 {
	return s3.New(session)
}

//------------------------------------------------------------------------------

type AWSManager struct {
	session      *session.Session
	disconnected chan bool
	closed       bool
	S3Manager    *S3Manager
}

func NewAWSManager(session *session.Session) *AWSManager {
	return &AWSManager{
		session:   session,
		S3Manager: NewS3Manager(session),
	}
}

//------------------------------------------------------------------------------

type S3Manager struct {
	session *session.Session
	client  *s3.S3
	S3      []StorageAccessLayer
	sync.Mutex
}

func NewS3Manager(session *session.Session) *S3Manager {
	return &S3Manager{
		client: NewS3Client(session),
		S3:     make([]StorageAccessLayer, 0),
	}
}

func (sm *S3Manager) Add(bucketName string) *S3Wrapper {
	sm.Lock()
	defer sm.Unlock()
	s3Wrapper := NewS3Wrapper(bucketName, sm.client)
	sm.S3 = append(sm.S3, s3Wrapper)

	return s3Wrapper
}

//------------------------------------------------------------------------------

// Storage Layer interface
type StorageAccessLayer interface {
	Upload(path string, filename string, data io.Reader) (string, error)
	Download(path string, filename string, data io.WriterAt) (int64, error)
	Read(path string, filename string, limit int64, readFrom string) (*s3.ListObjectsV2Output, error)
}

// S3 Adapter
type S3Wrapper struct {
	bucketName string
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	reader     func(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error)
}

func NewS3Wrapper(bucketName string, s3Client *s3.S3) *S3Wrapper {
	uploader := s3manager.NewUploaderWithClient(s3Client)
	downloader := s3manager.NewDownloaderWithClient(s3Client)
	reader := s3Client.ListObjectsV2
	return &S3Wrapper{
		bucketName: bucketName,
		uploader:   uploader,
		downloader: downloader,
		reader:     reader,
	}
}

func (s3w *S3Wrapper) Upload(path string, filename string, data io.Reader) (string, error) {
	upParams := &s3manager.UploadInput{
		Bucket: aws.String(s3w.bucketName),
		Key:    aws.String(path + filename),
		Body:   data,
	}
	result, err := s3w.uploader.Upload(upParams)
	var location string
	if result != nil {
		location = result.Location
	}
	return location, err
}

func (s3w *S3Wrapper) Download(path string, filename string, data io.WriterAt) (int64, error) {
	downParams := &s3.GetObjectInput{
		Bucket:  aws.String(s3w.bucketName),
		Key:    aws.String(path + filename),
	}
	err := downParams.Validate()
	if err != nil {
		return 0, utils.ErrorPrintf("Download params malformed: %v", err)
	}
	result, err := s3w.downloader.Download(data, downParams)
	return result, err
}

func (s3w *S3Wrapper) Read(path string, filename string, limit int64, readFrom string) (*s3.ListObjectsV2Output, error) {
	readParams := &s3.ListObjectsV2Input{
		Bucket:  aws.String(s3w.bucketName),
		Prefix: aws.String(path),
		MaxKeys: aws.Int64(limit),
	}

	if readFrom != "" {
		readParams.SetStartAfter(readFrom)
	}

	result, err := s3w.reader(readParams)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return nil, utils.ErrorPrint(fmt.Sprint(s3.ErrCodeNoSuchBucket, aerr.Error()))

			default:
				return nil, aerr
			}
		}
	}
	return result, err
}

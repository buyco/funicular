package client

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/buyco/keel/pkg/helper"
	"io"
	"sync"
)

// NewAWSConfig is AWS Config constructor
func NewAWSConfig() *aws.Config {
	return aws.NewConfig()
}

// NewAWSSession is AWS Session constructor
func NewAWSSession(config *aws.Config) *session.Session {
	return session.Must(session.NewSession(config))
}

// NewS3Client is AWS S3 constructor
func NewS3Client(session *session.Session) *s3.S3 {
	return s3.New(session)
}

//------------------------------------------------------------------------------

// AWSManager is a struct to manage AWS SDK
type AWSManager struct {
	session      *session.Session
	disconnected chan bool
	closed       bool
	S3Manager    *S3Manager
}

// NewAWSManager is AWSManager constructor
func NewAWSManager(session *session.Session) *AWSManager {
	return &AWSManager{
		session:   session,
		S3Manager: NewS3Manager(session),
	}
}

//------------------------------------------------------------------------------

// S3Manager is a struct to control S3 client
type S3Manager struct {
	session *session.Session
	client  *s3.S3
	S3      []StorageAccessLayer
	sync.Mutex
}

// NewS3Manager is a S3Manager constructor
func NewS3Manager(session *session.Session) *S3Manager {
	return &S3Manager{
		client: NewS3Client(session),
		S3:     make([]StorageAccessLayer, 0),
	}
}

// Add is used to add a bucket to the manager
func (sm *S3Manager) Add(bucketName string) *S3Wrapper {
	sm.Lock()
	defer sm.Unlock()
	s3Wrapper := NewS3Wrapper(bucketName, sm.client)
	sm.S3 = append(sm.S3, s3Wrapper)

	return s3Wrapper
}

//------------------------------------------------------------------------------

// StorageAccessLayer is a common interface for AWS storage
// Specific to S3 for now...
type StorageAccessLayer interface {
	Upload(path string, filename string, data io.Reader) (string, error)
	Download(path string, filename string, data io.WriterAt) (int64, error)
	Read(path string, limit int64, readFrom string) (*s3.ListObjectsV2Output, error)
	Delete(path string, files ...string) (*s3.DeleteObjectsOutput, error)
}

// S3Wrapper is a S3 Adapter
type S3Wrapper struct {
	bucketName string
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	deleter    func(input *s3.DeleteObjectsInput) (*s3.DeleteObjectsOutput, error)
	reader     func(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error)
}

// NewS3Wrapper is S3Wrapper struct constructor
func NewS3Wrapper(bucketName string, s3Client *s3.S3) *S3Wrapper {
	uploader := s3manager.NewUploaderWithClient(s3Client)
	downloader := s3manager.NewDownloaderWithClient(s3Client)
	deleter := s3Client.DeleteObjects
	reader := s3Client.ListObjectsV2
	return &S3Wrapper{
		bucketName: bucketName,
		uploader:   uploader,
		downloader: downloader,
		reader:     reader,
		deleter:    deleter,
	}
}

// Upload pushes a file to S3
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

// Download fetches a file from S3
func (s3w *S3Wrapper) Download(path string, filename string, data io.WriterAt) (int64, error) {
	downParams := &s3.GetObjectInput{
		Bucket: aws.String(s3w.bucketName),
		Key:    aws.String(path + filename),
	}
	err := downParams.Validate()
	if err != nil {
		return 0, helper.ErrorPrintf("Download params malformed: %v", err)
	}
	result, err := s3w.downloader.Download(data, downParams)
	return result, err
}

// Delete drops files from S3
func (s3w *S3Wrapper) Delete(path string, filename ...string) (*s3.DeleteObjectsOutput, error) {
	var objects []*s3.ObjectIdentifier
	for _, file := range filename {
		objects = append(objects, &s3.ObjectIdentifier{Key: aws.String(path + file)})
	}
	input := &s3.DeleteObjectsInput{
		Bucket: aws.String(s3w.bucketName),
		Delete: &s3.Delete{
			Objects: objects,
			Quiet:   aws.Bool(false),
		},
	}
	err := input.Validate()
	if err != nil {
		return nil, helper.ErrorPrintf("Delete params malformed: %v", err)
	}

	result, err := s3w.deleter(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return nil, helper.ErrorPrint(fmt.Sprint(s3.ErrCodeNoSuchBucket, aerr.Error()))
			case s3.ErrCodeNoSuchKey:
				return nil, helper.ErrorPrint(fmt.Sprint(s3.ErrCodeNoSuchKey, aerr.Error()))
			default:
				return nil, aerr
			}
		} else {
			return nil, helper.ErrorPrint(aerr.Error())
		}
	}
	return result, err
}

// Read gets file content from S3
func (s3w *S3Wrapper) Read(path string, limit int64, readFrom string) (*s3.ListObjectsV2Output, error) {
	readParams := &s3.ListObjectsV2Input{
		Bucket:  aws.String(s3w.bucketName),
		Prefix:  aws.String(path),
		MaxKeys: aws.Int64(limit),
	}
	if readFrom != "" {
		readParams.SetStartAfter(readFrom)
	}
	err := readParams.Validate()
	if err != nil {
		return nil, helper.ErrorPrintf("Read params malformed: %v", err)
	}

	result, err := s3w.reader(readParams)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return nil, helper.ErrorPrint(fmt.Sprint(s3.ErrCodeNoSuchBucket, aerr.Error()))

			default:
				return nil, aerr
			}
		} else {
			return nil, helper.ErrorPrint(aerr.Error())
		}
	}
	return result, err
}

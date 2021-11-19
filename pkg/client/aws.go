package client

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/buyco/keel/pkg/helper"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"io"
	"sync"
	"time"
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

type DownloadOptions struct {
	IfMatch                    string
	IfModifiedSince            *time.Time
	IfNoneMatch                string
	IfUnmodifiedSince          *time.Time
	ResponseCacheControl       string
	ResponseContentDisposition string
	ResponseContentEncoding    string
	ResponseContentLanguage    string
	ResponseContentType        string
	ResponseExpires            *time.Time
	VersionId                  string
}

type UploadOptions struct {
	ACL                       string
	CacheControl              string
	ContentDisposition        string
	ContentEncoding           string
	ContentLanguage           string
	ContentMD5                string
	ContentType               string
	Expires                   *time.Time
	GrantFullControl          string
	GrantRead                 string
	GrantReadACP              string
	GrantWriteACP             string
	Metadata                  map[string]string
	ObjectLockMode            string
	ObjectLockRetainUntilDate *time.Time
	ServerSideEncryption      string
	StorageClass              string
	Tagging                   string
}

// AWSManager is a struct to manage AWS SDK
type AWSManager struct {
	session   *session.Session
	S3Manager *S3Manager
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
	client *s3.S3
	s3     map[string]StorageAccessLayer
	mutex  sync.RWMutex
}

// NewS3Manager is a S3Manager constructor
func NewS3Manager(session *session.Session) *S3Manager {
	return &S3Manager{
		client: NewS3Client(session),
		s3:     make(map[string]StorageAccessLayer),
	}
}

// Add is used to add a bucket to the manager
func (sm *S3Manager) Add(bucketName string) *S3Wrapper {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	if sm.s3[bucketName] == nil {
		s3Wrapper := NewS3Wrapper(bucketName, sm.client)
		sm.s3[bucketName] = s3Wrapper
	}

	return sm.s3[bucketName].(*S3Wrapper)
}

// Get is used to fetch a bucket from the manager
func (sm *S3Manager) Get(bucketName string) *S3Wrapper {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	if sm.s3[bucketName] == nil {
		return nil
	}
	return sm.s3[bucketName].(*S3Wrapper)
}

// GetAll is used to fetch all buckets from the manager
func (sm *S3Manager) GetAll() map[string]StorageAccessLayer {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	return sm.s3
}

// Delete is used to delete a bucket from the manager
func (sm *S3Manager) Delete(bucketName string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	if sm.s3[bucketName] == nil {
		return helper.ErrorPrintf("bucket [%s] does not exist", bucketName)
	}
	delete(sm.s3, bucketName)
	return nil
}

//------------------------------------------------------------------------------

// StorageAccessLayer is a common interface for AWS storage
// Specific to S3 for now...
type StorageAccessLayer interface {
	Upload(path string, filename string, data io.Reader, options *UploadOptions) (string, error)
	Download(path string, filename string, data io.WriterAt, options *DownloadOptions) (int64, error)
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
func (s3w *S3Wrapper) Upload(path string, filename string, data io.Reader, options *UploadOptions) (string, error) {
	upParams, err := s3w.mergeUploadOptions(
		&s3manager.UploadInput{
			Bucket: aws.String(s3w.bucketName),
			Key:    aws.String(path + filename),
			Body:   data,
		},
		options,
	)
	if err != nil {
		return "", err
	}
	// We should implement multipart later on
	result, err := s3w.uploader.Upload(upParams)
	var location string
	if result != nil {
		location = result.Location
	}
	return location, err
}

// Generate a new copy of UploadInput filled with options
func (s3w *S3Wrapper) mergeUploadOptions(s3Params *s3manager.UploadInput, options *UploadOptions) (*s3manager.UploadInput, error) {
	if s3Params == nil {
		return nil, errors.New("s3Params argument must be of type UploadInput")
	}
	var newS3Input s3manager.UploadInput
	err := copier.Copy(&newS3Input, &s3Params)
	if err != nil {
		return nil, err
	}
	if options == nil {
		return &newS3Input, nil
	}
	if len(options.ContentType) > 0 {
		newS3Input.ContentType = aws.String(options.ContentType)
	}
	if options.Expires != nil {
		newS3Input.Expires = options.Expires
	}
	if len(options.ACL) > 0 {
		newS3Input.ACL = aws.String(options.ACL)
	}
	if options.ObjectLockRetainUntilDate != nil {
		newS3Input.ObjectLockRetainUntilDate = options.ObjectLockRetainUntilDate
	}
	if len(options.CacheControl) > 0 {
		newS3Input.CacheControl = aws.String(options.CacheControl)
	}
	if len(options.ContentDisposition) > 0 {
		newS3Input.ContentDisposition = aws.String(options.ContentDisposition)
	}
	if len(options.ContentEncoding) > 0 {
		newS3Input.ContentEncoding = aws.String(options.ContentEncoding)
	}
	if len(options.ContentLanguage) > 0 {
		newS3Input.ContentLanguage = aws.String(options.ContentLanguage)
	}
	if len(options.ContentMD5) > 0 {
		newS3Input.ContentMD5 = aws.String(options.ContentMD5)
	}
	if len(options.GrantFullControl) > 0 {
		newS3Input.GrantFullControl = aws.String(options.GrantFullControl)
	}
	if len(options.GrantRead) > 0 {
		newS3Input.GrantRead = aws.String(options.GrantRead)
	}
	if len(options.GrantReadACP) > 0 {
		newS3Input.GrantReadACP = aws.String(options.GrantReadACP)
	}
	if len(options.Metadata) > 0 {
		newS3Input.Metadata = aws.StringMap(options.Metadata)
	}
	if len(options.GrantWriteACP) > 0 {
		newS3Input.GrantWriteACP = aws.String(options.GrantWriteACP)
	}
	if len(options.ObjectLockMode) > 0 {
		newS3Input.ObjectLockMode = aws.String(options.ObjectLockMode)
	}
	if len(options.ServerSideEncryption) > 0 {
		newS3Input.ServerSideEncryption = aws.String(options.ServerSideEncryption)
	}
	if len(options.StorageClass) > 0 {
		newS3Input.StorageClass = aws.String(options.StorageClass)
	}
	if len(options.Tagging) > 0 {
		newS3Input.Tagging = aws.String(options.Tagging)
	}
	return &newS3Input, nil
}

// Download fetches a file from S3
func (s3w *S3Wrapper) Download(path string, filename string, data io.WriterAt, options *DownloadOptions) (int64, error) {
	downParams, err := s3w.mergeDownloadOptions(
		&s3.GetObjectInput{
			Bucket: aws.String(s3w.bucketName),
			Key:    aws.String(path + filename),
		},
		options,
	)
	if err != nil {
		return 0, err
	}
	if downParams.Validate() != nil {
		return 0, errors.Errorf("download params malformed: %v", err)
	}
	result, err := s3w.downloader.Download(data, downParams)
	return result, err
}

// Generate a new copy of GetObjectInput filled with options
func (s3w *S3Wrapper) mergeDownloadOptions(s3Params *s3.GetObjectInput, options *DownloadOptions) (*s3.GetObjectInput, error) {
	if s3Params == nil {
		return nil, errors.New("s3Params argument must be of type GetObjectInput")
	}
	var newS3Input s3.GetObjectInput
	err := copier.Copy(&newS3Input, &s3Params)
	if err != nil {
		return nil, err
	}
	if options == nil {
		return &newS3Input, nil
	}
	if len(options.IfMatch) > 0 {
		newS3Input.IfMatch = aws.String(options.IfMatch)
	}
	if options.IfModifiedSince != nil {
		newS3Input.IfModifiedSince = options.IfModifiedSince
	}
	if len(options.IfNoneMatch) > 0 {
		newS3Input.IfNoneMatch = aws.String(options.IfNoneMatch)
	}
	if options.IfUnmodifiedSince != nil {
		newS3Input.IfUnmodifiedSince = options.IfUnmodifiedSince
	}
	if len(options.ResponseCacheControl) > 0 {
		newS3Input.ResponseCacheControl = aws.String(options.ResponseCacheControl)
	}
	if len(options.ResponseContentDisposition) > 0 {
		newS3Input.ResponseContentDisposition = aws.String(options.ResponseContentDisposition)
	}
	if len(options.ResponseContentEncoding) > 0 {
		newS3Input.ResponseContentEncoding = aws.String(options.ResponseContentEncoding)
	}
	if len(options.ResponseContentLanguage) > 0 {
		newS3Input.ResponseContentLanguage = aws.String(options.ResponseContentLanguage)
	}
	if len(options.ResponseContentType) > 0 {
		newS3Input.ResponseContentType = aws.String(options.ResponseContentType)
	}
	if options.ResponseExpires != nil {
		newS3Input.ResponseExpires = options.ResponseExpires
	}
	if len(options.VersionId) > 0 {
		newS3Input.VersionId = aws.String(options.VersionId)
	}
	return &newS3Input, nil
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
		return nil, errors.Errorf("delete params malformed: %v", err)
	}

	result, err := s3w.deleter(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return nil, errors.New(fmt.Sprint(s3.ErrCodeNoSuchBucket, aerr.Error()))
			case s3.ErrCodeNoSuchKey:
				return nil, errors.New(fmt.Sprint(s3.ErrCodeNoSuchKey, aerr.Error()))
			default:
				return nil, aerr
			}
		} else {
			return nil, errors.New(aerr.Error())
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
		return nil, errors.Errorf("read params malformed: %v", err)
	}

	result, err := s3w.reader(readParams)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return nil, errors.New(fmt.Sprint(s3.ErrCodeNoSuchBucket, aerr.Error()))

			default:
				return nil, aerr
			}
		} else {
			return nil, errors.New(aerr.Error())
		}
	}
	return result, err
}

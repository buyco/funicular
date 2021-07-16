package client_test

import (
	"github.com/buyco/funicular/internal/mock"
	. "github.com/buyco/funicular/pkg/client"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/awstesting/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"strings"
)

var _ = Describe("AWS", func() {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	var awsConfig = &aws.Config{
		DisableSSL: aws.Bool(true),
		Endpoint:   aws.String(server.URL),
	}
	var awsSession = NewAWSSession(awsConfig)

	Describe("Using AWS Manager", func() {

		var manager *AWSManager

		BeforeEach(func() {
			manager = NewAWSManager(awsSession)
		})

		Context("From constructor function", func() {

			It("creates a valid instance", func() {
				Expect(manager).To(BeAssignableToTypeOf(&AWSManager{}))
			})

			It("contains same S3 client", func() {
				Expect(manager.S3Manager).To(BeAssignableToTypeOf(&S3Manager{}))
			})

			It("has no S3 clients", func() {
				Expect(manager.S3Manager.GetAll()).To(HaveLen(0))
			})
		})
	})

	Describe("Using AWS S3 Manager", func() {
		var s3Manager *S3Manager

		BeforeEach(func() {
			s3Manager = NewS3Manager(mock.Session)
		})

		Context("From constructor function", func() {

			It("creates a valid instance", func() {
				Expect(s3Manager).To(BeAssignableToTypeOf(&S3Manager{}))
			})

			It("has a S3 client", func() {
				s3Wrapper := s3Manager.Add("test-bucket")
				Expect(s3Manager.GetAll()).To(HaveLen(1))
				Expect(s3Wrapper).To(BeAssignableToTypeOf(&S3Wrapper{}))
			})

			It("fetches a S3 client", func() {
				bucketName := "test-bucket"
				s3Wrapper := s3Manager.Add(bucketName)
				Expect(s3Manager.Get(bucketName)).To(BeAssignableToTypeOf(&S3Wrapper{}))
				Expect(s3Manager.Get(bucketName)).To(BeIdenticalTo(s3Wrapper))
			})

			It("deletes a bucket from manager storage", func() {
				bucketName := "test-bucket-del"
				s3Manager.Add(bucketName)
				err := s3Manager.Delete(bucketName)
				Expect(err).To(BeNil())
			})

			It("fails to delete not existing bucket", func() {
				bucketName := "test-bucket-fail"
				err := s3Manager.Delete(bucketName)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Using AWS S3 Wrapper", func() {

		var (
			mockCtrl *gomock.Controller
			mockS3   *mock_clients.MockStorageAccessLayer
		)

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			mockS3 = mock_clients.NewMockStorageAccessLayer(mockCtrl)
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})

		Context("From constructor function", func() {

			It("creates a valid instance", func() {
				s3Wrapper := NewS3Wrapper("test-bucket", NewS3Client(awsSession))
				Expect(s3Wrapper).To(BeAssignableToTypeOf(&S3Wrapper{}))
			})

			It("does not fail to call uploader", func() {
				mockS3.EXPECT().Upload("test-path", "test-file", strings.NewReader("test-data"))
				_, respErr := mockS3.Upload("test-path", "test-file", strings.NewReader("test-data"))
				Expect(respErr).ToNot(HaveOccurred())
			})

			It("does not fail to call downloader", func() {
				var buffer io.WriterAt
				mockS3.EXPECT().Download("test-path", "test-file", buffer)
				_, respErr := mockS3.Download("test-path", "test-file", buffer)
				Expect(respErr).ToNot(HaveOccurred())
			})

			It("does not fail to call reader", func() {
				mockS3.EXPECT().Read("test-path", int64(1), "test-file")
				_, respErr := mockS3.Read("test-path", int64(1), "test-file")
				Expect(respErr).ToNot(HaveOccurred())
			})
		})
	})
})

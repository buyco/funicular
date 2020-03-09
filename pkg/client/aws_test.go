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

var _ = Describe("Aws", func() {

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

			It("should create a valid instance", func() {
				Expect(manager).To(BeAssignableToTypeOf(&AWSManager{}))
			})

			It("should contain same S3 client", func() {
				Expect(manager.S3Manager).To(BeAssignableToTypeOf(&S3Manager{}))
			})

			It("should have no S3 clients", func() {
				Expect(manager.S3Manager.S3).To(HaveLen(0))
			})
		})
	})

	Describe("Using AWS S3 Manager", func() {

		var s3Manager = NewS3Manager(mock.Session)

		Context("From constructor function", func() {

			It("should create a valid instance", func() {
				Expect(s3Manager).To(BeAssignableToTypeOf(&S3Manager{}))
			})

			It("should have a S3 client", func() {
				s3Wrapper := s3Manager.Add("test-bucket")
				Expect(s3Manager.S3).To(HaveLen(1))
				Expect(s3Wrapper).To(BeAssignableToTypeOf(&S3Wrapper{}))
			})

			//It("should upload a file", func() {
			//	s3Wrapper := s3Manager.AddS3BucketManager("test-bucket")
			//	Expect(s3Wrapper).To(BeAssignableToTypeOf(&S3Wrapper{}))
			//
			//	_, upError := s3Wrapper.UploadFile(
			//		"",
			//		"foo.bar",
			//		strings.NewReader("foo:bar"),
			//	)
			//	fmt.Print(upError)
			//})
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

			It("should create a valid instance", func() {
				s3Wrapper := NewS3Wrapper("test-bucket", NewS3Client(awsSession))
				Expect(s3Wrapper).To(BeAssignableToTypeOf(&S3Wrapper{}))
			})

			It("should not fail to call uploader", func() {
				mockS3.EXPECT().Upload("test-path", "test-file", strings.NewReader("test-data"))
				_, respErr := mockS3.Upload("test-path", "test-file", strings.NewReader("test-data"))
				Expect(respErr).ToNot(HaveOccurred())
			})

			It("should not fail to call downloader", func() {
				var buffer io.WriterAt
				mockS3.EXPECT().Download("test-path", "test-file", buffer)
				_, respErr := mockS3.Download("test-path", "test-file", buffer)
				Expect(respErr).ToNot(HaveOccurred())
			})

			It("should not fail to call reader", func() {
				mockS3.EXPECT().Read("test-path", int64(1), "test-file")
				_, respErr := mockS3.Read("test-path", int64(1), "test-file")
				Expect(respErr).ToNot(HaveOccurred())
			})
		})
	})
})

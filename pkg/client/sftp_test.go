package client_test

import (
	. "github.com/buyco/funicular/pkg/client"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sftp", func() {

	var manager *SFTPManager
	config := &ssh.ClientConfig{
		User:            "foo",
		Auth:            []ssh.AuthMethod{ssh.Password("bar")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         2 * time.Second,
	}

	Describe("Using Manager", func() {

		manager = NewSFTPManager("localhost", 22, config, 1, logrus.New())

		Context("From constructor function", func() {

			It("should create a valid instance", func() {
				Expect(manager).To(BeAssignableToTypeOf(&SFTPManager{}))
			})

			It("should contain zero clients", func() {
				sftpCli, err := manager.GetClient()
				Expect(sftpCli).To(BeNil())
				Expect(err).To(HaveOccurred())
			})

			It("should not fail to close without clients", func() {
				err := manager.Close()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		It("should fail to add new client", func() {
			addCliErr := manager.AddClient()
			Expect(addCliErr).To(HaveOccurred())
		})
	})
})

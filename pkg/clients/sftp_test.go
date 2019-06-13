package clients_test

import (
	. "github.com/buyco/funicular/pkg/clients"
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

		manager = NewSFTPManager("localhost", 22, config)

		Context("From constructor function", func() {

			It("should create a valid instance", func() {
				Expect(manager).To(BeAssignableToTypeOf(&SFTPManager{}))
			})

			It("should contain zero clients", func() {
				Expect(manager.Conns).To(HaveLen(0))
			})
		})

		It("should fail to add new client", func() {
			sftpWrapper, addCliErr := manager.AddClient()
			Expect(addCliErr).To(HaveOccurred())
			Expect(sftpWrapper).To(BeNil())
		})
	})
})

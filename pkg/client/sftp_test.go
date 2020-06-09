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

			It("creates a valid instance", func() {
				Expect(manager).To(BeAssignableToTypeOf(&SFTPManager{}))
			})

			It("contains zero clients", func() {
				sftpCli, err := manager.GetClient()
				Expect(sftpCli).To(BeNil())
				Expect(err).To(HaveOccurred())
			})

			It("does not fail to close without clients", func() {
				err := manager.Close()
				Expect(err).ToNot(HaveOccurred())
			})

			It("adds a Factory to pool", func() {
				manager.SetPoolFactory(func() interface{} { return &SFTPWrapper{}})
				Expect(manager.GetClient()).To(BeAssignableToTypeOf(&SFTPWrapper{}))
			})

		})

		It("fails to add new client", func() {
			addCliErr := manager.AddClient()
			Expect(addCliErr).To(HaveOccurred())
		})

		It("fails to get a client", func() {
			client, getCliErr := manager.GetClient()
			Expect(client).To(BeNil())
			Expect(getCliErr).To(HaveOccurred())
		})
	})
})

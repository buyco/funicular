package client_test

import (
	. "github.com/buyco/funicular/pkg/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("AMQP", func() {

	Describe("Using Manager", func() {
		var manager *AMQPManager
		config := NewAMQPManagerConfig("localhost", 5672, "guest", "guest", nil)
		BeforeEach(func() {
			var err error
			manager, err = NewAMQPManager(config, 22, logrus.New())
			Expect(err).ToNot(HaveOccurred())
		})

		Context("From constructor function", func() {

			It("creates a valid instance", func() {
				Expect(manager).To(BeAssignableToTypeOf(&AMQPManager{}))
			})

			It("contains zero clients", func() {
				cli, err := manager.GetClient()
				Expect(cli).To(BeNil())
				Expect(err).To(HaveOccurred())
			})

			It("does not fail to close without clients", func() {
				err := manager.Close()
				Expect(err).ToNot(HaveOccurred())
			})

			It("adds a Factory to pool", func() {
				manager.SetPoolFactory(func() interface{} { return &AMQPWrapper{} })
				Expect(manager.GetClient()).To(BeAssignableToTypeOf(&AMQPWrapper{}))
			})

		})

		It("adds a new client", func() {
			addCliErr := manager.AddClient()
			Expect(addCliErr).ToNot(HaveOccurred())
		})

		It("puts a client", func() {
			manager.PutClient(&AMQPWrapper{})
			client, getCliErr := manager.GetClient()
			Expect(client).ToNot(BeNil())
			Expect(getCliErr).ToNot(HaveOccurred())
		})

		It("fails to get a client", func() {
			client, getCliErr := manager.GetClient()
			Expect(client).To(BeNil())
			Expect(getCliErr).To(HaveOccurred())
		})
	})
})

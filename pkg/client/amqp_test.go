package client_test

import (
	. "github.com/buyco/funicular/pkg/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AMQP", func() {
	Describe("Using Manager", func() {
		var manager *AMQPConnection
		config := NewAMQPConnectionConfig("localhost", 5672, "guest", "guest", nil)
		BeforeEach(func() {
			var err error
			manager, err = NewAMQPConnection(config)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("From constructor function", func() {

			It("creates a valid instance", func() {
				Expect(manager).To(BeAssignableToTypeOf(&AMQPConnection{}))
			})

			It("returns a new channel", func() {
				Expect(manager.Channel()).To(BeAssignableToTypeOf(&AMQPChannel{}))
			})

		})
	})
})

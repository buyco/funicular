package client

import (
	"crypto/tls"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	amqp "github.com/rabbitmq/amqp091-go"
	"net"
)

var _ = Describe("AMQP", func() {
	Describe("Using Manager", func() {
		var manager AMQPConnection
		config := NewAMQPConnectionConfig("localhost", 5672, "guest", "guest", nil)
		BeforeEach(func() {
			var err error
			manager, err = NewAMQPConnection(config)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("From constructor function", func() {

			It("creates a valid instance", func() {
				Expect(manager).To(BeAssignableToTypeOf(&amqpConnection{}))
			})

			It("returns a new channel", func() {
				Expect(manager.Channel()).To(BeAssignableToTypeOf(&amqpChannel{}))
			})

		})

		Describe("LocalAddr", func() {
			It("returns value", func() {
				Expect(manager.LocalAddr()).To(BeAssignableToTypeOf(&net.TCPAddr{}))
			})
		})

		Describe("ConnectionState", func() {
			It("returns value", func() {
				Expect(manager.ConnectionState()).To(BeAssignableToTypeOf(tls.ConnectionState{}))
			})
		})

		Describe("IsClosed", func() {
			It("returns value", func() {
				Expect(manager.IsClosed()).To(BeFalse())
			})
		})

		Describe("Close", func() {
			It("returns value", func() {
				Expect(manager.Close()).ToNot(HaveOccurred())
			})
		})

		Describe("NotifyClose", func() {
			It("returns value on close connection", func() {
				receiver := make(chan *amqp.Error)
				manager.NotifyClose(receiver)
				Expect(manager.Close()).ToNot(HaveOccurred())

				Expect(<-receiver).To(BeAssignableToTypeOf(&amqp.Error{}))
			})
		})

		Describe("Channel", func() {
			var (
				channel AMQPChannel
				err     error
			)
			BeforeEach(func() {
				channel, err = manager.Channel()
				Expect(err).ToNot(HaveOccurred())
				Expect(channel).To(BeAssignableToTypeOf(&amqpChannel{}))
			})

			Describe("IsClosed", func() {
				It("returns value", func() {
					Expect(channel.IsClosed()).To(BeFalse())
				})
			})

			Describe("Close", func() {
				It("returns value", func() {
					Expect(channel.Close()).ToNot(HaveOccurred())
				})
			})

			Describe("NotifyClose", func() {
				It("returns value on close connection", func() {
					receiver := make(chan *amqp.Error)
					channel.NotifyClose(receiver)
					Expect(channel.Close()).ToNot(HaveOccurred())

					Expect(<-receiver).To(BeAssignableToTypeOf(&amqp.Error{}))
				})
			})
		})
	})
})

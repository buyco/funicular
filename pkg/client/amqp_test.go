package client_test

import (
	. "github.com/buyco/funicular/pkg/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

var _ = Describe("AMQP", func() {

	Describe("Using Manager", func() {

		Context("From constructor function", func() {
			It("fails to connect to RabbitMQ", func() {
				config := NewAMQPManagerConfig("localhost:5672", "guest", "guest", nil)
				manager, err := NewAMQPManager(config, 22, logrus.New())
				Expect(manager).To(BeNil())
				Expect(err).To(HaveOccurred())
			})

		})

		It("creates an AMQPConfig", func() {
			config := NewAMQPConfig("/", 12, 2*time.Millisecond)
			Expect(config).To(BeAssignableToTypeOf(amqp.Config{}))
		})

		It("creates an AMQPManagerConfig", func() {
			config := NewAMQPManagerConfig("localhost:5672", "guest", "guest", nil)
			Expect(config).To(BeAssignableToTypeOf(&AMQPManagerConfig{}))
		})
	})
})

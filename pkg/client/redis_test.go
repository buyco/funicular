package client_test

import (
	. "github.com/buyco/funicular/pkg/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"strconv"
)

var _ = Describe("Redis", func() {
	// Declaring var for tests
	port, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	var config = RedisConfig{
		Host: os.Getenv("REDIS_HOST"),
		Port: uint16(port),
		DB:   uint8(db),
	}

	Describe("Using Manager", func() {

		var category = "test"
		var manager *RedisManager

		BeforeEach(func() {
			manager = NewRedisManager(config)
		})

		Context("From constructor function", func() {

			It("should contain zero clients", func() {
				Expect(len(manager.Clients)).To(BeZero())
			})
		})

		It("should fail to add client with empty category", func() {
			client, err := manager.AddClient("")
			Expect(err).To(HaveOccurred())
			Expect(client).To(BeNil())
		})

		Context("Without Redis client in the stack", func() {

			It("should use category as channel if channel is empty and add client to manager", func() {
				client, err := manager.AddClient(category)
				Expect(err).ToNot(HaveOccurred())
				Expect(manager.Clients[category]).To(Equal(client))
			})

			It("should fail to close", func() {
				Expect(manager.Close()).To(HaveOccurred())
			})
		})

		Context("With Redis clients in the stack", func() {

			It("should use category as channel if channel is empty and add client to manager", func() {
				client, err := manager.AddClient(category)
				Expect(err).ToNot(HaveOccurred())
				_, err2 := manager.AddClient(category)
				Expect(err2).ToNot(HaveOccurred())
				Expect(manager.Clients[category]).To(Equal(client))
			})

			It("should close all clients", func() {
				_, _ = manager.AddClient(category)
				_, _ = manager.AddClient(category)
				var err error
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})

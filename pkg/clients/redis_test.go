package clients_test

import (
	. "github.com/buyco/funicular/pkg/clients"
	"github.com/buyco/funicular/pkg/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"strconv"
	"time"
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
	var wrapper, nilErr = NewRedisWrapper(config, "test-channel", "")

	Describe("Using Manager", func() {

		var manualManager = &RedisManager{
			Clients: make(map[string][]*RedisWrapper),
		}
		var category = "test"
		var manager *RedisManager

		BeforeEach(func() {
			manager = NewRedisManager()
		})

		Context("From constructor function", func() {

			It("should create a valid instance", func() {
				Expect(manager).To(Equal(manualManager))
			})

			It("should contain zero clients", func() {
				Expect(len(manager.Clients)).To(BeZero())
			})
		})

		It("should fail to add client with empty category", func() {
			client, err := manager.AddClient(config, "", "test", "")
			Expect(err).To(HaveOccurred())
			Expect(client).To(BeNil())
		})

		Context("Without Redis client in the stack", func() {

			It("should use category as channel if channel is empty and add client to manager", func() {
				client, err := manager.AddClient(config, category, "", "")
				Expect(err).ToNot(HaveOccurred())
				Expect(client.GetChannel()).To(Equal(category))
				Expect(manager.Clients[category]).To(HaveLen(1))
				Expect(manager.Clients[category][0]).To(Equal(client))
			})

			It("should not close clients", func() {
				var err error
				stdout := utils.CaptureStdout(func() { err = manager.Close() })
				Expect(err).ToNot(HaveOccurred())
				Expect(stdout).To(ContainSubstring("Manager have no clients to close..."))
			})
		})

		Context("With Redis clients in the stack", func() {

			It("should use category as channel if channel is empty and add client to manager", func() {
				client, err := manager.AddClient(config, category, "", "")
				client2, err2 := manager.AddClient(config, category, "", "")
				Expect(err).ToNot(HaveOccurred())
				Expect(client.GetChannel()).To(Equal(category))
				Expect(err2).ToNot(HaveOccurred())
				Expect(client2.GetChannel()).To(Equal(category))
				Expect(manager.Clients[category]).To(HaveLen(2))
				Expect(manager.Clients[category][0]).To(Equal(client))
				Expect(manager.Clients[category][1]).To(Equal(client2))
			})

			It("should close all clients", func() {
				_, _ = manager.AddClient(config, category, "", "")
				_, _ = manager.AddClient(config, category, "", "")
				var err error
				stdout := utils.CaptureStdout(func() { err = manager.Close() })
				Expect(err).ToNot(HaveOccurred())
				Expect(stdout).ToNot(ContainSubstring("Manager have no clients to close..."))
			})
		})
	})

	Describe("Using Wrapper", func() {

		var group = "foo-group"
		var validDefaultMsgId = "1538561700640-0"
		var malformedMsgId = "foo:bar"

		Context("From constructor function", func() {

			It("should create a valid instance", func() {
				Expect(nilErr).ToNot(HaveOccurred())
				Expect(wrapper.GetChannel()).To(Equal("test-channel"))
			})

			It("should fail with empty string for channel", func() {
				_, filledErr := NewRedisWrapper(config, "", "")
				Expect(filledErr).To(
					SatisfyAll(
						HaveOccurred(),
						MatchError("channel must be filled"),
					),
				)
			})
		})

		Context("When Redis stream channel is empty", func() {

			It("should fail to read message", func() {
				_, readErr := wrapper.ReadMessage("$", 1, 100*time.Millisecond)
				Expect(readErr).To(
					SatisfyAll(
						HaveOccurred(),
						MatchError("redis: nil"),
					),
				)

				_, readErr = wrapper.ReadRangeMessage("-", "+")
				Expect(readErr).ToNot(HaveOccurred())
			})

			It("should not have messages to delete", func() {
				id, readErr := wrapper.DeleteMessage(validDefaultMsgId)
				Expect(id).To(BeZero())
				Expect(readErr).ToNot(HaveOccurred())
			})

			It("should flush DB", func() {
				response, flushErr := wrapper.FlushDB()
				Expect(flushErr).ToNot(HaveOccurred())
				Expect(response).To(Equal("OK"))

				response, flushErr = wrapper.FlushDBAsync()
				Expect(flushErr).ToNot(HaveOccurred())
				Expect(response).To(Equal("OK"))
			})

			It("should flush all", func() {
				response, flushErr := wrapper.FlushAll()
				Expect(flushErr).ToNot(HaveOccurred())
				Expect(response).To(Equal("OK"))

				response, flushErr = wrapper.FlushDBAsync()
				Expect(flushErr).ToNot(HaveOccurred())
				Expect(response).To(Equal("OK"))
			})

			Context("When no group exists", func() {

				It("should not have messages to acknowledge", func() {
					id, readErr := wrapper.AckMessage(group, validDefaultMsgId)
					Expect(id).To(BeZero())
					Expect(readErr).ToNot(HaveOccurred())
				})

				It("should not have pending messages", func() {
					pendResp, readErr := wrapper.PendingMessage(group)
					Expect(pendResp).To(BeNil())
					Expect(readErr).To(HaveOccurred())
				})
			})

			Context("When a group exists", func() {

				BeforeEach(func() {
					cliAddGrpResponse, errAddGrp := wrapper.CreateGroup(group, "$")
					Expect(cliAddGrpResponse).To(Equal("OK"))
					Expect(errAddGrp).ToNot(HaveOccurred())
				})

				AfterEach(func() {
					_, flushErr := wrapper.FlushAll()
					Expect(flushErr).ToNot(HaveOccurred())
				})

				It("should fail to create same group", func() {
					failResp, errSameAddGrp := wrapper.CreateGroup(group, "$")
					Expect(failResp).To(BeEmpty())
					Expect(errSameAddGrp).To(
						SatisfyAll(
							HaveOccurred(),
							MatchError("BUSYGROUP Consumer Group name already exists"),
						),
					)
				})

				It("should fail to acknowledge malformed message ID", func() {
					ackMsgGrp, errAckMgsGrp := wrapper.AckMessage(group, malformedMsgId)
					Expect(ackMsgGrp).To(BeZero())
					Expect(errAckMgsGrp).To(
						SatisfyAll(
							HaveOccurred(),
							MatchError("ERR Invalid stream ID specified as stream command argument"),
						),
					)
				})

				It("should not have message to acknowledge", func() {
					ackMsgGrp, errAckMgsGrp := wrapper.AckMessage(group, validDefaultMsgId)
					Expect(ackMsgGrp).To(BeZero())
					Expect(errAckMgsGrp).ToNot(HaveOccurred())
				})
			})
		})

		Context("When Redis stream channel is filled", func() {

			message := map[string]interface{}{"foo": "bar"}
			var errAddMsg error

			BeforeEach(func() {
				_, errAddMsg = wrapper.AddMessage(message)
				Expect(errAddMsg).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				_, flushErr := wrapper.FlushAll()
				Expect(flushErr).ToNot(HaveOccurred())
			})

			It("should read message", func() {
				msg, readErr := wrapper.ReadMessage("0", 1, 10*time.Millisecond)
				Expect(readErr).ToNot(HaveOccurred())
				Expect(msg).To(HaveLen(1))
			})

			Context("When a group exists", func() {

				BeforeEach(func() {
					createGrp, errCreateGrp := wrapper.CreateGroup(group, "$")
					Expect(createGrp).To(Equal("OK"))
					Expect(errCreateGrp).ToNot(HaveOccurred())

					_, errAddMsg = wrapper.AddMessage(message)
					Expect(errAddMsg).ToNot(HaveOccurred())
				})

				AfterEach(func() {
					_, flushErr := wrapper.FlushAll()
					Expect(flushErr).ToNot(HaveOccurred())
				})

				It("should read group message", func() {
					msgs, errReadGrp := wrapper.ReadGroupMessage(group, 1, 100*time.Millisecond)
					Expect(msgs).To(HaveLen(1))
					Expect(errReadGrp).ToNot(HaveOccurred())
				})

			})
		})
	})
})

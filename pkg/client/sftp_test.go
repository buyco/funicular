package client_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"

	. "github.com/buyco/funicular/pkg/client"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SFTP", func() {
	config := &ssh.ClientConfig{
		User:            "foo",
		Auth:            []ssh.AuthMethod{ssh.Password("bar")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         2 * time.Second,
	}

	netPipe := func(port string) {
		// An SSH server is represented by a ServerConfig, which holds
		// certificate details and handles authentication of ServerConns.
		config := &ssh.ServerConfig{
			NoClientAuth: true,
		}

		// Once a ServerConfig has been configured, connections can be
		// accepted.
		listener, err := net.Listen("tcp", "127.0.0.1:"+port)
		if err != nil {
			log.Fatal("failed to listen for connection", err)
		}
		fmt.Printf("Listening on %v\n", listener.Addr())

		nConn, err := listener.Accept()
		if err != nil {
			log.Fatal("failed to accept incoming connection", err)
		}

		privateBytes, err := ioutil.ReadFile("fixture/id_rsa_sftp")
		if err != nil {
			log.Fatal("Failed to load private key", err)
		}

		private, err := ssh.ParsePrivateKey(privateBytes)
		if err != nil {
			log.Fatal("Failed to parse private key", err)
		}

		config.AddHostKey(private)

		// Before use, a handshake must be performed on the incoming
		// net.Conn.
		_, chans, reqs, err := ssh.NewServerConn(nConn, config)
		if err != nil {
			log.Fatal("failed to handshake", err)
		}

		// The incoming Request channel must be serviced.
		go ssh.DiscardRequests(reqs)

		// Service the incoming Channel channel.
		for newChannel := range chans {
			// Channels have a type, depending on the application level
			// protocol intended. In the case of an SFTP session, this is "subsystem"
			// with a payload string of "<length=4>sftp"
			if newChannel.ChannelType() != "session" {
				newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
				continue
			}
			channel, requests, err := newChannel.Accept()
			if err != nil {
				log.Fatal("could not accept channel.", err)
			}
			// Sessions have out-of-band requests such as "shell",
			// "pty-req" and "env".  Here we handle only the
			// "subsystem" request.
			go func(in <-chan *ssh.Request) {
				for req := range in {
					ok := false
					switch req.Type {
					case "subsystem":
						if string(req.Payload[4:]) == "sftp" {
							ok = true
						}
					}
					req.Reply(ok, nil)
				}
			}(requests)

			server, err := sftp.NewServer(
				channel,
			)
			if err != nil {
				log.Fatal(err)
			}
			if err := server.Serve(); err == io.EOF {
				server.Close()
				log.Fatal("sftp client exited session.")
			} else if err != nil {
				log.Fatal("sftp server completed with error:", err)
			}
		}
	}

	Describe("Using Manager", func() {
		go netPipe("18000")
		Context("From constructor function", func() {
			var badManager *SFTPManager
			BeforeEach(func() {
				badManager = NewSFTPManager("127.0.0.1", 0, config, 1)
			})
			AfterEach(func() {
				err := badManager.Close()
				Expect(err).ToNot(HaveOccurred())
			})

			It("creates a valid instance", func() {
				Expect(badManager).To(BeAssignableToTypeOf(&SFTPManager{}))
			})

			It("contains zero clients", func() {
				sftpCli, err := badManager.GetClient()
				Expect(sftpCli).To(BeNil())
				Expect(err).To(HaveOccurred())
			})

		})

		Context("From bad address", func() {
			manager := NewSFTPManager("127.0.0.1", 0, config, 1)
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

		Context("From active address", func() {
			manager := NewSFTPManager("127.0.0.1", 18000, config, 1)

			It("gets a client", func() {
				addCliErr := manager.AddClient()
				Expect(addCliErr).ToNot(HaveOccurred())
				client, getCliErr := manager.GetClient()
				Expect(client).To(BeAssignableToTypeOf(&SFTPWrapper{}))
				Expect(getCliErr).ToNot(HaveOccurred())
			})

			It("says when client is closed", func() {
				addCliErr := manager.AddClient()
				client, _ := manager.GetClient()
				Expect(client.Close()).To(BeTrue())
			})

			It("adds a Factory to pool", func() {
				manager.SetPoolFactory(func() interface{} { return &SFTPWrapper{} })
				Expect(manager.GetClient()).To(BeAssignableToTypeOf(&SFTPWrapper{}))
			})
		})
	})
})

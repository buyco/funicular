module github.com/buyco/funicular

go 1.16

require (
	github.com/aws/aws-sdk-go v1.37.23
	github.com/buyco/keel v0.5.2
	github.com/go-redis/redis/v7 v7.4.0
	github.com/golang/mock v1.5.0
	github.com/jinzhu/copier v0.2.5
	github.com/onsi/ginkgo v1.15.0
	github.com/onsi/gomega v1.10.5
	github.com/pkg/sftp v1.12.0
	github.com/sirupsen/logrus v1.8.0
	github.com/streadway/amqp v1.0.0
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
	gopkg.in/eapache/go-resiliency.v1 v1.2.0
)

// This version is in fact v0.11.0 tag
retract v1.11.0
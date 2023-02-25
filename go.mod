module github.com/buyco/funicular

go 1.17

// This version is in fact v0.11.0 tag
retract v1.11.0

require (
	github.com/aws/aws-sdk-go v1.44.58
	github.com/go-redis/redis/v7 v7.4.1
	github.com/golang/mock v1.6.0
	github.com/jinzhu/copier v0.3.5
	github.com/joho/godotenv v1.4.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.19.0
	github.com/pkg/sftp v1.13.5
	github.com/rabbitmq/amqp091-go v1.4.0
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f
	gopkg.in/eapache/go-resiliency.v1 v1.2.0
)

require (
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

module github.com/buyco/funicular

go 1.17

require (
	github.com/aws/aws-sdk-go v1.37.23
	github.com/go-redis/redis/v7 v7.4.0
	github.com/golang/mock v1.5.0
	github.com/jinzhu/copier v0.2.5
	github.com/joho/godotenv v1.3.0
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.16.0
	github.com/pkg/sftp v1.12.0
	github.com/rabbitmq/amqp091-go v1.3.4
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	gopkg.in/eapache/go-resiliency.v1 v1.2.0
)

require (
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.0.0-20210428140749-89ef3d95e781 // indirect
	golang.org/x/sys v0.0.0-20211019181941-9d821ace8654 // indirect
	golang.org/x/text v0.3.6 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

// This version is in fact v0.11.0 tag
retract v1.11.0

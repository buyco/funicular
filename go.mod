module github.com/buyco/funicular

go 1.17

require (
	github.com/aws/aws-sdk-go v1.37.23
	github.com/buyco/keel v0.5.2
	github.com/go-redis/redis/v7 v7.4.0
	github.com/golang/mock v1.5.0
	github.com/jinzhu/copier v0.2.5
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.16.0
	github.com/pkg/sftp v1.12.0
	github.com/sirupsen/logrus v1.8.0
	github.com/streadway/amqp v1.0.0
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
	gopkg.in/eapache/go-resiliency.v1 v1.2.0
)

require (
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/joho/godotenv v1.3.0 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/magefile/mage v1.10.0 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/thoas/go-funk v0.5.0 // indirect
	golang.org/x/net v0.0.0-20210428140749-89ef3d95e781 // indirect
	golang.org/x/sys v0.0.0-20210423082822-04245dca01da // indirect
	golang.org/x/text v0.3.6 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

// This version is in fact v0.11.0 tag
retract v1.11.0

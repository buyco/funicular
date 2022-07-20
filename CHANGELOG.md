CHANGE LOG
===================

# [v1.15.2](https://github.com/buyco/funicular/tree/v1.15.2)

* Upgrade dependencies (#31)

# [v1.15.1](https://github.com/buyco/funicular/tree/v1.15.1)

* Fix go vet error on CI

# [v1.15.0](https://github.com/buyco/funicular/tree/v1.15.0)

* Change AMQP client interface to make it mockable (#30)

# [v1.14.2](https://github.com/buyco/funicular/tree/v1.14.2)

* Fix CI and edit staticcheck checks (#28)
* Move AMQP client from streadway/amqp to rabbitmq/amqp091-go (#26)
* Migrate errors package to xerrors

# [v1.14.1](https://github.com/buyco/funicular/tree/v1.14.1)

* Delete logger on common usage and create a tag on debug build (#25)
* Update README on Makefile commands
* Run linter and fix code (#24)
* Edit Makefile, README and fix staticcheck command

# [v1.14.0](https://github.com/buyco/funicular/tree/v1.14.0)

* Add Makefile + linter + vet + bump to Go 1.17 (#23)

# [v1.13.0](https://github.com/buyco/funicular/tree/v1.13.0)

* Fix mutexes usage and add AWS S3 functions (#22)

# [v1.12.0](https://github.com/buyco/funicular/tree/v1.12.0)

* Fix race condition SFTP (#21)
* Move CI to Github actions (#19)

# [v1.11.2](https://github.com/buyco/funicular/tree/v1.11.2)

* Make the whole reconnection process lock actions on SFTP client (#18)

# [v1.11.1](https://github.com/buyco/funicular/tree/v1.11.1)

* Update README and tag over retracted version

# [v1.0.0](https://github.com/buyco/funicular/tree/v1.0.0)

* Retract a bad version of Funicular (#17)

# [v0.11.0](https://github.com/buyco/funicular/tree/v0.11.0)

* Upgrade libs & bump Go version to 1.16 (#16)
* Delete useless logger, add and fix tests (#14)

# [v0.10.2](https://github.com/buyco/funicular/tree/v0.10.2)

* Fix infinite loop in reconnect functions

# [v0.10.1](https://github.com/buyco/funicular/tree/v0.10.1)

* Fix errors in tests

# [v0.10.0](https://github.com/buyco/funicular/tree/v0.10.0)

* Change logger strategy

# [v0.9.0](https://github.com/buyco/funicular/tree/v0.9.0)

* Make connection / channel public to handle them manually (#13)

# [v0.8.0](https://github.com/buyco/funicular/tree/v0.8.0)

* Move AMQP config from address to host + port (#12)

# [v0.7.0](https://github.com/buyco/funicular/tree/v0.7.0)

* Add option structures for S3 uploader and downloader (#11)
* Reformat code

# [v0.6.2](https://github.com/buyco/funicular/tree/v0.6.2)

* Change pool Put() error in AddClient() as log

# [v0.6.2](https://github.com/buyco/funicular/tree/v0.6.2)

* Change pool Put() error in AddClient() as log

# [v0.6.1](https://github.com/buyco/funicular/tree/v0.6.1)

* Handle error returned in AddClient() of SFTP and AMQP managers (#10)

# [v0.6.0](https://github.com/buyco/funicular/tree/v0.6.0)

* Add AMQP client manager (#9)

# [v0.5.0](https://github.com/buyco/funicular/tree/v0.5.0)

* Replace sync.Pool to create our own Pool of objects (#8)

# [v0.4.1](https://github.com/buyco/funicular/tree/v0.4.1)

* Fix go-redis pkg dependency in go.mod

# [v0.4.0](https://github.com/buyco/funicular/tree/v0.4.0)

* Fix S3 manager and delete RedisWrapper to use redis client instead (#6)

# [v0.3.1](https://github.com/buyco/funicular/tree/v0.3.1)

* Add a PutClient function to SFTP Manager

# [v0.3.0](https://github.com/buyco/funicular/tree/v0.3.0)

* Move to sync.Pool for SFTP Clients

# [v0.2.0](https://github.com/buyco/funicular/tree/v0.2.0)

* Update deps and refactoring packages name

# [v0.1.15](https://github.com/buyco/funicular/tree/v0.1.15)

* Update Keel to v0.1.0
* Add comments on code to satisfy golint (we need more documentation)
* Change licence, run gofmt and linter on project

# [v0.1.14](https://github.com/buyco/funicular/tree/v0.1.14)

* Update Keel dep to v0.0.4

# [v0.1.13](https://github.com/buyco/funicular/tree/v0.1.13)

* Migrate common packages to Keel toolkit
* Move coverage formatter script to tools directory

# [v0.1.12](https://github.com/buyco/funicular/tree/v0.1.12)

* Delete "fake" copy of Redis clients (slice are pointers)

# [v0.1.11](https://github.com/buyco/funicular/tree/v0.1.11)

* Change AWS S3 deleter

# [v0.1.10](https://github.com/buyco/funicular/tree/v0.1.10)

* Add Deleter for AWS S3

# [v0.1.9](https://github.com/buyco/funicular/tree/v0.1.9)

* Fix AWS S3 Reader

# [v0.1.8](https://github.com/buyco/funicular/tree/v0.1.8)

* Add S3 reader and test

# [v0.1.7](https://github.com/buyco/funicular/tree/v0.1.7)

* Add Changelog
* Add S3 Downloader

# [v0.1.6](https://github.com/buyco/funicular/tree/v0.1.6)

* Delete log when SFTP circuit breaker is open

# [v0.1.5](https://github.com/buyco/funicular/tree/v0.1.5)

* Delete log when SFTP circuit breaker is open
* Add Circuit Breaker on SFTP wrapper and add logrus lib

# [v0.1.4](https://github.com/buyco/funicular/tree/v0.1.4)

* Add mutex for concurrency

# [v0.1.3](https://github.com/buyco/funicular/tree/v0.1.3)

* Move redis config to manager construct

# [v0.1.2](https://github.com/buyco/funicular/tree/v0.1.2)

* Edit SFTP : add close function on manager Edit Redis : add reconnect public function, add closed property on wrapper
* Add coverage option on tests

# [v0.1.1](https://github.com/buyco/funicular/tree/v0.1.1)

* Finally utils must be in internal

# [v0.1.0](https://github.com/buyco/funicular/tree/v0.1.0)

* Move utils from internal/ to pkg/
* Transfer repo to Buyco
* Initial commit

package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/buyco/funicular/pkg/client"
	"github.com/buyco/keel/pkg/helper"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const stream = "example-stream"
const category = stream + "-cat"
const bucketName = "buyco-foo-bar"
const storePath = "/foo/bar/"

func main() {
	helper.LoadEnvFile(os.Getenv("ENV"))
	fileChan := make(chan redis.XMessage)
	s3Chan := make(chan string)
	go func() {
		redisPort, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
		redisDb, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
		redisManager := client.NewRedisManager(
			client.RedisConfig{
				Host: os.Getenv("REDIS_HOST"),
				Port: uint16(redisPort),
				DB:   uint8(redisDb),
			},
			logrus.New(),
		)
		redisManager.AddClient(category)
		redisCli := redisManager.Clients[category]

		defer func() {
			err := redisCli.Close()
			if err != nil {
				log.Fatalf("Failed to close redis client: %v", err)
			}
		}()

		go func() {
			for {
				select {
				case filename := <-s3Chan:
					_, err := redisCli.XDel(stream, filename).Result()
					if err != nil {
						log.Fatalf("Failed to delete stream message: %v", err)
					}
					log.Printf("File message stream deleted for ID: %s", filename)
				}
			}
		}()
		lastID := "$"
		for {
			rArgs := &redis.XReadArgs{
				Streams: []string{stream, lastID},
				Count:   5,
				Block:   3000 * time.Millisecond,
			}
			vals, err := redisCli.XRead(rArgs).Result()
			if err != nil {
				log.Printf("Redis read error: %v", err)
			} else {
				NbStream := len(vals)
				NbMsgLastStreamEntry := len(vals[NbStream-1].Messages)
				lastID = vals[NbStream-1].Messages[NbMsgLastStreamEntry-1].ID
				for _, msgs := range vals {
					for _, msg := range msgs.Messages {
						log.Printf("Got message with file: %s", msg.Values["filename"].(string))
						fileChan <- msg
					}
				}
			}
		}
	}()
	var awsConfig = &aws.Config{
		MaxRetries: aws.Int(2),
	}
	awsManager := client.NewAWSManager(client.NewAWSSession(awsConfig))
	s3Bucket := awsManager.S3Manager.Add(bucketName)
	for {
		select {
		case fileData := <-fileChan:
			result, err := s3Bucket.Upload(
				storePath,
				fileData.Values["filename"].(string),
				strings.NewReader(fileData.Values["fileData"].(string)),
			)
			if err != nil {
				log.Printf("Failed to upload file, %v", err)
			} else {
				log.Printf("File uploaded to, %s\n", aws.StringValue(&result))
				s3Chan <- fileData.ID
			}
		}
	}
}

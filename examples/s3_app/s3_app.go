package main

import (
	"github.com/buyco/funicular/internal/utils"
	"github.com/buyco/funicular/pkg/clients"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-redis/redis"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const ENV_DIR = "../../.env"
const STREAM = "example-stream"
const CONSUMER_NAME = STREAM + "-consumer"
const BUCKET_NAME = "buyco-foo-bar"
const STORE_PATH = "/foo/bar/"

func main() {
	utils.LoadEnvFile(ENV_DIR, os.Getenv("ENV"))
	fileChan := make(chan redis.XMessage)
	s3Chan := make(chan string)
	go func() {
		redisPort, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
		redisDb, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
		redisCli, wrapperErr := clients.NewRedisWrapper(
			clients.RedisConfig{
				Host: os.Getenv("REDIS_HOST"),
				Port: uint16(redisPort),
				DB:   uint8(redisDb),
			},
			STREAM,
			CONSUMER_NAME,
		)
		if wrapperErr != nil {
			log.Fatalf("Redis read error: %v", wrapperErr)
		}

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
					_, err := redisCli.DeleteMessage(filename)
					if err != nil {
						log.Fatalf("Failed to delete stream message: %v", err)
					}
					log.Printf("File message stream deleted for ID: %s", filename)
				}
			}
		}()
		lastId := "$"
		for {
			vals, err := redisCli.ReadMessage(lastId, 5, 3000*time.Millisecond)
			if err != nil {
				log.Printf("Redis read error: %v", err)
			} else {
				NbStream := len(vals)
				NbMsgLastStreamEntry := len(vals[NbStream-1].Messages)
				lastId = vals[NbStream-1].Messages[NbMsgLastStreamEntry-1].ID
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
	awsManager := clients.NewAWSManager(clients.NewAWSSession(awsConfig))
	s3Bucket := awsManager.S3Manager.Add(BUCKET_NAME)
	for {
		select {
		case fileData := <-fileChan:
			result, err := s3Bucket.Upload(
				STORE_PATH,
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

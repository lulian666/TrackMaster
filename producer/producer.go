package main

import (
	"TrackMaster/initializer"
	"TrackMaster/model"
	"encoding/json"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"time"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.ConnectDB()
}

func NewProducer() (sarama.SyncProducer, error) {
	broker, ok := os.LookupEnv("BROKER")
	if !ok {
		broker = "localhost:9092"
	}
	brokerList := []string{broker}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokerList, config)

	if err != nil {
		return nil, err
	}
	return producer, nil
}

func ScanJobs() ([]model.Schedule, error) {
	var schedules []model.Schedule
	result := initializer.DB.Where("status = ?", true).Find(&schedules)
	return schedules, result.Error
}

func main() {
	producer, err := NewProducer()
	if err != nil {
		log.Fatalln("Failed to create producer:", err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln("Failed to close producer:", err)
		}
	}()

	topic, ok := os.LookupEnv("TOPIC")
	if !ok {
		topic = "pikachu.track"
	}

	ticker := time.Tick(5 * time.Minute)
	for range ticker {
		schedules, err := ScanJobs()
		if err != nil {
			log.Fatalln("Failed to read database:", err)
		}

		for _, schedule := range schedules {
			payload, err := json.Marshal(schedule)
			if err != nil {
				log.Fatalln("Failed to read marshal schedule:", err)
			}

			msg := &sarama.ProducerMessage{
				Topic:     topic,
				Value:     sarama.ByteEncoder(payload),
				Timestamp: time.Now(),
			}

			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				log.Printf("Failed to send message: %v\n", err)
			} else {
				log.Printf("Message sent to partition %d at offset %d\n", partition, offset)
			}
		}
	}
}

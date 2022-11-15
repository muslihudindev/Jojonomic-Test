package main

import (
	"context"
	"encoding/json"
	"fmt"
	"input-harga-storage-service/models"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"github.com/teris-io/shortid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	kafkaBrokerUrl := os.Getenv("KAFKA_BROKER_URL")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	kafkaGroupId := os.Getenv("KAFKA_GROUP_ID")
	postgresDsn := os.Getenv("POSTGRES_DSN")

	//gorm database
	connDb, err := gorm.Open(postgres.Open(postgresDsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to oppen connection to postgres:", err)
		panic(err)
	}

	ctx := context.Background()

	//kafka
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: strings.Split(kafkaBrokerUrl, ","),
		GroupID: kafkaGroupId,
		Topic:   kafkaTopic,
	})
	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			log.Println("could not fetch message : ", err.Error())
			break
		}

		fmt.Println("received: ", string(msg.Value))
		err = HandleInputHarga(connDb, msg.Value)
		if err != nil {
			log.Println("could not input to database : ", err.Error())
			continue
		}

		err = reader.CommitMessages(ctx, msg)
		if err != nil {
			log.Println("could not commit message : ", err.Error())
		}
	}

}

func HandleInputHarga(db *gorm.DB, msg []byte) error {
	message := models.Message{}

	err := json.Unmarshal(msg, &message)
	if err != nil {
		return err
	}

	reffID, err := shortid.Generate()
	if err != nil {
		return err
	}

	harga := models.Harga{
		ReffId:       reffID,
		AdminId:      message.AdminId,
		HargaTopup:   message.HargaBuyback,
		HargaBuyback: message.HargaBuyback,
	}

	err = db.Create(&harga).Error
	if err != nil {
		return err
	}
	return nil
}

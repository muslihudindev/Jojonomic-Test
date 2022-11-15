package main

import (
	"buyback-service/models"
	"buyback-service/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
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
	servicePort := os.Getenv("SERVICE_PORT")
	postgresDsn := os.Getenv("POSTGRES_DSN")

	//kafka
	connKafka, err := kafka.DialLeader(context.Background(), "tcp", kafkaBrokerUrl, kafkaTopic, 0)
	if err != nil {
		log.Fatal("failed to dial leader kafka:", err)
		panic(err)
	}
	defer connKafka.Close()

	//gorm database
	connDb, err := gorm.Open(postgres.Open(postgresDsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to oppen connection to postgres:", err)
		panic(err)
	}

	//router
	router := mux.NewRouter()
	router.HandleFunc("/api/buyback", HandleTopup(connKafka, connDb)).Methods(http.MethodPost)

	log.Printf("buyback-service run at : localhost:%s", servicePort)
	http.ListenAndServe(fmt.Sprintf(":%s", servicePort), router)

}

func HandleTopup(connKafka *kafka.Conn, db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reffID, err := shortid.Generate()
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, models.ResponseData{
				Error:   true,
				Message: err.Error(),
			})
			return
		}

		var req models.Params
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondError(w, http.StatusBadRequest, models.ResponseData{
				Error:   true,
				Message: err.Error(),
			})
			return
		}

		rekening := models.Rekening{}
		if err := db.Where("norek = ?", req.Norek).First(&rekening).Error; err != nil {
			utils.RespondError(w, http.StatusBadRequest, models.ResponseData{
				Error:   true,
				Message: err.Error(),
			})
			return
		}

		//validate
		if rekening.Saldo < req.Gram {
			utils.RespondError(w, http.StatusBadRequest, models.ResponseData{
				Error:   true,
				Message: "Saldo too low",
			})
			return

		}

		payloadBytes, err := json.Marshal(&req)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, models.ResponseData{
				Error:   true,
				Message: err.Error(),
			})
			return
		}

		connKafka.SetWriteDeadline(time.Now().Add(10 * time.Second))
		_, err = connKafka.WriteMessages(
			kafka.Message{
				Key:   []byte(reffID),
				Value: payloadBytes,
			},
		)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, models.ResponseData{
				Error:   true,
				ReffID:  reffID,
				Message: "Kafka not ready",
			})
			return
		}

		if err := connKafka.Close(); err != nil {
			utils.RespondError(w, http.StatusBadRequest, models.ResponseData{
				Error:   true,
				ReffID:  reffID,
				Message: "Kafka not ready",
			})
			return
		}

		utils.RespondJSON(w, http.StatusOK, models.ResponseData{
			Error:  false,
			ReffID: reffID,
		})
	}
}

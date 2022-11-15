package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"topup-service/models"
	"topup-service/utils"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"github.com/teris-io/shortid"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	kafkaBrokerUrl := os.Getenv("KAFKA_BROKER_URL")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	servicePort := os.Getenv("SERVICE_PORT")

	//kafka
	connKafka, err := kafka.DialLeader(context.Background(), "tcp", kafkaBrokerUrl, kafkaTopic, 0)
	if err != nil {
		log.Fatal("failed to dial leader kafka:", err)
		panic(err)
	}
	defer connKafka.Close()

	//router
	router := mux.NewRouter()
	router.HandleFunc("/api/topup", HandleTopup(connKafka)).Methods(http.MethodPost)

	log.Printf("topup-service run at : localhost:%s", servicePort)
	http.ListenAndServe(fmt.Sprintf(":%s", servicePort), router)

}

func HandleTopup(connKafka *kafka.Conn) func(w http.ResponseWriter, r *http.Request) {
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

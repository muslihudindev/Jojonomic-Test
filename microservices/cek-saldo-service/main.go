package main

import (
	"cek-saldo-service/models"
	"cek-saldo-service/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	postgresDsn := os.Getenv("POSTGRES_DSN")
	servicePort := os.Getenv("SERVICE_PORT")

	//gorm database
	connDb, err := gorm.Open(postgres.Open(postgresDsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to oppen connection to postgres:", err)
		panic(err)
	}

	//router
	router := mux.NewRouter()
	router.HandleFunc("/api/saldo", HandleSaldo(connDb)).Methods(http.MethodPost)

	log.Printf("input-saldo-service run at : localhost:%s", servicePort)
	http.ListenAndServe(fmt.Sprintf(":%s", servicePort), router)

}

func HandleSaldo(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.Params
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondError(w, http.StatusBadRequest, models.ResponseData{
				Error:   true,
				Message: err.Error(),
			})
			return
		}

		rekening := models.Rekening{}
		err := db.Where("norek = ?", req.Norek).Last(&rekening).Error
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, models.ResponseData{
				Error:   true,
				Message: err.Error(),
			})
			return
		}

		utils.RespondJSON(w, http.StatusOK, models.ResponseData{
			Error: false,
			Data: models.ResponseRekening{
				Norek: rekening.Norek,
				Saldo: rekening.Saldo,
			},
		})
	}
}

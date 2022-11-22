package main

import (
	"cek-harga-service/models"
	"cek-harga-service/utils"
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
	router.HandleFunc("/api/check-harga", HandleCheckHarga(connDb)).Methods(http.MethodGet)

	log.Printf("cek-harga-service run at : localhost:%s", servicePort)
	http.ListenAndServe(fmt.Sprintf(":%s", servicePort), router)

}

func HandleCheckHarga(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		harga := models.Harga{}
		err := db.Order("date DESC").First(&harga).Error
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, models.ResponseData{
				Error: true,
			})
			return
		}

		utils.RespondJSON(w, http.StatusOK, models.ResponseData{
			Error: false,
			Data: models.ResponseHarga{
				HargaTopup:   harga.HargaTopup,
				HargaBuyback: harga.HargaBuyback,
			},
		})
	}
}

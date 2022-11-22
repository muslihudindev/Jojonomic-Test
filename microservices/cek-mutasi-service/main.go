package main

import (
	"cek-mutasi-service/models"
	"cek-mutasi-service/utils"
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
	router.HandleFunc("/api/mutasi", HandleMutasi(connDb)).Methods(http.MethodPost)

	log.Printf("input-mutasi-service run at : localhost:%s", servicePort)
	http.ListenAndServe(fmt.Sprintf(":%s", servicePort), router)

}

func HandleMutasi(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.Params
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondError(w, http.StatusBadRequest, models.ResponseData{
				Error:   true,
				Message: err.Error(),
			})
			return
		}

		transaksi := []models.Transaksi{}
		err := db.Where("norek = ? AND date >= ? AND date <= ?", req.Norek, req.StartDate, req.EndDate).Find(&transaksi).Error
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, models.ResponseData{
				Error:   true,
				Message: err.Error(),
			})
			return
		}

		resData := []models.ResponseTransaksi{}
		for _, v := range transaksi {
			resData = append(resData, models.ResponseTransaksi{
				Norek:        v.Norek,
				Type:         v.Type,
				Gram:         v.Gram,
				HargaTopup:   v.HargaTopup,
				HargaBuyback: v.HargaBuyback,
				Saldo:        v.Saldo,
				Date:         v.Date,
			})
		}

		utils.RespondJSON(w, http.StatusOK, models.ResponseData{
			Error: false,
			Data:  resData,
		})
	}
}

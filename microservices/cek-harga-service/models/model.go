package models

type ResponseData struct {
	Error bool          `json:"error"`
	Data  ResponseHarga `json:"data"`
}

type ResponseHarga struct {
	HargaTopup   float64 `json:"harga_topup"`
	HargaBuyback float64 `json:"harga_buyback"`
}

type Harga struct {
	AdminId      string  `json:"admin_id"`
	ReffId       string  `json:"reff_id"`
	HargaTopup   float64 `json:"harga_topup"`
	HargaBuyback float64 `json:"harga_buyback"`
	Date         int64   `json:"date" gorm:"autoCreateTime"`
}

func (Harga) TableName() string {
	return "harga"
}

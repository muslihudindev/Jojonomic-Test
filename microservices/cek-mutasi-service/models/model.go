package models

type ResponseData struct {
	Error   bool                `json:"error"`
	Data    []ResponseTransaksi `json:"data"`
	ReffID  string              `json:"reff_id,omitempty"`
	Message string              `json:"message,omitempty"`
}

type ResponseTransaksi struct {
	Norek        string  `json:"norek"`
	Type         string  `json:"type"`
	Gram         float64 `json:"gram"`
	HargaTopup   float64 `json:"harga_topup"`
	HargaBuyback float64 `json:"harga_buyback"`
	Saldo        float64 `json:"saldo"`
	Date         int64   `json:"date"`
}

type Params struct {
	Norek     string `json:"norek"`
	StartDate int64  `json:"start_date"`
	EndDate   int64  `json:"end_date"`
}

type Transaksi struct {
	ReffID       string  `json:"reff_id"`
	Norek        string  `json:"norek"`
	Type         string  `json:"type"`
	Gram         float64 `json:"gram"`
	HargaTopup   float64 `json:"harga_topup"`
	HargaBuyback float64 `json:"harga_buyback"`
	Saldo        float64 `json:"saldo"`
	Date         int64   `json:"date" gorm:"autoCreateTime"`
}

func (Transaksi) TableName() string {
	return "transaksi"
}

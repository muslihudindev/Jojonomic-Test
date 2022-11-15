package models

const TopUp = `topup`

type Message struct {
	Gram  string `json:"gram"`
	Harga string `json:"harga"`
	Norek string `json:"norek"`
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

type Rekening struct {
	ReffID string  `json:"reff_id"`
	Norek  string  `json:"norek"`
	Saldo  float64 `json:"saldo"`
	Date   int64   `json:"date" gorm:"autoCreateTime"`
}

func (Rekening) TableName() string {
	return "rekening"
}

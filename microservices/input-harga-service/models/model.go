package models

type ResponseData struct {
	Error   bool   `json:"error"`
	ReffID  string `json:"reff_id,omitempty"`
	Message string `json:"message,omitempty"`
}

type Params struct {
	AdminID      string  `json:"admin_id"`
	HargaTopup   float64 `json:"harga_topup"`
	HargaBuyback float64 `json:"harga_buyback"`
}

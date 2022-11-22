package models

type ResponseData struct {
	Error   bool   `json:"error"`
	ReffID  string `json:"reff_id,omitempty"`
	Message string `json:"message,omitempty"`
}

type Params struct {
	Gram  float64 `json:"gram"`
	Harga float64 `json:"harga"`
	Norek string  `json:"norek"`
}

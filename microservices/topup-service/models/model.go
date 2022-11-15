package models

type ResponseData struct {
	Error   bool   `json:"error"`
	ReffID  string `json:"reff_id,omitempty"`
	Message string `json:"message,omitempty"`
}

type Params struct {
	Gram  string `json:"gram"`
	Saldo string `json:"saldo"`
	Norek string `json:"norek"`
}

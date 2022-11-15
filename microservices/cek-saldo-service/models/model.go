package models

type ResponseData struct {
	Error   bool             `json:"error"`
	Data    ResponseRekening `json:"data"`
	ReffID  string           `json:"reff_id,omitempty"`
	Message string           `json:"message,omitempty"`
}

type ResponseRekening struct {
	Norek string  `json:"norek"`
	Saldo float64 `json:"saldo"`
}

type Params struct {
	Norek string `json:"norek"`
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

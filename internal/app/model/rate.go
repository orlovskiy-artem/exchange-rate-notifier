package model

// TODO
type ExchangeRate struct {
	Target string  `json:"target_code"`
	Base   string  `json:"base_code"`
	Value  float32 `json:"conversion_rate"`
}

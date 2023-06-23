package structs

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string	`json:"request_id"`
	Currency     string	`json:"currency"`
	Provider     string	`json:"provider"`
	Amount       uint64	`json:"amount"`
	PaymentDt    uint64	`json:"payment_dt"`
	Bank         string	`json:"bank"`
	DeliveryCost uint64	`json:"delivery_cost"`
	GoodsTotal   uint64	`json:"goods_total"`
	CustomFee    uint64	`json:"custom_fee"`
}
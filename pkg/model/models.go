package model

import (
	"errors"
	data "withNats/pkg/structs"
)

var (
	ErrNoNatsData = errors.New(" No Data found")
)

type NatsData struct {
	OrderUID          string 			`json:"order_uid"`
	TrackNumber       string			`json:"track_number"`
	Entry             string			`json:"entry"`
	Delivery          data.Delivery		`json:"delivery"`
	Payment           data.Payment		`json:"payment"`
	Items             []data.Items		`json:"items"`
	Locale            string			`json:"locale"`
	InternalSignature string			`json:"internal_signature"`
	CustomerID        string			`json:"customer_id"`
	DeliveryService   string			`json:"delivery_service"`
	ShardKey          string			`json:"shardkey"`
	SmID              uint32			`json:"sm_id"`
	DateCreated       string			`json:"date_created"`
	OofShard          string			`json:"oof_shard"`
}

type NatsDataRepo interface {
	FindNatsData(ID string) (*NatsData, error)
	Add(nData *NatsData) error
}






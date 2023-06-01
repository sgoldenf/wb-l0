package application

import (
	"encoding/json"
	"time"

	"github.com/go-playground/validator"
	"github.com/sgoldenf/wb_l0/internal/interface/order"
)

type validationStruct struct {
	OrderUID    string `json:"order_uid" validate:"required"`
	TrackNumber string `json:"track_number" validate:"required"`
	Entry       string `json:"entry" validate:"required"`
	Delivery    struct {
		Name    string `json:"name" validate:"required"`
		Phone   string `json:"phone" validate:"required"`
		ZIP     string `json:"zip" validate:"required"`
		City    string `json:"city" validate:"required"`
		Address string `json:"address" validate:"required"`
		Region  string `json:"region" validate:"required"`
		Email   string `json:"email" validate:"required"`
	} `json:"delivery" validate:"required"`
	Payment struct {
		Transaction  string `json:"transaction" validate:"required"`
		RequestID    string `json:"request_id"`
		Currency     string `json:"currency" validate:"required"`
		Provider     string `json:"provider" validate:"required"`
		Amount       uint   `json:"amount" validate:"required"`
		PaymentDT    uint   `json:"payment_dt" validate:"required"`
		Bank         string `json:"bank" validate:"required"`
		DeliveryCost uint   `json:"delivery_cost" validate:"required"`
		GoodsTotal   uint   `json:"goods_total" validate:"required"`
		CustomFee    uint   `json:"custom_fee"`
	} `json:"payment" validate:"required"`
	Items []struct {
		ChrtID      uint   `json:"chrt_id" validate:"required"`
		TrackNumber string `json:"track_number" validate:"required"`
		Price       uint   `json:"price" validate:"required"`
		RID         string `json:"rid" validate:"required"`
		Name        string `json:"mascaras" validate:"required"`
		Sale        uint   `json:"sale" validate:"required"`
		Size        string `json:"size" validate:"required"`
		TotalPrice  uint   `json:"total_price" validate:"required"`
		NmID        uint   `json:"nm_id" validate:"required"`
		Brand       string `json:"brand" validate:"required"`
		Status      uint   `json:"status" validate:"required"`
	} `json:"items" validate:"required"`
	Locale            string    `json:"locale" validate:"required"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id" validate:"required"`
	DeliverService    string    `json:"delivery_service" validate:"required"`
	ShardKey          string    `json:"shardkey" validate:"required"`
	SmID              uint      `json:"sm_id" validate:"required"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard" validate:"required"`
}

func validateOrderJSON(data []byte) error {
	var o validationStruct
	if err := json.Unmarshal(data, &o); err != nil {
		return err
	}
	v := validator.New()
	if err := v.Struct(o); err != nil {
		return err
	}
	return nil
}

func parseOrderJSON(data []byte) (*order.Order, error) {
	o := &order.Order{Data: string(data)}
	if err := json.Unmarshal(data, o); err != nil {
		return nil, err
	}
	return o, nil
}

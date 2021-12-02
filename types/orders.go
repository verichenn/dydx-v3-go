package types

import (
	"net/url"
	"strconv"
	"time"
)

type OrderResponse struct {
	Order Order `json:"order"`
}

type CancelOrderResponse struct {
	CancelOrder Order `json:"cancelOrder"`
}

type Order struct {
	ID              string    `json:"id"`
	ClientID        string    `json:"clientId"`
	AccountID       string    `json:"accountId"`
	Market          string    `json:"market"`
	Side            string    `json:"side"`
	Price           string    `json:"price"`
	TriggerPrice    string    `json:"triggerPrice"`
	TrailingPercent string    `json:"trailingPercent"`
	Size            string    `json:"size"`
	RemainingSize   string    `json:"remainingSize"`
	Type            string    `json:"type"`
	CreatedAt       time.Time `json:"createdAt"`
	UnfillableAt    string    `json:"unfillableAt"`
	ExpiresAt       time.Time `json:"expiresAt"`
	Status          string    `json:"status"`
	TimeInForce     string    `json:"timeInForce"`
	PostOnly        bool      `json:"postOnly"`
	CancelReason    string    `json:"cancelReason"`
}

type OrderListResponse struct {
	Orders []Order `json:"orders"`
}

type OrderQueryParam struct {
	Market             string `json:"market"`
	Status             string `json:"status"`
	Type               string `json:"type"`
	Limit              int    `json:"limit"`
	Side               string `json:"side"`
	CreatedBeforeOrAt  string `json:"createdAt"`
	ReturnLatestOrders string `json:"returnLatestOrders"`
}

func (o OrderQueryParam) ToParams() url.Values {
	params := url.Values{}
	if o.Market != "" {
		params.Add("market", o.Market)
	}
	if o.Status != "" {
		params.Add("status", o.Status)
	}
	if o.Side != "" {
		params.Add("side", o.Side)
	}
	if o.Type != "" {
		params.Add("type", o.Type)
	}
	if o.Limit != 0 {
		params.Add("limit", strconv.Itoa(o.Limit))
	}
	if o.CreatedBeforeOrAt != "" {
		params.Add("createdBeforeOrAt", o.CreatedBeforeOrAt)
	}
	if o.ReturnLatestOrders != "" {
		params.Add("returnLatestOrders", o.ReturnLatestOrders)
	}
	return params
}

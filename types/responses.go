package types

import "time"

type AccountResponse struct {
	Account *Account `json:"account"`
}

type Account struct {
	StarkKey           string              `json:"starkKey"`
	PositionId         int64               `json:"positionId,string"`
	Equity             string              `json:"equity"`
	FreeCollateral     string              `json:"freeCollateral"`
	QuoteBalance       string              `json:"quoteBalance"`
	PendingDeposits    string              `json:"pendingDeposits"`
	PendingWithdrawals string              `json:"pendingWithdrawals"`
	CreatedAt          time.Time           `json:"createdAt"`
	OpenPositions      map[string]Position `json:"openPositions"`
	AccountNumber      string              `json:"accountNumber"`
	ID                 string              `json:"id"`
}

type Position struct {
	Market        string      `json:"market"`
	Status        string      `json:"status"`
	Side          string      `json:"side"`
	Size          string      `json:"size"`
	MaxSize       string      `json:"maxSize"`
	EntryPrice    string      `json:"entryPrice"`
	ExitPrice     interface{} `json:"exitPrice"`
	UnrealizedPnl string      `json:"unrealizedPnl"`
	RealizedPnl   string      `json:"realizedPnl"`
	CreatedAt     time.Time   `json:"createdAt"`
	ClosedAt      interface{} `json:"closedAt"`
	NetFunding    string      `json:"netFunding"`
	SumOpen       string      `json:"sumOpen"`
	SumClose      string      `json:"sumClose"`
}

type OrderResponse struct {
	Order *Order `json:"order"`
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

type OrderList struct {
	Orders []Order `json:"orders"`
}

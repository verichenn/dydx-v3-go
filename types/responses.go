package types

import "time"

type AccountResponseObject struct {
	starkKey           string
	positionId         string
	equity             string
	freeCollateral     string
	pendingDeposits    string
	pendingWithdrawals string
	openPositions      map[string]PositionResponseObject
	accountNumber      string
	id                 string
	quoteBalance       string
}

type PositionResponseObject struct {
	market        string
	status        string
	side          string
	size          string
	maxSize       string
	entryPrice    string
	exitPrice     string
	unrealizedPnl string
	realizedPnl   string
	createdAt     time.Time
	closedAt      time.Time
	sumOpen       string
	sumClose      string
	netFunding    string
}

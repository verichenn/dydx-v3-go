package types

type OrderQueryParam struct {
	Market             string `json:"market"`
	Status             string `json:"status"`
	Type               string `json:"type"`
	Limit              int    `json:"limit"`
	Side               string `json:"side"`
	CreatedBeforeOrAt  string `json:"createdAt"`
	ReturnLatestOrders string `json:"returnLatestOrders"`
}

package helpers

// API URLs
const (
	ApiHostMainnet = "https://api.dydx.exchange"
	ApiHostRopsten = "https://api.stage.dydx.exchange"
	WsHostMainnet  = "wss://api.dydx.exchange/v3/ws"
	WsHostRopsten  = "wss://api.stage.dydx.exchange/v3/ws"
)

// Signature Types
const (
	SignatureTypeNoPrepend   = 0
	SignatureTypeDecimal     = 1
	SignatureTypeHexadecimal = 2
)

// Off-Chain Ethereum-Signed Actions
const (
	OffChainOnboardingAction    = "dYdX Onboarding"
	OffChainKeyDerivationAction = "dYdX STARK Key"
)

// Ethereum Network IDs
const (
	NetworkIdMainnet = 1
	NetworkIdRopsten = 3
)

// Position Status Types
const (
	open       = "OPEN"
	CLOSED     = "CLOSED"
	LIQUIDATED = "LIQUIDATED"
)

// Order Types
const (
	OrderTypeLimit        = "LIMIT"
	OrderTypeMarket       = "MARKET"
	OrderTypeStop         = "STOP_LIMIT"
	OrderTypeTrailingStop = "TRAILING_STOP"
	OrderTypeTakeProfit   = "TAKE_PROFIT"
)

// Order Side
const (
	OrderSideBuy  = "BUY"
	OrderSideSell = "SELL"
)

// Time in Force Types
const (
	TimeInForceGtt = "GTT"
	TimeInForceFok = "FOK"
	TimeInForceIoc = "IOC"
)

const (
	OrderStatusPending     = "PENDING"
	OrderStatusOpen        = "OPEN"
	OrderStatusFilled      = "FILLED"
	OrderStatusCanceled    = "CANCELED"
	OrderStatusUntriggered = "UNTRIGGERED"
)

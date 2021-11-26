package helpers

const (
	SignatureTypeNoPrepend   = 0
	SignatureTypeDecimal     = 1
	SignatureTypeHexadecimal = 2

	NetworkIdMainnet = 5777
	NetworkIdRopsten = 3

	OffChainOnboardingAction    = "dYdX Onboarding"
	OffChainKeyDerivationAction = "dYdX STARK Key"

	ApiHostMainnet = "https://api.dydx.exchange"
	ApiHostRopsten = "https://api.stage.dydx.exchange"
	WsHostMainnet  = "wss://api.dydx.exchange/v3/ws"
	WsHostRopsten  = "wss://api.stage.dydx.exchange/v3/ws"
)

//position.status
const (
	open       = "OPEN"
	CLOSED     = "CLOSED"
	LIQUIDATED = "LIQUIDATED"
)

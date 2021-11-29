package dydx

import (
	"dydx-v3-go/helpers"
	"dydx-v3-go/modules"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
	"github.com/umbracle/go-web3/jsonrpc"
	"testing"
	"time"
)

const (
	DefaultHost     = "http://localhost:8080"
	EthereumAddress = "0x9Ff965Be98484736caD79C81152971E0AFe80493"
)

func TestCreateOrder(t *testing.T) {
	options := Options{}
	client := DefaultClient(options)
	account, _ := client.Private.GetAccount("")

	positions, _ := client.Private.GetPositions("BTC-USD")
	fmt.Println(positions)

	apiOrder := &modules.ApiOrder{
		ApiStarkwareSigned: modules.ApiStarkwareSigned{Expiration: expiration()},
		Market:             "BTC-USD",
		Side:               "BUY",
		Type:               "LIMIT",
		Size:               "0.001",
		Price:              "1",
		ClientId:           helpers.RandomClientId(),
		TimeInForce:        "GTT",
		PostOnly:           true,
		LimitFee:           "0.0015",
	}
	order, _ := client.Private.CreateOrder(apiOrder, account.PositionId)
	fmt.Println(order)
}

func expiration() string {
	return time.Now().Add(5 * time.Minute).UTC().Format("2006-01-02T15:04:05.999Z")
}

func TestSignViaLocalNode(t *testing.T) {
	web3, _ := jsonrpc.NewClient("http://localhost:8545")
	signer := &modules.EthWeb3Signer{Web3: web3}
	actionSinger := modules.NewSigner(signer, helpers.NetworkIdMainnet)
	sign := actionSinger.Sign(EthereumAddress,
		map[string]interface{}{"action": helpers.OffChainOnboardingAction})
	fmt.Println(sign)
}

func TestDeriveStarkKey(t *testing.T) {
	sha3 := solsha3.SoliditySHA3([]string{"address"}, "0x49EdDD3769c0712032808D86597B84ac5c2F5614")
	fmt.Println(hexutil.Encode(sha3))
}

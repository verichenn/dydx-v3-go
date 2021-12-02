package dydx

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
	"github.com/umbracle/go-web3/jsonrpc"
	"github.com/verichenn/dydx-v3-go/common"
	"github.com/verichenn/dydx-v3-go/modules"
	"github.com/verichenn/dydx-v3-go/types"
	"testing"
	"time"
)

const (
	DefaultHost     = "http://localhost:8080"
	EthereumAddress = "0x9Ff965Be98484736caD79C81152971E0AFe80493"
)

var options = Options{
	Host:                      common.ApiHostMainnet,
	StarkPublicKey:            "",
	StarkPrivateKey:           "",
	StarkPublicKeyYCoordinate: "",
	DefaultEthereumAddress:    "",
	ApiKeyCredentials: &modules.ApiKeyCredentials{
		Key:        "",
		Secret:     "",
		Passphrase: "",
	},
}

func TestGetAccount(t *testing.T) {
	client := NewClient(options)
	account, _ := client.Private.GetAccount("")
	fmt.Println(account)
}

func TestGetPositions(t *testing.T) {
	client := NewClient(options)
	positions, _ := client.Private.GetPositions("BTC-USD")
	fmt.Println(positions)
}

func TestCreateOrder(t *testing.T) {
	client := NewClient(options)
	apiOrder := &modules.ApiOrder{
		ApiBaseOrder: modules.ApiBaseOrder{Expiration: common.ExpireAfter(5 * time.Minute)},
		Market:       "BTC-USD",
		Side:         "BUY",
		Type:         "LIMIT",
		Size:         "0.001",
		Price:        "1",
		ClientId:     common.RandomClientId(),
		TimeInForce:  "GTT",
		PostOnly:     true,
		LimitFee:     "0.0015",
	}
	order, _ := client.Private.CreateOrder(apiOrder, 144336)
	fmt.Println(order)
}

func TestSignViaLocalNode(t *testing.T) {
	web3, _ := jsonrpc.NewClient("http://localhost:8545")
	signer := &modules.EthWeb3Signer{Web3: web3}
	actionSinger := modules.NewSigner(signer, common.NetworkIdMainnet)
	sign := actionSinger.Sign(EthereumAddress,
		map[string]interface{}{"action": common.OffChainOnboardingAction})
	fmt.Println(sign)
}

func TestDeriveStarkKey(t *testing.T) {
	sha3 := solsha3.SoliditySHA3([]string{"address"}, "0x49EdDD3769c0712032808D86597B84ac5c2F5614")
	fmt.Println(hexutil.Encode(sha3))
}

func TestCancelOrder(t *testing.T) {
	client := NewClient(options)
	data, err := client.Private.CancelOder("4bf8757c3ed8fb70a9c6e22f5b2fef5f4b4bd67113ed73c00f15874b2029b37")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(data)
	}
}

func TestGetOrderById(t *testing.T) {
	client := NewClient(options)
	data, err := client.Private.GetOderById("4bf8757c3ed8fb70a9c6e22f5b2fef5f4b4bd67113ed73c00f15874b2029b37")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(data.Order.ID)
	}
}

func TestGetOrders(t *testing.T) {
	client := NewClient(options)
	req := types.OrderQueryParam{
		Market: "BTC-USD",
		Limit:  100,
		Type:   "LIMIT",
	}
	data, err := client.Private.GetOrders(&req)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(data.Orders)
	}
}

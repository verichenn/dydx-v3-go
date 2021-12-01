package dydx

import (
	"dydx-v3-go/helpers"
	"dydx-v3-go/modules"
	"dydx-v3-go/types"
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

var options = Options{
	Host:                      helpers.ApiHostMainnet,
	StarkPublicKey:            "0x2c256a659da55071d90cdb27c247264b2544d4129746a07df90c97c601cbf39",
	StarkPrivateKey:           "0x24796a5e3f3b00a90553b7e97a5f76f46d8ed2f7c315cd4cd99614f717de1fe",
	starkPublicKeyYCoordinate: "0x49257237c10719d38ef7a1523fa41af41e57d8ecbdc9e5294ac2f89781c533a",
	DefaultEthereumAddress:    "0x93A0b678674BB2bAF5D47B12d33723070d2c8783",
	ApiKeyCredentials: &modules.ApiKeyCredentials{
		Key:        "ca3ace6b-849f-ff9b-9a1f-cc6e5c9d978e",
		Secret:     "ErdvqOj_YSt61LRA_71z4xcZPS29p3DWfl_KBxRb",
		Passphrase: "491trEHGp4uJ4bZ-c75R",
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
		ApiBaseOrder: modules.ApiBaseOrder{Expiration: expiration()},
		Market:       "BTC-USD",
		Side:         "BUY",
		Type:         "LIMIT",
		Size:         "0.001",
		Price:        "1",
		ClientId:     helpers.RandomClientId(),
		TimeInForce:  "GTT",
		PostOnly:     true,
		LimitFee:     "0.0015",
	}
	order, _ := client.Private.CreateOrder(apiOrder, 144336)
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

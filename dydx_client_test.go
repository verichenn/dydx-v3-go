package dydx

import (
	"dydx-v3-go/modules"
	"fmt"
	"github.com/umbracle/go-web3/jsonrpc"
	"testing"
)

const (
	//Production (Mainnet): https://api.dydx.exchange
	//Staging (Ropsten): https://api.stage.dydx.exchange
	DefaultHost     = "http://localhost:8080"
	EthereumAddress = "0x9Ff965Be98484736caD79C81152971E0AFe80493"
)

func TestCreateOrder(t *testing.T) {
	web3, _ := jsonrpc.NewClient("http://localhost:8545")
	c := DefaultClient(3, modules.ApiHostRopsten, EthereumAddress, web3)
	c.Private.GetAccount("")
}

func TestRecoverDefaultApiKeyCredentialsOnRopstenFromWeb3(t *testing.T) {
	/*web3, _ := jsonrpc.NewClient("http://localhost:8545")
	client := DefaultClient(3, DefaultHost, "", web3)
	fmt.Println(client.OnBoarding.RecoverDefaultApiCredentials(client.DefaultAddress))*/
	sData := [][]interface{}{{"bool"}, {true}}
	fmt.Println(modules.SolidityKeccak(sData))
}

func TestSignViaLocalNode(t *testing.T) {
	web3, _ := jsonrpc.NewClient("http://localhost:8545")
	signer := &modules.EthWeb3Signer{Web3: web3}
	actionSinger := modules.NewSigner(signer, modules.NetworkIdMainnet)
	sign := actionSinger.Sign(EthereumAddress,
		map[string]interface{}{"action": modules.OffChainOnboardingAction})
	fmt.Println(sign)
}

func TestDeriveStarkKey(t *testing.T) {
	web3, _ := jsonrpc.NewClient("http://localhost:8545")
	c := DefaultClient(3, DefaultHost, "", web3)

	key := c.OnBoarding.DeriveStarkKey(EthereumAddress)
	fmt.Println(key)
}

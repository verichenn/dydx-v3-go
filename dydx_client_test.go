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
	DefaultHost = "http://localhost:8080"
)

func TestCreateOrder(t *testing.T) {
	web3, _ := jsonrpc.NewClient("http://localhost:8545")
	c := DefaultClient(1, "", "", web3)
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
	signer := &modules.EthWeb3Signer{web3}
	actionSinger := modules.NewSigner(signer, modules.NetworkIdMainnet)
	sign := actionSinger.Sign("0xE165Cc13943fb8ba027E2FB5121092f63447a2a4",
		map[string]interface{}{"action": modules.OffChainOnboardingAction})
	fmt.Println(sign)
}

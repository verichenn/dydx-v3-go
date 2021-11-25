package dydx

import (
	"dydx-v3-go/modules"
	"github.com/umbracle/go-web3/jsonrpc"
	"strings"
)

const (
	WebProviderUrl  = "http://localhost:8545"
	EthereumAddress = "0x22d491Bde2303f2f43325b2108D26f1eAbA1e32b"
)

type Client struct {
	Host                      string
	ApiTimeout                int
	EthSendOptions            interface{}
	StarkPrivateKey           string
	ApiKeyCredentials         modules.ApiKeyCredentials
	starkPublicKeyYCoordinate string
	Web3                      *jsonrpc.Client
	EthSigner                 *modules.EthWeb3Signer
	DefaultAddress            string
	NetworkId                 int

	Private    *modules.Private
	OnBoarding *modules.OnBoarding
}

type Options struct {
}

func DefaultClient(networkId int, host, ethAddress string, web3 *jsonrpc.Client) *Client {
	if strings.HasSuffix(host, "/") {
		host = host[:len(host)-1]
	}
	ethSigner := &modules.EthWeb3Signer{Web3: web3}
	accounts, _ := web3.Eth().Accounts()
	defaultAddress := accounts[1].String()
	if ethAddress != "" {
		defaultAddress = ethAddress
	}
	net, _ := web3.Net().Version()
	if networkId == 0 {
		networkId = int(net)
	}
	client := &Client{
		Host:           host,
		ApiTimeout:     3000,
		NetworkId:      networkId,
		Web3:           web3,
		DefaultAddress: defaultAddress,
		EthSigner:      ethSigner,
	}
	client.OnBoarding = &modules.OnBoarding{
		Host:       host,
		EthSigner:  ethSigner,
		NetworkId:  networkId,
		EthAddress: defaultAddress,
		Singer:     modules.NewSigner(ethSigner, networkId),
	}
	//client.ApiKeyCredentials = client.OnBoarding.RecoverDefaultApiCredentials(defaultAddress)
	return client
}

package dydx

import (
	"dydx-v3-go/modules"
	"github.com/umbracle/go-web3/jsonrpc"
	"log"
	"os"
	"strings"
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
	Logger                    *log.Logger

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
	defaultAddress := accounts[0].String()
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
		Logger:         log.New(os.Stderr, "dydx-v3-go ", log.LstdFlags),
	}
	client.OnBoarding = &modules.OnBoarding{
		Host:       host,
		EthSigner:  ethSigner,
		NetworkId:  networkId,
		EthAddress: defaultAddress,
		Singer:     modules.NewSigner(ethSigner, networkId),
		Logger:     client.Logger,
	}

	apiKeyCredentials := client.OnBoarding.RecoverDefaultApiCredentials(defaultAddress)
	starkKey := client.OnBoarding.DeriveStarkKey(defaultAddress)

	client.ApiKeyCredentials = apiKeyCredentials
	client.StarkPrivateKey = starkKey

	client.Private = &modules.Private{
		Host:              host,
		NetworkId:         networkId,
		StarkPrivateKey:   starkKey,
		DefaultAddress:    defaultAddress,
		ApiKeyCredentials: apiKeyCredentials,
		Logger:            client.Logger,
	}
	return client
}

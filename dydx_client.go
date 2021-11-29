package dydx

import (
	"dydx-v3-go/helpers"
	"dydx-v3-go/modules"
	"github.com/umbracle/go-web3/jsonrpc"
	"log"
	"os"
	"strings"
)

type Client struct {
	Host                      string
	ApiTimeout                int
	StarkPublicKey            string
	StarkPrivateKey           string
	StarkPublicKeyYCoordinate string
	ApiKeyCredentials         *modules.ApiKeyCredentials

	Web3                   *jsonrpc.Client
	EthSigner              modules.EthSigner
	DefaultEthereumAddress string
	NetworkId              int
	Logger                 *log.Logger

	Private    *modules.Private
	OnBoarding *modules.OnBoarding
}

type Options struct {
	Host                      string
	StarkPublicKey            string
	StarkPrivateKey           string
	starkPublicKeyYCoordinate string
	DefaultEthereumAddress    string
	ApiKeyCredentials         *modules.ApiKeyCredentials

	Web3      *jsonrpc.Client
	NetworkId int
}

func DefaultClient(options Options) *Client {
	host := options.Host
	if strings.HasSuffix(host, "/") {
		host = host[:len(host)-1]
	}
	client := &Client{
		Host:       host,
		ApiTimeout: 3000,
		Logger:     log.New(os.Stderr, "dydx-v3-go ", log.LstdFlags),
	}
	web3 := options.Web3
	if web3 != nil {
		net, _ := web3.Net().Version()
		networkId := options.NetworkId
		if networkId == 0 {
			networkId = int(net)
		}

		client.Web3 = web3
		client.EthSigner = &modules.EthWeb3Signer{Web3: web3}
		client.NetworkId = networkId
	}
	if options.StarkPrivateKey != "" {
		client.StarkPrivateKey = options.StarkPrivateKey
		client.EthSigner = &modules.EthKeySinger{PrivateKey: options.StarkPrivateKey}
	}
	client.DefaultEthereumAddress = options.DefaultEthereumAddress
	if client.NetworkId == 0 {
		client.NetworkId = helpers.NetworkIdMainnet
	}

	if options.StarkPublicKey != "" && options.starkPublicKeyYCoordinate != "" {
		client.StarkPublicKey = options.StarkPublicKey
		client.StarkPublicKeyYCoordinate = options.starkPublicKeyYCoordinate
	}
	client.OnBoarding = &modules.OnBoarding{
		Host:       host,
		EthSigner:  client.EthSigner,
		NetworkId:  client.NetworkId,
		EthAddress: client.DefaultEthereumAddress,
		Singer:     modules.NewSigner(client.EthSigner, client.NetworkId),
		Logger:     client.Logger,
	}
	if options.ApiKeyCredentials != nil {
		client.ApiKeyCredentials = options.ApiKeyCredentials
	} else {
		client.ApiKeyCredentials = client.OnBoarding.RecoverDefaultApiCredentials(client.DefaultEthereumAddress)
	}

	client.Private = &modules.Private{
		Host:              host,
		NetworkId:         client.NetworkId,
		StarkPrivateKey:   client.StarkPrivateKey,
		DefaultAddress:    client.DefaultEthereumAddress,
		ApiKeyCredentials: client.ApiKeyCredentials,
		Logger:            client.Logger,
	}
	return client
}

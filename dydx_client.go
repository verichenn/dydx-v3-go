package dydx

import (
	"github.com/umbracle/go-web3/jsonrpc"
	"github.com/verichenn/dydx-v3-go/common"
	"github.com/verichenn/dydx-v3-go/modules"
	"log"
	"os"
	"strings"
	"time"
)

type Client struct {
	options Options

	Host                      string
	ApiTimeout                time.Duration
	StarkPublicKey            string
	StarkPrivateKey           string
	StarkPublicKeyYCoordinate string
	ApiKeyCredentials         *modules.ApiKeyCredentials

	Web3           *jsonrpc.Client
	EthSigner      modules.EthSigner
	DefaultAddress string
	NetworkId      int
	Logger         *log.Logger

	Private    *modules.Private
	OnBoarding *modules.OnBoarding
}

type Options struct {
	Host                      string
	StarkPublicKey            string
	StarkPrivateKey           string
	StarkPublicKeyYCoordinate string
	DefaultEthereumAddress    string
	ApiKeyCredentials         *modules.ApiKeyCredentials

	Web3      *jsonrpc.Client
	NetworkId int
}

func NewClient(options Options) *Client {
	client := &Client{
		Host:              strings.TrimPrefix(options.Host, "/"),
		ApiTimeout:        3 * time.Second,
		DefaultAddress:    options.DefaultEthereumAddress,
		StarkPublicKey:    options.StarkPublicKey,
		StarkPrivateKey:   options.StarkPrivateKey,
		ApiKeyCredentials: options.ApiKeyCredentials,
		Logger:            log.New(os.Stderr, "dydx-v3-go ", log.LstdFlags),
	}

	if options.Web3 != nil {
		networkId := options.NetworkId
		if networkId == 0 {
			net, _ := options.Web3.Net().Version()
			networkId = int(net)
		}

		client.Web3 = options.Web3
		client.EthSigner = &modules.EthWeb3Signer{Web3: options.Web3}
		client.NetworkId = networkId
	}

	if client.NetworkId == 0 {
		client.NetworkId = common.NetworkIdMainnet
	}

	if options.StarkPrivateKey != "" {
		client.StarkPrivateKey = options.StarkPrivateKey
		client.EthSigner = &modules.EthKeySinger{PrivateKey: options.StarkPrivateKey}
	}

	client.OnBoarding = &modules.OnBoarding{
		Host:       client.Host,
		EthSigner:  client.EthSigner,
		NetworkId:  client.NetworkId,
		EthAddress: client.DefaultAddress,
		Singer:     modules.NewSigner(client.EthSigner, client.NetworkId),
		Logger:     client.Logger,
	}
	if options.ApiKeyCredentials == nil {
		client.ApiKeyCredentials = client.OnBoarding.RecoverDefaultApiCredentials(client.DefaultAddress)
	}

	client.Private = &modules.Private{
		Host:              client.Host,
		NetworkId:         client.NetworkId,
		StarkPrivateKey:   client.StarkPrivateKey,
		DefaultAddress:    client.DefaultAddress,
		ApiKeyCredentials: client.ApiKeyCredentials,
		Logger:            client.Logger,
	}
	return client
}

package modules

import (
	"github.com/umbracle/go-web3/jsonrpc"
	"github.com/verichenn/dydx-v3-go/common"
)

type EthSigner interface {
	sign(eip712Message map[string]interface{}, messageHash, optSingerAddress string) string
}

type EthWeb3Signer struct {
	Web3 *jsonrpc.Client
}

func (web3Singer *EthWeb3Signer) sign(eip712Message map[string]interface{}, messageHash, address string) string {
	rawSignature := signTypedData(eip712Message, web3Singer, address)
	return common.CreateTypedSignature(rawSignature, common.SignatureTypeNoPrepend)
}

//https://github.com/dydxprotocol/dydx-v3-python/issues/62
func signTypedData(eip712Message map[string]interface{}, web3Singer *EthWeb3Signer, address string) string {
	var out string
	if err := web3Singer.Web3.Call("eth_signTypedData", &out, address, eip712Message); err != nil {
		panic(err)
	}
	return out
}

type EthKeySinger struct {
	Address    string
	PrivateKey string
}

func (keySinger EthKeySinger) sign(eip712Message map[string]interface{}, messageHash, optSingerAddress string) string {
	panic("implement me")
}

package modules

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
	"github.com/verichenn/dydx-v3-go/common"
)

const (
	Domain                       = "dYdX"
	Version                      = "1.0"
	Eip712DomainStringNoContract = "EIP712Domain(string name,string version,uint256 chainId)"
)

var (
	Eip712OnboardingActionStruct = []map[string]string{
		{"type": "string", "name": "action"},
		{"type": "string", "name": "onlySignOn"},
	}
	Eip712OnboardingActionStructString = "dYdX(string action,string onlySignOn)"

	Eip712OnboardingActionStructTestnet = []map[string]string{
		{"type": "string", "name": "action"},
	}
	Eip712OnboardingActionStructStringTestnet = "dYdX(string action)"
	Eip712StructName                          = "dYdX"
	OnlySignOnDomainMainnet                   = "https://trade.dydx.exchange"
)

type SignOnboardingAction struct {
	Signer    EthSigner
	NetworkId int
}

func NewSigner(signer EthSigner, networkId int) *SignOnboardingAction {
	return &SignOnboardingAction{signer, networkId}
}

func (a *SignOnboardingAction) Sign(signerAddress string, message map[string]interface{}) string {
	eip712Message := a.GetEIP712Message(message)
	action := message["action"].(string)
	messageHash := a.GetHash(action)
	typedSignature := a.Signer.sign(eip712Message, messageHash, signerAddress)
	return typedSignature
}

func (a *SignOnboardingAction) GetEIP712Message(message map[string]interface{}) map[string]interface{} {
	structName := a.GetEIP712StructName()
	eip712Message := map[string]interface{}{
		"types": map[string]interface{}{
			"EIP712Domain": []map[string]string{
				{
					"name": "name",
					"type": "string",
				},
				{
					"name": "version",
					"type": "string",
				},
				{
					"name": "chainId",
					"type": "uint256",
				},
			},
			structName: a.GetEIP712Struct(),
		},
		"domain": map[string]interface{}{
			"name":    Domain,
			"version": Version,
			"chainId": a.NetworkId,
		},
		"primaryType": structName,
		"message":     message,
	}
	if a.NetworkId == common.NetworkIdMainnet {
		msg := eip712Message["message"].(map[string]interface{})
		msg["onlySignOn"] = "https://trade.dydx.exchange"
	}

	return eip712Message
}

func (a *SignOnboardingAction) GetEip712Hash(structHash string) string {
	fact := solsha3.SoliditySHA3(
		[]string{"bytes2", "bytes32", "bytes32"},
		[]interface{}{"0x1901", a.GetDomainHash(), structHash},
	)
	return fmt.Sprintf("0x%x", fact)
}

func (a *SignOnboardingAction) GetDomainHash() string {
	fact := solsha3.SoliditySHA3(
		[]string{"bytes32", "bytes32", "bytes32", "uint256"},
		[]interface{}{common.HashString(Eip712DomainStringNoContract), common.HashString(Domain), common.HashString(Version), a.NetworkId},
	)
	return fmt.Sprintf("0x%x", fact)
}

func (a *SignOnboardingAction) GetEIP712Struct() []map[string]string {
	if a.NetworkId == common.NetworkIdMainnet {
		return Eip712OnboardingActionStruct
	} else {
		return Eip712OnboardingActionStructTestnet
	}
}

func (a *SignOnboardingAction) GetEIP712StructName() string {
	return Eip712StructName
}

func (a *SignOnboardingAction) GetHash(action string) string {
	var eip712StructStr string
	if a.NetworkId == common.NetworkIdMainnet {
		eip712StructStr = Eip712OnboardingActionStructString
	} else {
		eip712StructStr = Eip712OnboardingActionStructStringTestnet
	}
	data := [][]string{
		{"bytes32", "bytes32"},
		{common.HashString(eip712StructStr), common.HashString(action)},
	}
	if a.NetworkId == common.NetworkIdMainnet {
		data[0] = append(data[0], "bytes32")
		data[1] = append(data[1], common.HashString(OnlySignOnDomainMainnet))
	}
	structHash := solsha3.SoliditySHA3(data[0], data[1])
	return a.GetEip712Hash(hexutil.Encode(structHash))
}

package modules

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
	"github.com/verichenn/dydx-v3-go/helpers"
	"log"
	"math/big"
	"strings"
)

type OnBoarding struct {
	Host                      string
	EthSigner                 EthSigner
	NetworkId                 int
	EthAddress                string
	StarkPublicKey            string
	StarkPublicKeyYCoordinate string
	Singer                    *SignOnboardingAction
	Logger                    *log.Logger
}

type ApiKeyCredentials struct {
	Key        string
	Secret     string
	Passphrase string
}

func (board OnBoarding) RecoverDefaultApiCredentials(ethereumAddress string) *ApiKeyCredentials {
	signature := board.Singer.Sign(ethereumAddress, map[string]interface{}{"action": helpers.OffChainOnboardingAction})
	rHex := signature[2:66]
	rInt, _ := math.MaxBig256.SetString(rHex, 16)

	hashedRBytes := solsha3.SoliditySHA3([]string{"uint256"}, rInt.String())
	secretBytes := hashedRBytes[:30]
	sHex := signature[66:130]
	sInt, _ := math.MaxBig256.SetString(sHex, 16)

	hashedSBytes := solsha3.SoliditySHA3([]string{"uint256"}, sInt.String())
	keyBytes := hashedSBytes[:16]
	passphraseBytes := hashedSBytes[16:31]

	keyHex := hex.EncodeToString(keyBytes)
	keyUuid := strings.Join([]string{
		keyHex[:8],
		keyHex[8:12],
		keyHex[12:16],
		keyHex[16:20],
		keyHex[20:],
	}, "-")
	return &ApiKeyCredentials{
		Secret:     base64.URLEncoding.EncodeToString(secretBytes),
		Key:        keyUuid,
		Passphrase: base64.URLEncoding.EncodeToString(passphraseBytes),
	}
}

func (board OnBoarding) DeriveStarkKey(ethereumAddress string) string {
	signature := board.Singer.Sign(ethereumAddress, map[string]interface{}{"action": helpers.OffChainKeyDerivationAction})
	sig, _ := new(big.Int).SetString(signature, 0)

	sha3 := solsha3.SoliditySHA3([]string{"uint256"}, sig.String())
	hashedSignature := hexutil.Encode(sha3)

	privateKey, _ := new(big.Int).SetString(hashedSignature, 0)
	privateKey = new(big.Int).Rsh(privateKey, 5)
	return fmt.Sprintf("0x%s", privateKey.Text(16))
}

func (board OnBoarding) sign(signerAddress, action string) string {
	return board.Singer.Sign(signerAddress, map[string]interface{}{"action": action})
}

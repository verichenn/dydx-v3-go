package modules

import (
	"fmt"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
	"math/rand"
	"strconv"
	"strings"
)

func RandomClientId() string {
	return strconv.Itoa(rand.Int())
}

func createTypedSignature(signature string, sigType int) string {
	return fmt.Sprintf("%s0%s", fixRawSignature(signature), strconv.Itoa(sigType))
}

func fixRawSignature(signature string) string {
	stripped := stripHexPrefix(signature)
	if len(stripped) != 130 {
		panic(fmt.Sprintf("Invalid raw signature: %s", signature))
	}
	rs := stripped[:128]
	v := stripped[128:130]
	if v == "00" {
		return "0x" + rs + "1b"
	}
	if v == "01" {
		return "0x" + rs + "1c"
	}
	if v == "1b" || v == "1c" {
		return "0x" + stripped
	}
	panic(fmt.Sprintf("Invalid v value: %s", v))
}

func stripHexPrefix(input string) string {
	if strings.HasPrefix(input, "0x") {
		return input[2:]
	}
	return input
}

func HashString(input string) string {

	fact := solsha3.SoliditySHA3([]string{"string"}, []string{input})
	return fmt.Sprintf("0x%x", fact)
}

func SolidityKeccak(data ...interface{}) string {
	fact := solsha3.SoliditySHA3(data)
	return fmt.Sprintf("0x%x", fact)
}

func EcRecoverTypedSignature(hashVal, typedSignature string) string {
	return ""
}

func AddressAreEqual(addrOne, AddrTow string) bool {
	return false
}

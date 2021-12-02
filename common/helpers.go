package helpers

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func RandomClientId() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%016d", rand.Intn(10000000000000000))
}

func GenerateQueryPath(url string, params url.Values) string {
	if len(params) == 0 {
		return url
	}
	return fmt.Sprintf("%s?%s", url, params.Encode())
}

func CreateTypedSignature(signature string, sigType int) string {
	return fmt.Sprintf("%s0%s", fixRawSignature(signature), strconv.Itoa(sigType))
}

func fixRawSignature(signature string) string {
	stripped := strings.TrimPrefix(signature, "0x")
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

func HashString(input string) string {
	return hexutil.Encode(solsha3.SoliditySHA3([]string{"string"}, input))
}

func ExpireAfter(duration time.Duration) string {
	return time.Now().Add(duration).UTC().Format("2006-01-02T15:04:05.999Z")
}

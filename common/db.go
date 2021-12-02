package common

import (
	"github.com/satori/go.uuid"
	"strconv"
	"strings"
)

var namespace = Must(FromString("0f9da948-a6fb-4c45-9edc-4685c3f3317d"))

func getUserId(address string) string {
	return uuid.NewV5(namespace, address).String()
}

func GetAccountId(address string) string {
	return uuid.NewV5(namespace, getUserId(strings.ToLower(address))+strconv.Itoa(0)).String()
}

func FromString(input string) (u uuid.UUID, err error) {
	err = u.UnmarshalText([]byte(input))
	return
}

func Must(u uuid.UUID, err error) uuid.UUID {
	if err != nil {
		panic(err)
	}
	return u
}

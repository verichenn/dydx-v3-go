package modules

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"dydx-v3-go/helpers"
	"dydx-v3-go/types"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/yanue/starkex"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Private struct {
	Host              string
	NetworkId         int
	StarkPrivateKey   string
	DefaultAddress    string
	ApiKeyCredentials *ApiKeyCredentials
	Logger            *log.Logger
}

type ApiStarkwareSigned struct {
	Signature  string `json:"signature"`
	Expiration string `json:"expiration"`
}

type ApiOrder struct {
	ApiStarkwareSigned
	Market          string `json:"market"`
	Side            string `json:"side"`
	Type            string `json:"type"`
	Size            string `json:"size"`
	Price           string `json:"price"`
	ClientId        string `json:"clientId"`
	TimeInForce     string `json:"timeInForce"`
	PostOnly        bool   `json:"postOnly"`
	LimitFee        string `json:"limitFee"`
	CancelId        string `json:"cancelId,omitempty"`
	TriggerPrice    string `json:"triggerPrice,omitempty"`
	TrailingPercent string `json:"trailingPercent,omitempty"`
	//NetworkId              int    `json:"network_id"`
	//PositionId             int64  `json:"position_id"`
}

func (p Private) GetAccount(ethereumAddress string) (*types.Account, error) {
	if ethereumAddress == "" {
		ethereumAddress = p.DefaultAddress
	}
	uri := fmt.Sprintf("accounts/%s", helpers.GetAccountId(ethereumAddress))
	res, _ := p.get(uri, nil)
	var accountResponse *types.AccountResponse
	err := json.Unmarshal(res, &accountResponse)
	if err != nil {
		return nil, err
	}
	return accountResponse.Account, nil
}

func (p Private) CreateOrder(input *ApiOrder, positionId int64) {
	orderSignParam := starkex.OrderSignParam{
		NetworkId:  p.NetworkId,
		PositionId: positionId,
		Market:     input.Market,
		Side:       input.Side,
		HumanSize:  input.Size,
		HumanPrice: input.Price,
		LimitFee:   input.LimitFee,
		ClientId:   input.ClientId,
		Expiration: input.Expiration,
	}
	privateKey := p.StarkPrivateKey[2:]
	signature, err := starkex.OrderSign(privateKey, orderSignParam)
	if err != nil {
		fmt.Println(err)
		return
	}
	input.Signature = signature
	p.post("orders", input)
}

func (p Private) get(endpoint string, params url.Values) ([]byte, error) {
	return p.execute(http.MethodGet, helpers.GenerateQueryPath(endpoint, params), "")
}

func (p Private) post(endpoint string, data interface{}) ([]byte, error) {
	marshalData, _ := json.Marshal(data)
	return p.execute(http.MethodPost, endpoint, string(marshalData))
}

func (p Private) execute(method, endpoint string, data string) ([]byte, error) {
	isoTimestamp := generateNowISO()
	requestPath := fmt.Sprintf("/v3/%s", endpoint)
	headers := map[string]string{
		"DYDX-SIGNATURE":  p.sign(requestPath, method, isoTimestamp, data),
		"DYDX-API-KEY":    p.ApiKeyCredentials.Key,
		"DYDX-TIMESTAMP":  isoTimestamp,
		"DYDX-PASSPHRASE": p.ApiKeyCredentials.Passphrase,
	}
	resp, err := p.doExecute(method, requestPath, headers, data)
	if err != nil {
		p.Logger.Panic(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		p.Logger.Printf("wrong status code: %d,err msg:%s", resp.StatusCode, buf.String())
		return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		p.Logger.Panic(err)
		return nil, err
	}
	return respBody, nil
}

func (p Private) doExecute(method string, requestPath string, headers map[string]string, data string) (*http.Response, error) {
	var req *http.Request
	requestPath = fmt.Sprintf("%s%s", p.Host, requestPath)
	req, err := http.NewRequest(method, requestPath, strings.NewReader(data))
	if err != nil {
		return nil, errors.New("new request is fail: %v ")
	}

	for key, val := range headers {
		req.Header.Add(key, val)
	}
	req.Header.Add("Content-Type", "application/json")

	c := &http.Client{
		Timeout: time.Second * 3,
	}
	return c.Do(req)

}

func generateNowISO() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
}

func (p Private) sign(requestPath, method, isoTimestamp, body string) string {
	message := fmt.Sprintf("%s%s%s%s", isoTimestamp, method, requestPath, body)
	secret, _ := base64.URLEncoding.DecodeString(p.ApiKeyCredentials.Secret)
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(message))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

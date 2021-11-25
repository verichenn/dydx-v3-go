package modules

import (
	"crypto/hmac"
	"crypto/sha256"
	"dydx-v3-go/helpers"
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
	ApiKeyCredentials ApiKeyCredentials
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
	TimeInForce     string `json:"time_in_force"`
	PostOnly        string `json:"post_only"`
	LimitFee        string `json:"limit_fee"`
	CancelId        string `json:"cancel_id"`
	TriggerPrice    string `json:"trigger_price"`
	TrailingPercent string `json:"trailing_percent"`
	//NetworkId              int    `json:"network_id"`
	//PositionId             int64  `json:"position_id"`
}

func (p Private) GetAccount(ethereumAddress string) {
	if ethereumAddress == "" {
		ethereumAddress = p.DefaultAddress
	}
	uri := fmt.Sprintf("accounts/%s", helpers.GetAccountId(ethereumAddress))
	res, _ := p.get(uri, nil)
	fmt.Println(string(res))
}

func (p Private) CreateOrder(input *ApiOrder, positionId int64) {
	if input.ClientId == "0" {
		input.ClientId = RandomClientId()
	}
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
	signature, _ := starkex.OrderSign(p.StarkPrivateKey, orderSignParam)
	input.Signature = signature
	p.post("orders", input)
}

func (p Private) get(endpoint string, params url.Values) ([]byte, error) {
	return p.execute(http.MethodGet, endpoint, "{}")
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
		p.Logger.Println("wrong status code: ", resp.StatusCode)
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

	c := &http.Client{
		Timeout: time.Second * 3,
	}
	return c.Do(req)

}

func generateNowISO() string {
	return time.Now().Format("2006-01-02T15:04:05.999Z")
}

func (p Private) sign(requestPath, method, isoTimestamp, body string) string {
	message := fmt.Sprintf("%s%s%s%s", isoTimestamp, method, requestPath, body)
	key := []byte(base64.URLEncoding.EncodeToString([]byte(p.ApiKeyCredentials.Secret)))
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	signData := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return signData
}

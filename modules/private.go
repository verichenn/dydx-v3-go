package modules

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/yanue/starkex"
	"io/ioutil"
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

func NewPrivate(host string) *Private {
	p := &Private{
		Host:              host,
		NetworkId:         0,
		StarkPrivateKey:   "",
		DefaultAddress:    "",
		ApiKeyCredentials: ApiKeyCredentials{},
	}
	return p
}

func (p Private) GetAccount(ethereumAddress string) {
	uri := fmt.Sprintf("/accounts/%s", ethereumAddress)
	p.get(uri, nil)
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
	encodeParams := params.Encode()
	requestUrl := fmt.Sprintf("%s?%s", endpoint, encodeParams)
	return p.execute(http.MethodGet, requestUrl, "")

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
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
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
	messageString := fmt.Sprintf("%s%s%s%s", isoTimestamp, method, requestPath, body)
	h := hmac.New(sha256.New, []byte(base64.StdEncoding.EncodeToString([]byte(p.ApiKeyCredentials.Secret))))
	h.Write([]byte(messageString))
	sha := hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(sha))

}

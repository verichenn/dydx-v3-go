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

type ApiBaseOrder struct {
	Signature  string `json:"signature"`
	Expiration string `json:"expiration"`
}

type ApiOrder struct {
	ApiBaseOrder
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
}

// GetAccount 查询账户
// see https://docs.dydx.exchange/?json#get-account
func (p Private) GetAccount(ethereumAddress string) (types.AccountResponse, error) {
	if ethereumAddress == "" {
		ethereumAddress = p.DefaultAddress
	}
	uri := fmt.Sprintf("accounts/%s", helpers.GetAccountId(ethereumAddress))
	res, _ := p.get(uri, nil)
	accountResponse := types.AccountResponse{}
	if err := json.Unmarshal(res, &accountResponse); err != nil {
		return accountResponse, err
	}
	return accountResponse, nil
}

// CreateOrder 创建订单
// see https://docs.dydx.exchange/?json#create-a-new-order
func (p Private) CreateOrder(input *ApiOrder, positionId int64) (types.OrderResponse, error) {
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
	signature, err := starkex.OrderSign(p.StarkPrivateKey[2:], orderSignParam)

	orderResponse := types.OrderResponse{}
	if err != nil {
		return orderResponse, errors.New("sign error")
	}
	input.Signature = signature
	res, _ := p.post("orders", input)

	if err = json.Unmarshal(res, &orderResponse); err != nil {
		return orderResponse, err
	}
	return orderResponse, nil
}

// GetPositions 查询持仓
// see https://docs.dydx.exchange/?json#get-positions
func (p Private) GetPositions(market string) (types.PositionResponse, error) {
	params := url.Values{}
	params.Add("market", market)
	res, err := p.get("positions", params)
	position := types.PositionResponse{}
	if err != nil {
		return position, errors.New("request error")
	}
	if err = json.Unmarshal(res, &position); err != nil {
		return position, errors.New("json parser error")
	}
	return position, nil
}

// GetOrders 查询订单列表
// see https://docs.dydx.exchange/?json#get-orders
func (p Private) GetOrders(input *types.OrderQueryParam) (types.OrderListResponse, error) {
	var result types.OrderListResponse
	data, err := p.get("orders", input.ToParams())
	if err == nil {
		if err := json.Unmarshal(data, &result); err != nil {
			return result, err
		}
		return result, nil
	}
	return result, err
}

// CancelOder 取消订单
// see https://docs.dydx.exchange/?json#cancel-an-order
func (p Private) CancelOder(orderId string) (types.OrderCancelResponse, error) {
	data, err := p.delete("orders/"+orderId, nil)
	result := types.OrderCancelResponse{}
	if err != nil {
		return result, err
	}
	if err = json.Unmarshal(data, &result); err != nil {
		return result, err
	}
	return result, nil
}

// GetOderById 查询订单
// see https://docs.dydx.exchange/?json#get-order-by-id
func (p Private) GetOderById(orderId string) (types.OrderResponse, error) {
	res, err := p.get("orders/"+orderId, nil)

	var orderResponse types.OrderResponse
	if err == nil {
		if err := json.Unmarshal(res, &orderResponse); err != nil {
			return orderResponse, err
		}
		return orderResponse, nil
	}
	return orderResponse, err
}

func (p Private) get(endpoint string, params url.Values) ([]byte, error) {
	return p.request(http.MethodGet, helpers.GenerateQueryPath(endpoint, params), "")
}

func (p Private) post(endpoint string, data interface{}) ([]byte, error) {
	marshalData, _ := json.Marshal(data)
	return p.request(http.MethodPost, endpoint, string(marshalData))
}

func (p Private) delete(endpoint string, params url.Values) ([]byte, error) {
	return p.request(http.MethodGet, helpers.GenerateQueryPath(endpoint, params), "")
}

func (p Private) request(method, endpoint string, data string) ([]byte, error) {
	isoTimestamp := generateNowISO()
	requestPath := fmt.Sprintf("/v3/%s", endpoint)
	headers := map[string]string{
		"DYDX-SIGNATURE":  p.sign(requestPath, method, isoTimestamp, data),
		"DYDX-API-KEY":    p.ApiKeyCredentials.Key,
		"DYDX-TIMESTAMP":  isoTimestamp,
		"DYDX-PASSPHRASE": p.ApiKeyCredentials.Passphrase,
	}
	resp, err := p.execute(method, requestPath, headers, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		p.Logger.Printf("uri:%s, code: %d, err msg:%s", requestPath, resp.StatusCode, buf.String())
		return nil, fmt.Errorf("uri:%v , status code: %d", requestPath, resp.StatusCode)
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	p.Logger.Printf("uri:%s,response body:%s", requestPath, responseBody)
	return responseBody, err
}

func (p Private) execute(method string, requestPath string, headers map[string]string, data string) (*http.Response, error) {
	requestPath = fmt.Sprintf("%s%s", p.Host, requestPath)
	req, _ := http.NewRequest(method, requestPath, strings.NewReader(data))

	for key, val := range headers {
		req.Header.Add(key, val)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "dydx/go")

	c := &http.Client{
		Timeout: time.Second * 5,
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

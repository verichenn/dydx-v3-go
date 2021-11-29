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
	"strconv"
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
	if err := json.Unmarshal(res, &accountResponse); err != nil {
		return nil, err
	}
	return accountResponse.Account, nil
}

func (p Private) CreateOrder(input *ApiOrder, positionId int64) (*types.Order, error) {
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
	if err != nil {
		return nil, errors.New("sign error")
	}
	input.Signature = signature
	res, _ := p.post("orders", input)
	var orderResponse *types.OrderResponse
	if err = json.Unmarshal(res, &orderResponse); err != nil {
		return nil, err
	}
	return orderResponse.Order, nil
}

func (p Private) GetPositions(market string) (*types.Position, error) {
	params := url.Values{}
	params.Add("market", market)
	res, err := p.get("positions", params)
	if err != nil {
		return nil, errors.New("request error")
	}
	var position *types.Position
	if err = json.Unmarshal(res, position); err != nil {
		return nil, errors.New("json passer error")
	}
	return position, nil
}

func (p Private) GetOrder(input *types.OrderQueryParam) ([]byte, error) {
	orderQuery := url.Values{}
	if input.Market != "" {
		orderQuery.Add("market", input.Market)
	}
	if input.Status != "" {
		orderQuery.Add("status", input.Status)
	}

	if input.Side != "" {
		orderQuery.Add("side", input.Side)
	}

	if input.Type != "" {
		orderQuery.Add("type", input.Type)
	}

	if input.Limit != 0 {
		orderQuery.Add("limit", strconv.Itoa(input.Limit))
	}

	if input.CreatedBeforeOrAt != "" {
		orderQuery.Add("createdBeforeOrAt", input.CreatedBeforeOrAt)
	}

	if input.ReturnLatestOrders != "" {
		orderQuery.Add("returnLatestOrders", input.ReturnLatestOrders)
	}
	return p.get("orders", orderQuery)
}

//取消订单
func (p Private) CancelOder(orderId string) ([]byte, error) {
	return p.delete("orders/"+orderId, nil)
}

//取消订单
func (p Private) GetOderById(orderId string) (*types.OrderResponse, error) {
	res, reqErr := p.get("orders/"+orderId, nil)

	if reqErr == nil {
		var orderResponse *types.OrderResponse
		if err := json.Unmarshal(res, &orderResponse); err != nil {
			return nil, err
		}
		return orderResponse, nil
	}
	return nil, reqErr
}

func (p Private) get(endpoint string, params url.Values) ([]byte, error) {
	return p.execute(http.MethodGet, helpers.GenerateQueryPath(endpoint, params), "")
}

func (p Private) post(endpoint string, data interface{}) ([]byte, error) {
	marshalData, _ := json.Marshal(data)
	return p.execute(http.MethodPost, endpoint, string(marshalData))
}

func (p Private) delete(endpoint string, data interface{}) ([]byte, error) {
	if data != nil {
		marshalData, _ := json.Marshal(data)
		return p.execute(http.MethodDelete, endpoint, string(marshalData))
	} else {
		return p.execute(http.MethodDelete, endpoint, "")
	}
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

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusIMUsed {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		p.Logger.Printf("wrong status code: %d,err msg:%s", resp.StatusCode, buf.String())
		return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("---->response body:%s\n", respBody)
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

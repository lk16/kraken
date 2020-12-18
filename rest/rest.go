package rest

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	APIUrl     = "https://api.kraken.com"
	APIVersion = "0"
)

type Client struct {
	key           string
	secret        string
	decodedSecret []byte
}

func NewClient() *Client {
	return &Client{}
}

func (client *Client) SetAuth(key string, secret string) error {
	client.key = key
	client.secret = secret

	var err error
	client.decodedSecret, err = base64.StdEncoding.DecodeString(client.secret)
	return err
}

func (client *Client) sign(requestURL string, data url.Values) string {
	// note: calling Write() on a hash object cannot fail, the returned error is always nil

	sha256State := sha256.New()
	sha256State.Write([]byte(data.Get("nonce")))
	sha256State.Write([]byte(data.Encode()))

	hmacState := hmac.New(sha512.New, client.decodedSecret)
	hmacState.Write([]byte(requestURL))
	hmacState.Write(sha256State.Sum(nil))

	return base64.StdEncoding.EncodeToString(hmacState.Sum(nil))
}

func (client *Client) prepareRequest(method string, isPrivate bool, data url.Values) (*http.Request, error) {

	if data == nil {
		data = url.Values{}
	}

	if isPrivate {
		data.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))
	}

	var publicOrPrivate string
	if isPrivate {
		publicOrPrivate = "private"
	} else {
		publicOrPrivate = "public"
	}

	urlPath := fmt.Sprintf("/%s/%s/%s", APIVersion, publicOrPrivate, method)

	request, err := http.NewRequest("POST", APIUrl+urlPath, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	if isPrivate {
		request.Header.Add("API-Key", client.key)
		request.Header.Add("API-Sign", client.sign(urlPath, data))
	}
	return request, nil
}

func parseResponse(response *http.Response, retType interface{}) error {
	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected status code %d", response.StatusCode)
	}

	if response.Body == nil {
		return fmt.Errorf("response body is nil")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body (%s)", err.Error())
	}

	var retData Response
	if retType != nil {
		retData.Result = retType
	}

	if err = json.Unmarshal(body, &retData); err != nil {
		return fmt.Errorf("parsing JSON body failed: %w", err)
	}

	return nil
}

func (client *Client) request(method string, isPrivate bool, data url.Values, retType interface{}) error {
	req, err := client.prepareRequest(method, isPrivate, data)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error during request execution: %s", err.Error())
	}

	defer resp.Body.Close()
	return parseResponse(resp, retType)
}

// GetWebSocketsToken - WebSockets authentication
func (client *Client) GetWebSocketsToken() (WebSocketToken, error) {
	var response WebSocketToken

	if err := client.request("GetWebSocketsToken", true, nil, &response); err != nil {
		return response, err
	}
	return response, nil
}

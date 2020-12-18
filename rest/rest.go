package rest

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
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

func (client *Client) getSign(requestURL string, data url.Values) string {
	sha := sha256.New()

	sha.Write([]byte(data.Get("nonce") + data.Encode()))

	hashData := sha.Sum(nil)
	hmacObj := hmac.New(sha512.New, client.decodedSecret)

	hmacObj.Write(append([]byte(requestURL), hashData...))
	hmacData := hmacObj.Sum(nil)

	return base64.StdEncoding.EncodeToString(hmacData)
}

func (client *Client) prepareRequest(method string, isPrivate bool, data url.Values) (*http.Request, error) {
	if data == nil {
		data = url.Values{}
	}
	requestURL := ""
	if isPrivate {
		requestURL = fmt.Sprintf("%s/%s/private/%s", APIUrl, APIVersion, method)
		data.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))
	} else {
		requestURL = fmt.Sprintf("%s/%s/public/%s", APIUrl, APIVersion, method)
	}
	req, err := http.NewRequest("POST", requestURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("Error during request creation: %s", err.Error())
	}

	if isPrivate {
		urlPath := fmt.Sprintf("/%s/private/%s", APIVersion, method)
		req.Header.Add("API-Key", client.key)
		signature := client.getSign(urlPath, data)
		req.Header.Add("API-Sign", signature)
	}
	return req, nil
}

/* TODO
func (api *Kraken) parseResponse(response *http.Response, retType interface{}) error {
	if response.StatusCode != 200 {
		return fmt.Errorf("Error during response parsing: invalid status code %d", response.StatusCode)
	}

	if response.Body == nil {
		return fmt.Errorf("Error during response parsing: can not read response body")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Error during response parsing: can not read response body (%s)", err.Error())
	}

	var retData KrakenResponse
	if retType != nil {
		retData.Result = retType
	}

	if err = json.Unmarshal(body, &retData); err != nil {
		return fmt.Errorf("Error during response parsing: json marshalling (%s)", err.Error())
	}

	if len(retData.Error) > 0 {
		return fmt.Errorf("Kraken return errors: %s", retData.Error)
	}

	return nil
} */

func (client *Client) request(method string, isPrivate bool, data url.Values, retType interface{}) error {
	req, err := client.prepareRequest(method, isPrivate, data)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error during request execution: %s", err.Error())
	}

	bytes, _ := ioutil.ReadAll(resp.Body) // TODO remove
	log.Printf("got: %s", string(bytes))  // TODO remove

	defer resp.Body.Close()
	return nil // TODO api.parseResponse(resp, retType)
}

// GetWebSocketsToken - WebSockets authentication
func (client *Client) GetWebSocketsToken() (WebSocketTokenResponse, error) {
	var response WebSocketTokenResponse

	if err := client.request("GetWebSocketsToken", true, nil, &response); err != nil {
		return response, err
	}
	return response, nil
}

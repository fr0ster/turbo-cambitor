package spot_rest_api_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/bitly/go-simplejson"
	spot "github.com/fr0ster/turbo-cambitor/rest_api/binance/spot"
	signature "github.com/fr0ster/turbo-signer/signature"
	"github.com/stretchr/testify/assert"
)

var (
	apiKey = os.Getenv("SPOT_TEST_BINANCE_API_KEY")
	secret = os.Getenv("SPOT_TEST_BINANCE_SECRET_KEY")
	sign   = signature.NewSignHMAC(signature.PublicKey(apiKey), signature.SecretKey(secret))
)

func CallRestAPI(req *http.Request) (
	response *simplejson.Json, err error) {
	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Parse the response body as JSON
	response, err = simplejson.NewJson(body)
	if err != nil {
		response = simplejson.New()
		response.Set("message", string(body))
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	return response, nil
}

func TestAccount(t *testing.T) {
	// ra := spot.New(sign, true)
	// response, err := ra.Account().SetTimestamp().SetSignature().Do()
	// assert.NoError(t, err)
	// assert.NotNil(t, response)

	params := simplejson.New()
	params.Set("timestamp", int64(time.Nanosecond)*time.Now().UnixNano()/int64(time.Millisecond))
	params, _ = sign.SignParameters(params)
	paramsStr := ""
	for key, value := range params.MustMap() {
		paramsStr += fmt.Sprintf("%v=%v&", key, value)
	}
	paramsStr = paramsStr[:len(paramsStr)-1]
	baseUrl := "https://testnet.binance.vision"
	endpoint := "/api/v3/account"
	if paramsStr != "" {
		endpoint = endpoint + "?" + paramsStr
	}
	req, err := http.NewRequest("GET", baseUrl+endpoint, nil)
	req.Header.Set("X-MBX-APIKEY", apiKey)
	assert.NoError(t, err)
	response, err := CallRestAPI(req)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestQueryCurrentOrderCountUsage(t *testing.T) {
	ra := spot.New(sign, true)
	response, err := ra.QueryCurrentOrderCountUsage().SetTimestamp().SetSignature().Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestQueryAllocations(t *testing.T) {
	ra := spot.New(sign, true)
	response, err := ra.QueryAllocations("BTCUSDT").SetTimestamp().SetSignature().Do()
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

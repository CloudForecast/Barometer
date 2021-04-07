package barometerApi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"time"
)

const BAROMETER_API_URL = "http://localhost"

var barometerApiKey string

func init() {
	barometerApiKey = viper.GetString("apiKey")
}

type ApiClient interface {
	makePostRequest(interface{}) (int, error)
	makeGetRequest(string) ([]byte, error)

	GetApiKey() string
	GetPromQlInstructions() (*PromQlQueryInstruction, error)

	SendHealthCheckEvent() error
}

type BarometerApi struct {
	barometerApiKey string
}

func NewBarometerApi(apiKey string) BarometerApi {
	return BarometerApi{barometerApiKey: apiKey}
}

func (b BarometerApi) GetApiKey() string {
	return b.barometerApiKey
}

func (b BarometerApi) makeGetRequest(path string) ([]byte, error) {
	var request *http.Request
	timeout := 5 * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", fmt.Sprint(BAROMETER_API_URL, path), bytes.NewBuffer([]byte{}))
	if err != nil {
		return []byte{}, err
	}
	request.Header.Set("X-API-Key", barometerApiKey)
	resp, err := client.Do(request)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func (b BarometerApi) makePostRequest(payload interface{}) (statusCode int, err error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return
	}

	timeout := 5 * time.Second
	client := http.Client{
		Timeout: timeout,
	}

    var request *http.Request
	if request, err = http.NewRequest("POST", BAROMETER_API_URL, bytes.NewBuffer(jsonData)); err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-Key", barometerApiKey)

	dryRun := viper.GetBool("dryrun")
	if dryRun {
		return 200, nil
	}

	resp, err := client.Do(request)
	if err != nil {
		return
	}

	statusCode = resp.StatusCode
	return
}


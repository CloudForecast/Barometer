package barometerApi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"time"
)

const BAROMETER_API_URL = "https://barometer.perfectweather.io"

type ApiClient interface {
	makePostRequest(interface{}) (int, error)
	makeGetRequest(string) ([]byte, error)

	GetApiKey() string
	GetPromQlInstructions() (*PromQlQueryInstruction, error)
	GetKubeInstructions() (*KubeQueryInstruction, error)

	SendHealthCheckEvent() error
	SendK8sAPIResultsEvent(BarometerK8sApiResultsEventData) error
	SendExceptionEvent(inputError error) error
}

type BarometerApi struct {
	barometerApiKey string
	clusterUUID string
	HTTPClient *http.Client
}

func NewBarometerApi(apiKey string, clusterUUID string) BarometerApi {
	return BarometerApi{
		barometerApiKey: apiKey,
		clusterUUID: clusterUUID,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (b BarometerApi) GetApiKey() string {
	return b.barometerApiKey
}

func (b BarometerApi) makeGetRequest(path string) ([]byte, error) {
	var request *http.Request

	request, err := http.NewRequest("GET", fmt.Sprint(BAROMETER_API_URL, path), nil)
	if err != nil {
		return []byte{}, err
	}
	request.Header.Set("BM-API-Key", b.barometerApiKey)
	request.Header.Set("BM-Cluster-UUID", b.clusterUUID)
	resp, err := b.HTTPClient.Do(request)
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
	log.Debug().Msgf("before converting payload to JSON: %v", payload)
	jsonData, err := json.Marshal(payload)
	log.Debug().Msgf("POSTing this JSON: %s", jsonData)
	if err != nil {
		return
	}

    var request *http.Request
	if request, err = http.NewRequest("POST", BAROMETER_API_URL, bytes.NewBuffer(jsonData)); err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("BM-API-Key", b.barometerApiKey)
	request.Header.Set("BM-Cluster-UUID", b.clusterUUID)

	dryRun := viper.GetBool("dryrun")
	if dryRun {
		return 200, nil
	}

	resp, err := b.HTTPClient.Do(request)
	if err != nil {
		return
	}

	statusCode = resp.StatusCode
	return
}


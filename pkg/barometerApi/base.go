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

const BAROMETER_API_EVENTS_PATH = "/api/barometer/v1/events"

// This should likely be broken up into smaller client interfaces,
// as we only have it as an interface for now to enable easier testing later
type ApiClient interface {
	makePostRequest(interface{}) (int, error)
	makeGetRequest(string) ([]byte, error)

	GetApiKey() string
	GetPromQlInstructions() (*PromQlQueryInstruction, error)
	GetKubeInstructions() (*KubeQueryInstruction, error)

	SendHealthCheckEvent() error
	SendK8sAPIResultsEvent(BarometerK8sApiResultsEventData) error
	SendPromQlResultsEvent(PromQlQueryInstruction, []PromQLResult) error
	SendExceptionEvent(inputError error) error
}

type BarometerApi struct {
	barometerApiKey string
	clusterUUID     string
	HTTPClient      *http.Client
	ApiHost 		string
}

func NewBarometerApi(apiKey string, clusterUUID string) BarometerApi {
	apiHost := viper.GetString("apiHost")
	return BarometerApi{
		barometerApiKey: apiKey,
		clusterUUID:     clusterUUID,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		ApiHost: apiHost,
	}
}

func (b BarometerApi) GetApiKey() string {
	return b.barometerApiKey
}

func (b BarometerApi) makeGetRequest(path string) ([]byte, error) {
	var request *http.Request

	request, err := http.NewRequest("GET", fmt.Sprint(b.ApiHost, path), nil)
	if err != nil {
		return []byte{}, err
	}
	request.Header.Set("bm-api-key", b.barometerApiKey)
	request.Header.Set("bm-cluster-uuid", b.clusterUUID)
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
	log.Trace().Msgf("before converting payload to JSON: %v", payload)
	jsonData, err := json.Marshal(payload)
	log.Trace().Msgf("POSTing this JSON: %s", jsonData)
	if err != nil {
		return
	}

	var request *http.Request
	urlPath := fmt.Sprint(b.ApiHost, BAROMETER_API_EVENTS_PATH)
	if request, err = http.NewRequest("POST", urlPath, bytes.NewBuffer(jsonData)); err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("bm-api-key", b.barometerApiKey)
	request.Header.Set("bm-cluster-uuid", b.clusterUUID)

	dryRun := viper.GetBool("dryrun")
	if dryRun {
		return 200, nil
	}

	resp, err := b.HTTPClient.Do(request)
	if err != nil {
		return
	}

	statusCode = resp.StatusCode
	log.Debug().Msgf("POST request result: %s", statusCode)
	return
}

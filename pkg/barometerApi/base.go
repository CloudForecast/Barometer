package barometerApi

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"time"
)

const BAROMETER_API_EVENTS_PATH = "/api/barometer/v2/events"

// This should likely be broken up into smaller client interfaces,
// as we only have it as an interface for now to enable easier testing later
type ApiClient interface {
	makePostRequest(interface{}) (int, error)
	makeGetRequest(string) ([]byte, error)

	GetApiKey() string
	GetPromQlInstructions() (*PromQlQueryInstruction, error)
	GetKubeInstructions() (*KubeQueryInstruction, error)

	SendHealthCheckEvent(info map[string]interface{}) error
	SendK8sAPIResults(*KubeQueryInstruction, BarometerK8sApiResultsEventData) error
	SendPromQlResults(PromQlQueryInstruction, PromQLResultsWrapper) error
	SendExceptionEvent(inputError error) error
	SendUploadedDataEvent(eventUUID string, uploadedDataType BarometerEventType) error

	UploadDataToS3(configuration UploadConfiguration, data interface{}) (int, error)
}

type BarometerApi struct {
	barometerApiKey string
	clusterUUID     string
	HTTPClient      *http.Client
	ApiHost         string
}

type UploadConfiguration struct {
	EventUUID      string `json:"event_uuid"`
	S3PreSignedUrl string `json:"s3_pre_signed_url"`
}

func NewBarometerApi(apiKey string, clusterUUID string) BarometerApi {
	apiHost := viper.GetString("apiHost")
	return BarometerApi{
		barometerApiKey: apiKey,
		clusterUUID:     clusterUUID,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
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

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(jsonData)
	gz.Close()

	var request *http.Request
	urlPath := fmt.Sprint(b.ApiHost, BAROMETER_API_EVENTS_PATH)
	if request, err = http.NewRequest("POST", urlPath, &buf); err != nil {
		return
	}
	request.Header.Set("Content-Encoding", "gzip")
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
	log.Debug().Msgf("POST request result: %v", statusCode)
	return
}

func (b BarometerApi) SendUploadedDataEvent(eventUUID string, uploadedDataType BarometerEventType) error {
	log.Debug().Msg("Sending Uploaded Data Event...")
	eventData := map[string]interface{}{"event_uuid": eventUUID, "uploaded_data_type": uploadedDataType}
	event := BarometerEvent{
		EventType: UploadedData,
		EventTs:   time.Now().Unix(),
		Event:     eventData,
	}
	statusCode, err := b.makePostRequest(event)
	if err != nil {
		return err
	}
	if statusCode != 200 {
		return fmt.Errorf("received unexpected status code %d from UploadedDataEvent", statusCode)
	}
	log.Info().Msgf("Uploaded Data Done - eventUUID: %s", eventUUID)
	return nil
}

func (b BarometerApi) UploadDataToS3(configuration UploadConfiguration, payload interface{}) (statusCode int, err error) {
	log.Trace().Msgf("before converting payload to JSON: %v", payload)
	jsonData, err := json.Marshal(payload)
	log.Trace().Msgf("PUTing this JSON to S3 url (%s): %s",  configuration.S3PreSignedUrl, jsonData)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(jsonData)
	gz.Close()

	var request *http.Request
	if request, err = http.NewRequest("PUT", configuration.S3PreSignedUrl, &buf); err != nil {
		return
	}
	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Content-Type", "application/json")

	resp, err := b.HTTPClient.Do(request)
	if err != nil {
		return
	}

	statusCode = resp.StatusCode
	body, err := ioutil.ReadAll(resp.Body)
	log.Info().Msgf("Data was sent to S3; statusCode:%v; body:%s", statusCode, body)
	return
}
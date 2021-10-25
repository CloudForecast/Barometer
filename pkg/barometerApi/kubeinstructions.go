package barometerApi

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"time"
)

const KUBERNETES_INSTRUCTION_PATH = "/api/barometer/v2/kubectl_instructions"

type KubeKind string

type KubeQueryInstruction struct {
	KindsToFetch []KubeKind `json:"kubectl_get_queries"`
	UploadConfiguration UploadConfiguration `json:"upload_configuration"`
}

func NewK8sApiResultsEvent(d BarometerK8sApiResultsEventData) BarometerEvent {
	var outputMap map[string]interface{}
	err := mapstructure.Decode(d, &outputMap)
	if err != nil {
		panic(err)
	}
	return BarometerEvent{
		EventType: K8sApiResults,
		EventTs:   time.Now().Unix(),
		Event:     outputMap,
	}
}

func (b BarometerApi) GetKubeInstructions() (*KubeQueryInstruction, error) {
	response, err := b.makeGetRequest(KUBERNETES_INSTRUCTION_PATH)
	if err != nil {
		return nil, err
	}

	var instructions KubeQueryInstruction
	if err := json.Unmarshal(response, &instructions); err != nil {
		return nil, err
	}
	return &instructions, nil
}

func (b BarometerApi) SendK8sAPIResults(instruction *KubeQueryInstruction, eventData BarometerK8sApiResultsEventData) error {
	log.Info().Msg("Sending k8s API data...")
	event := NewK8sApiResultsEvent(eventData)

	// Upload data to S3
	statusCode, err := b.UploadDataToS3(instruction.UploadConfiguration, event)
	if err != nil {
		return err
	}
	if statusCode != 200 {
		return fmt.Errorf("received unexpected status code %d from SendK8sAPIResults when sending the data to S3", statusCode)
	}

	// Send UploadedDataEvent
	err = b.SendUploadedDataEvent(instruction.UploadConfiguration.EventUUID, K8sApiResults)
	if err != nil {
		return err
	}
	return nil
}

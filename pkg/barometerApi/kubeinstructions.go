package barometerApi

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

const KUBERNETES_INSTRUCTION_PATH = "/api/barometer/v1/kubectl_instructions"

type KubeKind string

type KubeQueryInstruction struct {
	KindsToFetch []KubeKind `json:"kubectl_get_queries"`
}

func NewK8sApiResultsEvent(d BarometerK8sApiResultsEventData) BarometerEvent {
	var outputMap map[string]interface{}
	err := mapstructure.Decode(d, &outputMap)
	if err != nil {
		panic(err)
	}
	return BarometerEvent{
		EventType: K8sApiResults,
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

func (b BarometerApi) SendK8sAPIResultsEvent(eventData BarometerK8sApiResultsEventData) error {
	log.Debug().Msg("Sending k8s API data...")
	event := NewK8sApiResultsEvent(eventData)
	statusCode, err := b.makePostRequest(event)
	if err != nil {
		return err
	}
	if statusCode != 200 {
		return fmt.Errorf("received unexpected status code %d from kubernetes API data submission", statusCode)
	}
	return nil
}

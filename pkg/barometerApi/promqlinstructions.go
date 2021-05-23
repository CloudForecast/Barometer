package barometerApi

import (
	"encoding/json"
)

const PromqlInstructionPath = "/api/barometer/v1/promql_instructions"

type PromQlConfiguration struct {
	StartTs int `json:"start_ts"`
	EndTs   int `json:"end_ts"`
	StepSec int `json:"step_sec"`
}

// PromQlQueryName is the name of an individual Prometheus query provided by the Barometer API.
type PromQlQueryName string

// PromQlQuery is a string containing the query itself to send to Prometheus.
type PromQlQuery string

type PromQlQueryInstruction struct {
	Configurations []PromQlConfiguration           `json:"promql_configurations"`
	Queries        map[PromQlQueryName]PromQlQuery `json:"promql_queries"`
}

// GetPromQlInstructions fetches the list of PromQL queries to execute for this cluster
// from the Barometer API.
func (b BarometerApi) GetPromQlInstructions() (*PromQlQueryInstruction, error) {
	response, err := b.makeGetRequest(PromqlInstructionPath)
	if err != nil {
		return nil, err
	}

	var instructions PromQlQueryInstruction
	if err := json.Unmarshal(response, &instructions); err != nil {
		return nil, err
	}
	return &instructions, nil
}

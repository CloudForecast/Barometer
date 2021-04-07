package barometerApi

import "encoding/json"

const PROMQL_INSTRUCTION_PATH = "/promql"

type PromQlConfiguration struct {
	StartTs int `json:"start_ts"`
	EndTs int `json:"end_ts"`
	StepSec int `json:"step_sec"`
}

type PromQlQueryInstruction struct {
	Configurations []PromQlConfiguration `json:"promql_configurations"`
	Queries map[string]string `json:"promql_queries"`
}

func (b BarometerApi) GetPromQlInstructions() (*PromQlQueryInstruction, error) {
	response, err := b.makeGetRequest(PROMQL_INSTRUCTION_PATH)
	if err != nil {
		return nil, err
	}

	var instructions PromQlQueryInstruction
	if err := json.Unmarshal(response, &instructions); err != nil {
		return nil, err
	}
	return &instructions, nil
}
package barometerApi

import (
	"encoding/json"
	"testing"
)

func TestKubeInstructionsUnmarshal(t *testing.T) {
	const KUBE_INSTRUCTIONS_RESPONSE = `{"kubectl_get_queries":["namespaces","services","pods","nodes","persistentvolumes","ingress"]}`

	var instructions KubeQueryInstruction
	err := json.Unmarshal([]byte(KUBE_INSTRUCTIONS_RESPONSE), &instructions)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
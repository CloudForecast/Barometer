package barometerApi

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKubeInstructionsUnmarshal(t *testing.T) {
	const KUBE_INSTRUCTIONS_RESPONSE = `{ "kubectl_get_queries": ["namespaces","services","pods","nodes","persistentvolumes","ingress"], "upload_configuration": { "event_uuid": "13f6b8dd-72fa-4872-97cd-6cabbaa4f5d4", "s3_pre_signed_url": "https://s3.url/13f6b8dd-72fa-4872-97cd-6cabbaa4f5d4.json?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAI6JZBXNP3Q4QKDIA%2F20211023%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20211023T223241Z&X-Amz-Expires=300&X-Amz-SignedHeaders=host&X-Amz-Signature=0820e7fcc8cf07f6ed67401c2d750a103db73538e2577c7e4599efbc0d54daa7" } }`

	var instructions KubeQueryInstruction
	err := json.Unmarshal([]byte(KUBE_INSTRUCTIONS_RESPONSE), &instructions)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	} else {
		assert.Equal(t, instructions.UploadConfiguration.EventUUID, "13f6b8dd-72fa-4872-97cd-6cabbaa4f5d4", "UUID should be '13f6b8dd-72fa-4872-97cd-6cabbaa4f5d4'")
		assert.Equal(t, instructions.UploadConfiguration.S3PreSignedUrl, "https://s3.url/13f6b8dd-72fa-4872-97cd-6cabbaa4f5d4.json?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAI6JZBXNP3Q4QKDIA%2F20211023%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20211023T223241Z&X-Amz-Expires=300&X-Amz-SignedHeaders=host&X-Amz-Signature=0820e7fcc8cf07f6ed67401c2d750a103db73538e2577c7e4599efbc0d54daa7")
		assert.Equal(t, len(instructions.KindsToFetch), 6, "Should have 6 kinds")
	}
}

package barometerApi

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPromQLInstructionsUnmarshal(t *testing.T) {
	const mockApiResponse = `{"promql_configurations": [{"start_ts": 1635019200,"end_ts": 1635026400,"step_sec": 3600,"upload_configuration": {"event_uuid": "4d1228f9-8e73-42ee-90c5-22ebe601c2a0","s3_pre_signed_url": "https://s3.url/4d1228f9-8e73-42ee-90c5-22ebe601c2a0.json?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAI6JZBXNP3Q4QKDIA%2F20211023%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20211023T223427Z&X-Amz-Expires=300&X-Amz-SignedHeaders=host&X-Amz-Signature=734da9a884f9f34848d5a1ba9184ebb0d9db2886f8f0e4847fb99b2cb138971f"}}],"promql_queries": {"container_cpu_request": "sum(max_over_time(kube_pod_container_resource_requests{unit='core'}[1h])) by (namespace, service, pod, container)","container_cpu_usage": "sum(irate(container_cpu_usage_seconds_total[1h])) by (namespace, service, pod, container)","container_memory_request": "sum(max_over_time(kube_pod_container_resource_requests{unit='byte'}[1h])) by (namespace, service, pod, container)","container_memory_usage": "sum(max_over_time(container_memory_working_set_bytes[1h])) by (namespace, service, pod, container)","container_disc_read": "sum(rate(container_fs_reads_total[1h])) by (namespace, service, pod, container)","container_disc_write": "sum(rate(container_fs_writes_total[1h])) by (namespace, service, pod, container)","pod_start_time": "avg(avg_over_time(kube_pod_start_time[1h])) by (namespace, service, pod, container)","pod_interval_last_timestamp": "max(max_over_time(timestamp(kube_pod_start_time)[1h:])) by (namespace, service, pod, container)","pod_traffic_in": "sum(increase(container_network_receive_bytes_total{namespace!='', pod!='', instance!=''}[1h])) by (namespace, pod)","pod_traffic_out": "sum(increase(container_network_transmit_bytes_total{namespace!='', pod!='', instance!=''}[1h])) by (namespace, pod)","namespace_labels": "rate(kube_namespace_labels[1h])","service_labels": "rate(kube_service_labels[1h])","service_info": "rate(kube_service_info[1h])","pod_labels": "rate(kube_pod_labels[1h])","pod_info": "rate(kube_pod_info[1h])","node_labels": "rate(kube_node_labels[1h])","node_info": "rate(kube_node_info[1h])","node_cpu_usage": "sum(kube_pod_container_resource_requests{resource=\"cpu\"}) by (kubernetes_node, node)","node_cpu_limit": "sum(kube_node_status_allocatable{resource='cpu', unit='core'}) by (kubernetes_node, node)","node_memory_usage": "sum(kube_pod_container_resource_requests{resource=\"memory\"}) by (kubernetes_node, node)","node_memory_limit": "sum(kube_node_status_allocatable{resource=\"memory\", unit=\"byte\"}) by (kubernetes_node, node)","prometheus_target_interval_length_seconds": "prometheus_target_interval_length_seconds","scrape_duration_seconds": "scrape_duration_seconds"}}`
	var instructions PromQlQueryInstruction

	err := json.Unmarshal([]byte(mockApiResponse), &instructions)

	if err != nil {
		t.Errorf("received unexepected error: %v", err)
	} else {
		assert.Equal(t, len(instructions.Configurations), 1, "Should contains one configuration")
		assert.Equal(t, instructions.Configurations[0].StartTs, 1635019200, "StartTs should be 1635019200")
		assert.Equal(t, instructions.Configurations[0].EndTs, 1635026400, "EndTs should be 1635026400")
		assert.Equal(t, instructions.Configurations[0].StepSec, 3600, "StepSec should be 3600")
		assert.Equal(t, instructions.Configurations[0].UploadConfiguration.EventUUID, "4d1228f9-8e73-42ee-90c5-22ebe601c2a0", "UUID should be '4d1228f9-8e73-42ee-90c5-22ebe601c2a0'")
		assert.Equal(t, instructions.Configurations[0].UploadConfiguration.S3PreSignedUrl, "https://s3.url/4d1228f9-8e73-42ee-90c5-22ebe601c2a0.json?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAI6JZBXNP3Q4QKDIA%2F20211023%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20211023T223427Z&X-Amz-Expires=300&X-Amz-SignedHeaders=host&X-Amz-Signature=734da9a884f9f34848d5a1ba9184ebb0d9db2886f8f0e4847fb99b2cb138971f", "S3 Presigned shoudl be 'https://dev-barometer-event-raw-data.s3.amazonaws.com/tmp/c258/bmc40/4d1228f9-8e73-42ee-90c5-22ebe601c2a0.json?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAI6JZBXNP3Q4QKDIA%2F20211023%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20211023T223427Z&X-Amz-Expires=300&X-Amz-SignedHeaders=host&X-Amz-Signature=734da9a884f9f34848d5a1ba9184ebb0d9db2886f8f0e4847fb99b2cb138971f'")
		assert.Equal(t, len(instructions.Queries), 23, "Should have 10 queries")
	}
}

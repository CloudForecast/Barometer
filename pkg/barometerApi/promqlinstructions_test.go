package barometerApi

import (
	"encoding/json"
	"testing"
)

func TestPromQLInstructionsUnmarshal(t *testing.T) {
	const mockApiResponse = `{"promql_configurations":[{"start_ts":1619434800,"end_ts":1619445600,"step_sec":300},{"start_ts":1619424000,"end_ts":1619434800,"step_sec":300}],"promql_queries":{"container_cpu_requested":"sum(max_over_time(kube_pod_container_resource_requests_cpu_cores[5m])) by (namespace, service, pod, container, node)","container_cpu_usage":"sum(irate(container_cpu_usage_seconds_total[5m])) by (namespace, service, pod, container, node)","container_memory_request":"sum(max_over_time(kube_pod_container_resource_requests_memory_bytes[5m])) by (namespace, service, pod, container, node)","container_memory_usage":"sum(irate(container_memory_working_set_bytes[5m])) by (namespace, service, pod, container, node)","container_traffic_in":"sum(irate(container_network_receive_bytes_total[5m])) by (namespace, service, pod, container, node)","container_traffic_out":"sum(irate(container_network_transmit_bytes_total[5m])) by (namespace, service, pod, container, node)","container_disc_read":"sum(irate(container_fs_reads_total[5m])) by (namespace, service, pod, container, node)","container_disc_write":"sum(irate(container_fs_writes_total[5m])) by (namespace, service, pod, container, node)","namespace_labels":"rate(kube_namespace_labels[5m])","service_labels":"rate(kube_service_labels[5m])","service_info":"rate(kube_service_info[5m])","pod_labels":"rate(kube_pod_labels[5m])","pod_info":"rate(kube_pod_info[5m])","node_labels":"rate(kube_node_labels[5m])","node_info":"rate(kube_node_info[5m])","node_cpu_usage":"sum (rate (container_cpu_usage_seconds_total[5m])) by (node)","node_cpu_requested":"sum(machine_cpu_cores) by (node)","node_memory_usage":"sum(rate(container_memory_usage_bytes[5m])) by (node)","node_memory_limit":"sum(machine_memory_bytes) by (node)"}}`
	var instructions PromQlQueryInstruction

	err := json.Unmarshal([]byte(mockApiResponse), &instructions)
	if err != nil {
		t.Errorf("received unexepected error: %v", err)
	}
}

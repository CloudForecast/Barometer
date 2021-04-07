package barometerApi

type BarometerEventType string
const (
	HealthCheck   BarometerEventType = "health_check"
	Exception     BarometerEventType = "exception"
	PromQlResults BarometerEventType = "promql_results"
	K8sApiResults BarometerEventType = "k8s_api_results"
)

type BarometerEvent struct {
	EventKey BarometerEventType     `json:"event_key"`
	Event    map[string]interface{} `json:"event"`
}

type BarometerHealthCheckEventData struct {}

type BarometerExceptionEventData struct {
	Message string `json:"message"`
}

type BarometerPromQlResultsEventData struct {
	PromQlResults []PromQLResult `json:"promql_results"`
}

type BarometerK8sApiResultsEventData struct {
	K8sClusterInformation map[string]interface{} `json:"k8s_cluster_information"`
}

type PromQLResult struct {
	QueryId string `json:"query_id"`
	Query string `json:"query"`
	PromQlConfiguration map[string]interface{} `json:"promql_configuration"`
	// TODO: type this better
	Result []map[string]interface{} `json:"result"`
}


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
	Event    map[string]interface{} `json:"event_data"`
}

type BarometerHealthCheckEventData struct {}

type BarometerExceptionEventData struct {
	ErrorBool bool `json:"error" mapstructure:"error"`
	Message string `json:"message" mapstructure:"message"`
}

// Error implements the error interface on BarometerExceptionEventData
//  so it can be passed around or used to wrap other errors conveniently
func (d BarometerExceptionEventData) Error() string {
	return d.Message
}

type BarometerPromQlResultsEventData struct {
	PromQlResults []PromQLResult `json:"promql_results" mapstructure:"promql_results"`
}

type BarometerK8sApiResultsEventData struct {
	K8sClusterInformation map[string][]interface{} `json:"k8s_cluster_information" mapstructure:"k8s_cluster_information"`
}

type PromQLResult struct {
	QueryId string `json:"query_id"`
	Query string `json:"query"`
	PromQlConfiguration map[string]interface{} `json:"promql_configuration"`
	// TODO: type this better
	Result []map[string]interface{} `json:"result"`
}


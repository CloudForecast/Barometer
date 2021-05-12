package barometerApi

func newPromQlResultsEvent(results ...PromQLResult) *BarometerPromQlResultsEventData {
	return &BarometerPromQlResultsEventData{
		PromQlResults: results,
	}
}
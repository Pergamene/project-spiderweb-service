package nextbatch

// NextBatch represents the nextBatch object that is passed back with GET requests for a list of entities.
type NextBatch struct {
	ParamKey   string `json:"paramKey"`
	ParamValue string `json:"paramValue"`
}
